package moodle

import (
	"context"
	"github.com/k-yomo/moodle/pkg/maputil"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"net/http"
	"strconv"
)

type quizResponse struct {
	ID                    int    `json:"id"`
	CourseID              int    `json:"course"`
	CourseModuleID        int    `json:"coursemodule"`
	Name                  string `json:"name"`
	Intro                 string `json:"intro"`
	IntroFormat           int    `json:"introfomat"`
	TimeOpenUnix          int64  `json:"timeopen"`
	TimeCloseUnix         int64  `json:"timeclose"`
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

type quizAttemptResponse struct {
	ID                      int    `json:"id"`
	QuizID                  int    `json:"quiz"`
	UserID                  int    `json:"userid"`
	Attempt                 int    `json:"attempt"`
	UniqueID                int    `json:"uniqueid"`
	Layout                  string `json:"lauout"`
	CurrentPage             int    `json:"currentpage"`
	Preview                 int    `json:"preview"`
	State                   string `json:"state"`
	TimeStartUnix           int64  `json:"timestart"`
	TimeFinishUnix          int64  `json:"timefinish"`
	TimeModifiedUnix        int64  `json:"timemodified"`
	TimeModifiedOfflineUnix int64  `json:"timemodifiedoffline"`
	TimeCheckStateUnix      *int64 `json:"timecheckstate,omitempty"`
	SumGrades               int    `json:"sumgrades"`
}

type quizQuestionResponse struct {
	Slot              int    `json:"slot"`
	Type              string `json:"type"`
	Page              int    `json:"page"`
	Html              string `json:"html"`
	SequenceCheck     int    `json:"sequencecheck"`
	LastActionTime    int    `json:"lastactiontime"`
	HasAutoSavedStep  bool   `json:"hasautosavedstep"`
	Flagged           bool   `json:"flagged"`
	Number            int    `json:"number"`
	State             string `json:"state"`
	Status            string `json:"status"`
	BlockedByPrevious bool   `json:"blockedbyprevious"`
	Mark              string `json:"mark"`
	MaxMark           int    `json:"maxmark"`
}

type getQuizzesByCourseResponse struct {
	Quizzes []*quizResponse `json:"quizzes"`
}

func (q *quizAPI) getQuizzesByCourse(ctx context.Context, courseID int) (*getQuizzesByCourseResponse, error) {
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

	res := getQuizzesByCourseResponse{}
	if err := doAndMap(q.httpClient, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type getUserAttemptsResponse struct {
	Attempts []*quizAttemptResponse `json:"attempts"`
}

func (q *quizAPI) getUserAttempts(ctx context.Context, quizID int) (*getUserAttemptsResponse, error) {
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

	res := getUserAttemptsResponse{}
	if err := doAndMap(q.httpClient, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type getAttemptReviewResponse struct {
	Grade     int                     `json:"grade"`
	Attempt   *quizAttemptResponse    `json:"attempt"`
	Questions []*quizQuestionResponse `json:"questions"`
}

func (q *quizAPI) getAttemptReview(ctx context.Context, attemptID int) (*getAttemptReviewResponse, error) {
	u := urlutil.CopyWithQueries(
		q.apiURL,
		map[string]string{
			"wsfunction": "mod_quiz_get_attempt_review",
			"attemptid":  strconv.Itoa(attemptID),
		},
	)

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	res := getAttemptReviewResponse{}
	if err := doAndMap(q.httpClient, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
