package entity

import "time"

type FeatureDetail struct {
	ID          int
	FeatureID   int
	PathOfFile  string
	Title       string
	Code        string
	Description string
	AuthorName  string
	FeatureName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
