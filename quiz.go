package moodle

import (
	"context"
	"github.com/k-yomo/moodle/pkg/maputil"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"net/http"
	"net/url"
	"strconv"
)

type QuizAPI interface {
	GetQuizzesByCourse(ctx context.Context, courseID int) ([]*Quiz, error)
	GetUserAttempts(ctx context.Context, quizID int) ([]*QuizAttempt, error)
	GetAttemptReview(ctx context.Context, attemptID int) error
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

type Quiz struct {
	ID                    int    `json:"id"`
	CourseID              int    `json:"course"`
	CourseModuleID        int    `json:"coursemodule"`
	Name                  string `json:"name"`
	Intro                 string `json:"intro"`
	IntroFormat           int    `json:"introfomat"`
	TimeOpen              int    `json:"timeopen"`
	TimeClose             int    `json:"timeclose"`
	TimeLimit             int    `json:"timelimit"`
	PreferredBehaviour    string `json:"preferredbehaviour"`
	Attempts              int    `json:"attempts"`
	GradeMethod           int    `json:"grademethod"`
	DecimalPoints         int    `json:"decimalpoints"`
	QuestionDecimalPoints int    `json:"questiondecimalpoints"`
	SumGrades             int    `json:"sumgrades"`
	Grade                 int    `json:"grade"`
	HasFeedback           int    `json:"hasfeedback"`
	Section               int    `json:"section"`
	Visible               int    `json:"visible"`
	GroupMode             int    `json:"groupmode"`
	GroupingID            int    `json:"groupingid"`
}

type QuizAttempt struct {
	ID                  int    `json:"id"`
	QuizID              int    `json:"quiz"`
	UserID              int    `json:"userid"`
	Attempt             int    `json:"attempt"`
	UniqueID            int    `json:"uniqueid"`
	Layout              string `json:"lauout"`
	CurrentPage         int    `json:"currentpage"`
	Preview             int    `json:"preview"`
	State               string `json:"state"`
	TimeStart           int    `json:"timestart"`
	TimeFinish          int    `json:"timefinish"`
	TimeModified        int    `json:"timemodified"`
	TimeModifiedOffline int    `json:"timemodifiedoffline"`
	timecheckstate      *int   `json:"timecheckstate,omitempty"`
	SumGrades           int    `json:"sumgrades"`
}

type QuizQuestion struct {
	Slot              int    `json:"slot"`
	Type              string `json:"type"`
	Page              int    `json:"page"`
	Html              string `json:"string"`
	SequenceCheck     int    `json:"sequencecheck"`
	LastActionTime    int    `json:"lastactiontime"`
	HasAutoSavedStep  bool   `json:"hasautosavedstep"`
	Flagged           int    `json:"flagged"`
	Number            int    `json:"number"`
	State             string `json:"state"`
	Status            string `json:"status"`
	BlockedByPrevious bool   `json:"blockedbyprevious"`
	Mark              string `json:"mark"`
	MaxMark           int    `json:"maxmark"`
}

type GetQuizzesByCourseResponse struct {
	Quizzes []*Quiz `json:"quizzes"`
}

func (q *quizAPI) GetQuizzesByCourse(ctx context.Context, courseID int) ([]*Quiz, error) {
	u := urlutil.CopyWithQueries(
		q.apiURL,
		maputil.MergeStrMap(
			map[string]string{"wsfunction": "mod_quiz_get_quizzes_by_courses"},
			strArrayToQueryParams("courseids", []string{strconv.Itoa(courseID)}),
		),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	res := GetQuizzesByCourseResponse{}
	if err := doAndMap(q.httpClient, req, &res); err != nil {
		return nil, err
	}
	return res.Quizzes, nil
}

type GetUserAttemptsResponse struct {
	Attempts []*QuizAttempt `json:"attempts"`
}

func (q *quizAPI) GetUserAttempts(ctx context.Context, quizID int) ([]*QuizAttempt, error) {
	u := urlutil.CopyWithQueries(
		q.apiURL,
		map[string]string{
			"wsfunction": "mod_quiz_get_user_attempts",
			"quizid":     strconv.Itoa(quizID),
		},
	)

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	res := GetUserAttemptsResponse{}
	if err := doAndMap(q.httpClient, req, &res); err != nil {
		return nil, err
	}
	return res.Attempts, nil
}

type GetAttemptReviewResponse struct {
	Grade   int          `json:"grade"`
	Attempt *QuizAttempt `json:"attempt"`
}

func (q *quizAPI) GetAttemptReview(ctx context.Context, attemptID int) error {
	u := urlutil.CopyWithQueries(
		q.apiURL,
		map[string]string{
			"wsfunction": "mod_quiz_get_attempt_review",
			"attemptid":  strconv.Itoa(attemptID),
		},
	)

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}

	res := GetAttemptReviewResponse{}
	if err := doAndMap(q.httpClient, req, &res); err != nil {
		return err
	}
	return nil
}
