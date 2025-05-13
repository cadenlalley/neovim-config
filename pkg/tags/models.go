package tags

import "errors"

const (
	TagTypeDifficulty = "difficulty"
	TagTypeCourse     = "course"
	TagTypeClass      = "class"
	TagTypeCuisine    = "cuisine"
	TagTypeDiet       = "diet"
	TagTypeIngredient = "ingredient"
	TagTypeEquipment  = "equipment"
	TagTypeKeyword    = "keyword"
)

type Tag struct {
	TagID int    `json:"tagId" db:"tag_id"`
	Type  string `json:"type" db:"tag_type"`
	Value string `json:"value" db:"tag_value"`
}

func ValidateTagType(t string) error {
	switch t {
	case TagTypeDifficulty,
		TagTypeCourse,
		TagTypeClass,
		TagTypeCuisine,
		TagTypeDiet,
		TagTypeIngredient,
		TagTypeEquipment,
		TagTypeKeyword:
		return nil
	default:
		return errors.New("invalid tag type")
	}
}
