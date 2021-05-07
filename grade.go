package moodle

import "time"

type UserGrade struct {
	CourseID     int
	UserID       int
	UserFullname string
	MaxDepth     int
	GradeItems   []*GradeItem
}

// GradeItem represents an grade
// If you want to know percentage of the grade out of course total grade, see GradeTableItem
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

// GradeTable represents a grade table for a course
type GradeTable struct {
	CourseID     int
	UserID       int
	UserFullname string
	MaxDepth     int
	ItemGroups   []*GradeTableItemGroup
}

// GradeTableItemGroup represents a group of grade items
type GradeTableItemGroup struct {
	Name  string
	Items []*GradeTableItem
}

type GradeTableItem struct {
	ItemName                  string
	ItemNameRawHTML           string
	ItemURL                   *string
	IsGraded                  bool
	Grade                     float64
	GradeRangeMin             float64
	GradeRangeMax             float64
	Feedback                  string
	FeedBackRawHTML           string
	ContributionToCourseTotal float64
}
