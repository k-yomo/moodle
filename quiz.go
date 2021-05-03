package moodle

import (
	"time"
)

type Quiz struct {
	ID                    int
	CourseID              int
	CourseModuleID        int
	Name                  string
	Intro                 string
	IntroFormat           int
	TimeOpen              time.Time
	TimeClose             time.Time
	TimeLimit             int
	PreferredBehaviour    string
	Attempts              int
	GradeMethod           int
	DecimalPoints         int
	QuestionDecimalPoints int
	SumGrades             int
	Grade                 int
	HasFeedback           int
	Section               int
	Visible               int
	GroupMode             int
	GroupingID            int
}

type QuizAttempt struct {
	ID                  int
	QuizID              int
	UserID              int
	Attempt             int
	UniqueID            int
	Layout              string
	CurrentPage         int
	Preview             int
	State               string
	TimeStart           time.Time
	TimeFinish          *time.Time
	TimeModified        time.Time
	TimeModifiedOffline time.Time
	TimeCheckState      *time.Time
	SumGrades           int
}

type QuizQuestion struct {
	Slot              int
	Type              string
	Page              int
	HtmlRaw           string
	SequenceCheck     int
	LastActionTime    time.Time
	HasAutoSavedStep  bool
	Flagged           bool
	Number            int
	State             string
	Status            string
	BlockedByPrevious bool
	Mark              string
	MaxMark           int
}
