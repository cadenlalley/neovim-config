package api

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/kitchens-io/kitchens-api/internal/ai"
	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/rs/zerolog/log"
)

var testApp *App
var testDSN string

var testBearerToken = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjNJS090Z0ZtWi02ejYzMjdFWWFYdiJ9.eyJuaWNrbmFtZSI6InRlc3Qtc2VydmljZSIsIm5hbWUiOiJ0ZXN0LXNlcnZpY2VAa2l0Y2hlbnMtYXBwLmNvbSIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci8yY2M3ODg2Y2Y1NTUwMDliZTNmNjE3ZjM5NDMyNjAxMj9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRnRlLnBuZyIsInVwZGF0ZWRfYXQiOiIyMDI1LTA0LTEyVDE0OjQ2OjI1LjY2MVoiLCJlbWFpbCI6InRlc3Qtc2VydmljZUBraXRjaGVucy1hcHAuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vZGV2LWI3bmRvN2l5MWNhbGl5bzQudXMuYXV0aDAuY29tLyIsImF1ZCI6IkZxYjNZcG85TFRSOThOb0M5QktZcE5hT2kyZEN0MThyIiwic3ViIjoiYXV0aDB8NjY1ZTM2NDYxMzlkOWY2MzAwYmFkNWU5IiwiaWF0IjoxNzQ0NDY5MTg3LCJleHAiOjE3NDUwNzM5ODcsInNpZCI6InlfSjkzWkFqQ2Z1UnhUYWZXX3NWcGRHbmplemd1MmpsIiwibm9uY2UiOiJYMGxIY2xsVmEyUXlkbVJEWlZkR01tZ3hOMFJLZGpkWE4wVm1hazkxT0dWUGMxUTNORE5FYkY5elNBPT0ifQ.n1G3wQt3SlzxP111TgZlASRN978vImYqzpCL3kDt32e3NCLMznT4y3juy6ujduz0iB42aCCy3ffQYEn2UyKWJMtzzNXKZOOUBNLQ1iIUMRqwVOxPZr1UgBjxrRYOYrecVn3MSCLzPcNPMLT84LK44b4I16yxJYQA8jLdAH6GM9V7cIHs00jjrLa95psUATSwlgnjm76diGHOO8JPzFUNwl28hfsyv4C_pTipfnZECgaxIIGQbaTcWEKEtWzzj-rFGZVX1Syt_VJA6Sk2Cwr0PWcvuQdLBh-vrR2TBPkPwvuly1ZQO_IfyAe7XkrzYzH5tUUZksQkeRl7Ajpd9hBm3Q"

// Application Configuration
type AppConfig struct {
	Debug           bool          `default:"false" envconfig:"DEBUG"`
	Port            string        `default:"1313" envconfig:"PORT"`
	ShutdownTimeout time.Duration `default:"10s" envconfig:"SHUTDOWN_TIMEOUT"`
	Env             string        `default:"test" envconfig:"APP_ENV"`

	// Database configurations
	DB struct {
		User string `required:"true" envconfig:"DB_USER"`
		Pass string `required:"true" envconfig:"DB_PASS"`
		Host string `required:"true" envconfig:"DB_HOST"`
		Name string `required:"true" envconfig:"DB_NAME"`
	}

	Migrations struct {
		Fixtures *string `default:"fixtures_migrations" envconfig:"MIGRATIONS_FIXTURES"`
	}

	// S3 Object Storage configurations
	S3 struct {
		Host        string `envconfig:"S3_LOCAL_HOST"`
		MediaBucket string `required:"true" envconfig:"S3_MEDIA_BUCKET"`
	}

	CDN struct {
		Host string `required:"true" envconfig:"CDN_HOST"`
	}

	// OpenAI
	OpenAI struct {
		Host  string `required:"true" envconfig:"OPENAI_HOST"`
		Token string `required:"true" envconfig:"OPENAI_TOKEN"`
	}

	// Brave Search
	Brave struct {
		Host  string `required:"true" envconfig:"BRAVE_HOST"`
		Token string `required:"true" envconfig:"BRAVE_TOKEN"`
	}
}

func TestMain(m *testing.M) {
	// .env file is optional, so ignore the error returned from loading.
	_ = godotenv.Load(filepath.Join("../../", ".env"))

	// Parse environemnt variables into the configuration struct.
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal().Err(err).Msg("could not parse application config")
	}

	// AWS Configurations
	//===========================================
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().Err(err).Msg("could not load aws configurations")
	}

	// Handle database migrations and connections.
	// ===========================================
	testDSN = mysql.DSN(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Name)

	db, err := mysql.Connect(testDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("could not start database")
	}
	defer db.Close()

	// Handle Object storage
	// ==========================
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.S3.Host != "" {
			o.BaseEndpoint = &cfg.S3.Host
			o.UsePathStyle = true
		}
	})

	fileManager := media.NewS3FileManager(s3Client, cfg.S3.MediaBucket)

	// Handle AI Client
	// ==========================
	aiClient := ai.NewClient(cfg.OpenAI.Token, cfg.OpenAI.Host)

	// Handle application server.
	// ==========================

	// Create an API instance.
	testApp = Create(CreateInput{
		DB:            db,
		FileManager:   fileManager,
		AuthValidator: nil,
		Env:           cfg.Env,
		CDNHost:       cfg.CDN.Host,
		AIClient:      aiClient,
	})

	// Run all tests
	exitCode := m.Run()

	// Exit with the test result code
	os.Exit(exitCode)
}

func resetFixtures() error {
	return mysql.ResetFixtures("file://../../fixtures", testDSN, nil)
}

func request(method, path string, payload interface{}) (status int, body *bytes.Buffer, err error) {
	var data []byte
	if payload != nil {
		data, err = json.Marshal(payload)
		if err != nil {
			return 0, nil, err
		}
	}

	// Create a request
	req := httptest.NewRequest(method, path, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", testBearerToken)

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the handler directly
	testApp.API.ServeHTTP(w, req)

	// Return the response details
	return w.Code, w.Body, err
}

func formRequest(method, path string, payload map[string]string) (status int, body *bytes.Buffer, err error) {
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	for k, v := range payload {
		writer.WriteField(k, v)
	}

	err = writer.Close()
	if err != nil {
		return 0, nil, err
	}

	// Create a request
	req := httptest.NewRequest(method, path, data)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", testBearerToken)

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the handler directly
	testApp.API.ServeHTTP(w, req)

	// Return the response details
	return w.Code, w.Body, err
}
