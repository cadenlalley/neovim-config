package models

import (
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"gopkg.in/guregu/null.v4"
)

type Profile struct {
	ProfileID   string    `json:"profile_id" db:"profile_id"`
	AccountID   string    `json:"account_id" db:"account_id"`
	Name        string    `json:"name" db:"name"`
	Bio         string    `json:"bio" db:"bio"`
	Handle      string    `json:"handle" db:"handle"`
	AvatarPhoto string    `json:"avatarPhoto" db:"avatar_photo"`
	CoverPhoto  string    `json:"coverPhoto" db:"cover_photo"`
	Public      bool      `json:"public" db:"public"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt   null.Time `json:"deletedAt" db:"deleted_at"`
}

func CreateProfileID() string {
	return fmt.Sprintf("pro_%s", ksuid.New().String())
}
