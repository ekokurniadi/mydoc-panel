package entity

import "time"

type Feature struct {
	ID                 int
	FeatureName        string
	FeatureDescription string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
