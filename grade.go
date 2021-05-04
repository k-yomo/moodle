package moodle

import "time"

type UserGrade struct {
	CourseID     int          `json:"courseid"`
	UserID       int          `json:"userid"`
	UserFullname string       `json:"userfullname"`
	MaxDepth     int          `json:"maxdepth"`
	GradeItems   []*GradeItem `json:"gradeitems"`
}

type GradeItem struct {
	ID                 int
	ItemName           string
	ItemType           string
	ItemModule         *string
	ItemInstance       int
	ItemNumber         *int
	CategoryID         *int
	OutcomeID          *int
	ScaleID            *int
	Locked             *bool
	CmID               int
	GradeRaw           float64
	GradeDateSubmitted *time.Time
	GradeDateGraded    *time.Time
	GradeHiddenByDate  bool
	GradeNeedsUpdate   bool
	GradeIsHidden      bool
	GradeIsLocked      *bool
	GradeIsOverridden  *bool
	GradeFormatted     string
	GradeMin           int
	GradeMax           int
	RangeFormatted     string
	Feedback           string
	FeedbackFormat     int
}
