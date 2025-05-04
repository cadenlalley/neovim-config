package main

import (
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/folders"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/kitchens-io/kitchens-api/pkg/sdk"

	"github.com/rs/zerolog/log"
)

// Application Configuration
type AppConfig struct {
	Source string `default:"https://dbm2jjspeb.us-east-1.awsapprunner.com"`
	Target string `default:"http://localhost:1313"`
	JWT    string `required:"true" envconfig:"JWT"`

	// Database configurations
	DB struct {
		User string `required:"true" envconfig:"DB_USER"`
		Pass string `required:"true" envconfig:"DB_PASS"`
		Host string `required:"true" envconfig:"DB_HOST"`
		Name string `required:"true" envconfig:"DB_NAME"`
	}
}

func main() {
	// .env file is optional, so ignore the error returned from loading.
	_ = godotenv.Load()

	// Parse environemnt variables into the configuration struct.
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal().Err(err).Msg("could not parse application config")
	}

	// Handle database connections.
	// ===========================================
	db, err := mysql.Connect(mysql.DSN(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Name))
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to database")
	}
	defer db.Close()

	sourceSDK := sdk.NewClient(sdk.NewClientOptions{
		Host: cfg.Source,
	})
	sourceSDK.SetAuth(cfg.JWT)

	targetSDK := sdk.NewClient(sdk.NewClientOptions{
		Host: cfg.Target,
	})
	targetSDK.SetAuth(cfg.JWT)

	// Sync Accounts, Kitchens, and Folders to database.
	// These are handled directly due to content-type limitations on the API.
	// =========================================
	accountRes, err := sourceSDK.AdminListAccounts()
	if err != nil {
		log.Fatal().Err(err).Msg("could not list accounts")
	}

	var accountsCount int
	var testAccountID string
	for _, account := range accountRes {
		// Skip test account
		if account.UserID == "auth0|665e3646139d9f6300bad5e9" {
			testAccountID = account.AccountID
			continue
		}

		_, err := accounts.CreateAccount(context.TODO(), db, accounts.CreateAccountInput{
			AccountID: account.AccountID,
			UserID:    account.UserID,
			Email:     account.Email,
			FirstName: account.FirstName,
			LastName:  account.LastName,
		})
		if err != nil {
			log.Fatal().Err(err).Msgf("could not sync account '%s'", account.AccountID)
		}

		accountsCount++
	}
	log.Info().Int("accounts", accountsCount).Msg("synced accounts")

	kitchensRes, err := sourceSDK.SearchKitchens()
	if err != nil {
		log.Fatal().Err(err).Msg("could not list kitchens")
	}

	// Create Kitchen IDs to track for kitchen dependencies syncing.
	kitchenIDs := make([]string, 0)

	var kitchensCount int
	for _, kitchen := range kitchensRes {
		// Skip test account
		if kitchen.AccountID == testAccountID {
			continue
		}

		_, err := kitchens.CreateKitchen(context.TODO(), db, kitchens.CreateKitchenInput{
			KitchenID:   kitchen.KitchenID,
			AccountID:   kitchen.AccountID,
			KitchenName: kitchen.Name,
			Bio:         kitchen.Bio.String,
			Handle:      kitchen.Handle,
			Avatar:      kitchen.Avatar.String,
			Cover:       kitchen.Cover.String,
			Private:     kitchen.Private,
		})
		if err != nil {
			log.Fatal().Err(err).Msgf("could not sync kitchen '%s'", kitchen.KitchenID)
		}

		kitchensCount++
		kitchenIDs = append(kitchenIDs, kitchen.KitchenID)
	}
	log.Info().Int("kitchens", kitchensCount).Msg("synced kitchens")

	// Sync Kitchen dependencies to database.
	// =========================================
	for _, kitchenID := range kitchenIDs {

		// Sync Recipes to database.
		// =========================================
		recipesRes, err := sourceSDK.ListKitchenRecipes(kitchenID)
		if err != nil {
			log.Fatal().Err(err).Msg("could not list recipes")
		}

		var recipeCount int
		for _, recipe := range recipesRes {
			recipeRes, err := sourceSDK.GetKitchenRecipe(kitchenID, recipe.RecipeID)
			if err != nil {
				log.Fatal().Err(err).Msg("could not get recipe")
			}

			_, err = targetSDK.AdminCreateRecipe(kitchenID, recipeRes)
			if err != nil {
				log.Error().Err(err).Msgf("could not sync recipe '%s'", recipe.RecipeID)
			}

			recipeCount++
			time.Sleep(100 * time.Millisecond)
		}

		log.Info().Int("recipes", recipeCount).Msgf("synced recipes for kitchen '%s'", kitchenID)

		// Sync Folders to database.
		// =========================================
		foldersRes, err := sourceSDK.ListKitchenFolders(kitchenID)
		if err != nil {
			log.Fatal().Err(err).Msg("could not list folders")
		}

		var folderCount int
		for _, folder := range foldersRes {
			_, err := folders.CreateFolder(context.TODO(), db, folders.CreateFolderInput{
				FolderID:  folder.FolderID,
				KitchenID: folder.KitchenID,
				Name:      folder.Name,
				Cover:     folder.Cover,
			})
			if err != nil {
				log.Fatal().Err(err).Msgf("could not sync folder '%s'", folder.FolderID)
			}

			folderCount++
		}
		log.Info().Int("folders", folderCount).Msgf("synced folders for kitchen '%s'", kitchenID)
	}

	// Sync Folder Recipes to Database
	// =========================================
	for _, kitchenID := range kitchenIDs {
		foldersRes, err := sourceSDK.ListKitchenFolders(kitchenID)
		if err != nil {
			log.Fatal().Err(err).Msg("could not list folders")
		}

		var folderCount int
		for _, folder := range foldersRes {
			folderRes, err := sourceSDK.GetKitchenFolder(kitchenID, folder.FolderID)
			if err != nil {
				log.Fatal().Err(err).Msgf("could not get folder '%s'", folder.FolderID)
			}

			var recipeIDs []string
			for _, recipe := range folderRes.Recipes {
				recipeIDs = append(recipeIDs, recipe.RecipeID)
			}

			if len(recipeIDs) > 0 {
				err = targetSDK.AdminAddFolderRecipes(kitchenID, folder.FolderID, recipeIDs)
				if err != nil {
					log.Error().Err(err).Msgf("could not create kitchen folder recipes for kitchen '%s' and folder '%s'", kitchenID, folder.FolderID)
				}
			}

			folderCount++
			time.Sleep(100 * time.Millisecond)
		}
		log.Info().Int("folders", folderCount).Msgf("synced folder recipes for kitchen '%s'", kitchenID)
	}

	// Sync Following
	// =========================================
	for _, kitchenID := range kitchenIDs {
		followed, err := sourceSDK.ListKitchenFollowing(kitchenID)
		if err != nil {
			log.Fatal().Err(err).Msg("could not list followed kitchens")
		}

		var followingCount int
		for _, following := range followed {
			err = targetSDK.AdminFollowKitchen(kitchenID, following.KitchenID)
			if err != nil {
				log.Error().Err(err).Msgf("could not follow kitchen '%s' from kitchen '%s'", following.KitchenID, kitchenID)
			}

			followingCount++
		}
		log.Info().Int("followings", followingCount).Msgf("synced followed kitchens for kitchen '%s'", kitchenID)
	}

	log.Info().Msg("sync complete")
}
