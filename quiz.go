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
type GetQuizzesByCourseResponse struct {
	Quizzes    []*Quiz `json:"quizzes"`
	NextOffset int     `json:"nextoffset"`
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
