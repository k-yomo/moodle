package moodle

import (
	"context"
	"net/http"
	"net/url"
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
	TimeFinish          time.Time
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
	LastActionTime    int
	HasAutoSavedStep  bool
	Flagged           bool
	Number            int
	State             string
	Status            string
	BlockedByPrevious bool
	Mark              string
	MaxMark           int
}

type QuizAPI interface {
	GetQuizzesByCourse(ctx context.Context, courseID int) ([]*Quiz, error)
	GetUserAttempts(ctx context.Context, quizID int) ([]*QuizAttempt, error)
	GetAttemptReview(ctx context.Context, attemptID int) (*QuizAttempt, []*QuizQuestion, error)
}

type quizAPI struct {
	httpClient *http.Client
	apiURL     *url.URL
}

func newQuizAPI(httpClient *http.Client, apiURL *url.URL) *quizAPI {
	return &quizAPI{
		httpClient: httpClient,
		apiURL:     apiURL,
	}
}

func (q *quizAPI) GetQuizzesByCourse(ctx context.Context, courseID int) ([]*Quiz, error) {
	res, err := q.getQuizzesByCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}
	return mapFromQuizListResponse(res.Quizzes), nil
}

func (q *quizAPI) GetUserAttempts(ctx context.Context, quizID int) ([]*QuizAttempt, error) {
	res, err := q.getUserAttempts(ctx, quizID)
	if err != nil {
		return nil, err
	}
	return mapFromQuizAttemptListResponse(res.Attempts), nil
}

func (q *quizAPI) GetAttemptReview(ctx context.Context, attemptID int) (*QuizAttempt, []*QuizQuestion, error) {
	res, err := q.getAttemptReview(ctx, attemptID)
	if err != nil {
		return nil, nil, err
	}
	return mapFromQuizAttemptResponse(res.Attempt), mapFromQuizQuestionListResponse(res.Questions), nil
}

func mapFromQuizListResponse(quizResList []*quizResponse) []*Quiz {
	quizzes := make([]*Quiz, len(quizResList))
	for _, quizRes := range quizResList {
		quizzes = append(quizzes, mapFromQuizResponse(quizRes))
	}
	return quizzes
}

func mapFromQuizResponse(quizRes *quizResponse) *Quiz {
	return &Quiz{
		ID:                    quizRes.ID,
		CourseID:              quizRes.CourseID,
		CourseModuleID:        quizRes.CourseModuleID,
		Name:                  quizRes.Name,
		Intro:                 quizRes.Intro,
		IntroFormat:           quizRes.IntroFormat,
		TimeOpen:              time.Unix(quizRes.TimeOpenUnix, 0),
		TimeClose:             time.Unix(quizRes.TimeCloseUnix, 0),
		TimeLimit:             quizRes.TimeLimit,
		PreferredBehaviour:    quizRes.PreferredBehaviour,
		Attempts:              quizRes.Attempts,
		GradeMethod:           quizRes.GradeMethod,
		DecimalPoints:         quizRes.DecimalPoints,
		QuestionDecimalPoints: quizRes.QuestionDecimalPoints,
		SumGrades:             quizRes.SumGrades,
		Grade:                 quizRes.GradeMethod,
		HasFeedback:           quizRes.HasFeedback,
		Section:               quizRes.Section,
		Visible:               quizRes.Visible,
		GroupMode:             quizRes.GroupMode,
		GroupingID:            quizRes.GroupingID,
	}
}

func mapFromQuizAttemptListResponse(attemptResList []*quizAttemptResponse) []*QuizAttempt {
	attempts := make([]*QuizAttempt, 0, len(attemptResList))
	for _, attemptRes := range attemptResList {
		attempts = append(attempts, mapFromQuizAttemptResponse(attemptRes))
	}
	return attempts
}

func mapFromQuizAttemptResponse(attemptRes *quizAttemptResponse) *QuizAttempt {
	var timeCheckState *time.Time
	if attemptRes.TimeCheckStateUnix != nil {
		t := time.Unix(*attemptRes.TimeCheckStateUnix, 0)
		timeCheckState = &t
	}
	return &QuizAttempt{
		ID:                  attemptRes.ID,
		QuizID:              attemptRes.QuizID,
		UserID:              attemptRes.UserID,
		Attempt:             attemptRes.Attempt,
		UniqueID:            attemptRes.UniqueID,
		Layout:              attemptRes.Layout,
		CurrentPage:         attemptRes.CurrentPage,
		Preview:             attemptRes.Preview,
		State:               attemptRes.State,
		TimeStart:           time.Unix(attemptRes.TimeStartUnix, 0),
		TimeFinish:          time.Unix(attemptRes.TimeFinishUnix, 0),
		TimeModified:        time.Unix(attemptRes.TimeModifiedUnix, 0),
		TimeModifiedOffline: time.Unix(attemptRes.TimeModifiedOfflineUnix, 0),
		TimeCheckState:      timeCheckState,
		SumGrades:           attemptRes.SumGrades,
	}
}

func mapFromQuizQuestionListResponse(quizQuestionResList []*quizQuestionResponse) []*QuizQuestion {
	questions := make([]*QuizQuestion, 0, len(quizQuestionResList))
	for _, questionRes := range quizQuestionResList {
		questions = append(questions, mapFromQuizQuestionResponse(questionRes))
	}
	return questions
}

func mapFromQuizQuestionResponse(quizQuestionRes *quizQuestionResponse) *QuizQuestion {
	return &QuizQuestion{
		Slot:              quizQuestionRes.Slot,
		Type:              quizQuestionRes.Type,
		Page:              quizQuestionRes.Page,
		HtmlRaw:           quizQuestionRes.Html,
		SequenceCheck:     quizQuestionRes.SequenceCheck,
		LastActionTime:    quizQuestionRes.LastActionTime,
		HasAutoSavedStep:  quizQuestionRes.HasAutoSavedStep,
		Flagged:           quizQuestionRes.Flagged,
		Number:            quizQuestionRes.Number,
		State:             quizQuestionRes.State,
		Status:            quizQuestionRes.Status,
		BlockedByPrevious: quizQuestionRes.BlockedByPrevious,
		Mark:              quizQuestionRes.Mark,
		MaxMark:           quizQuestionRes.MaxMark,
	}
}
