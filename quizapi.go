package moodle

import (
	"context"
	"strconv"
	"time"
)

type QuizAPI interface {
	GetQuizzesByCourse(ctx context.Context, courseID int) ([]*Quiz, error)
	GetUserAttempts(ctx context.Context, quizID int) ([]*QuizAttempt, error)
	GetAttemptReview(ctx context.Context, attemptID int) (*QuizAttempt, []*QuizQuestion, error)
	StartAttempt(ctx context.Context, quizID int) (*QuizAttempt, error)
	FinishAttempt(ctx context.Context, attemptID int, timeUp bool) error
}

type quizAPI struct {
	*apiClient
}

func newQuizAPI(apiClient *apiClient) *quizAPI {
	return &quizAPI{apiClient}
}

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
	Slot               int    `json:"slot"`
	Type               string `json:"type"`
	Page               int    `json:"page"`
	Html               string `json:"html"`
	SequenceCheck      int    `json:"sequencecheck"`
	LastActionTimeUnix int64  `json:"lastactiontime"`
	HasAutoSavedStep   bool   `json:"hasautosavedstep"`
	Flagged            bool   `json:"flagged"`
	Number             int    `json:"number"`
	State              string `json:"state"`
	Status             string `json:"status"`
	BlockedByPrevious  bool   `json:"blockedbyprevious"`
	Mark               string `json:"mark"`
	MaxMark            int    `json:"maxmark"`
}

type getQuizzesByCourseResponse struct {
	Quizzes []*quizResponse `json:"quizzes"`
}

func (q *quizAPI) GetQuizzesByCourse(ctx context.Context, courseID int) ([]*Quiz, error) {
	res := getQuizzesByCourseResponse{}

	err := q.callMoodleFunction(
		ctx,
		&res,
		map[string]string{"wsfunction": "mod_quiz_get_quizzes_by_courses"},
		mapStrArrayToQueryParams("courseids", []string{strconv.Itoa(courseID)}),
	)
	if err != nil {
		return nil, err
	}
	return mapToQuizList(res.Quizzes), nil
}

type getUserAttemptsResponse struct {
	Attempts []*quizAttemptResponse `json:"attempts"`
}

func (q *quizAPI) GetUserAttempts(ctx context.Context, quizID int) ([]*QuizAttempt, error) {
	res := getUserAttemptsResponse{}
	err := q.callMoodleFunction(
		ctx,
		&res,
		map[string]string{
			"wsfunction": "mod_quiz_get_user_attempts",
			"quizid":     strconv.Itoa(quizID),
		},
	)
	if err != nil {
		return nil, err
	}
	return mapToQuizAttemptList(res.Attempts), nil
}

type getAttemptReviewResponse struct {
	Grade     int                     `json:"grade"`
	Attempt   *quizAttemptResponse    `json:"attempt"`
	Questions []*quizQuestionResponse `json:"questions"`
}

func (q *quizAPI) GetAttemptReview(ctx context.Context, attemptID int) (*QuizAttempt, []*QuizQuestion, error) {
	res := getAttemptReviewResponse{}
	err := q.callMoodleFunction(
		ctx,
		&res,
		map[string]string{
			"wsfunction": "mod_quiz_get_attempt_review",
			"attemptid":  strconv.Itoa(attemptID),
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return mapToQuizAttempt(res.Attempt), mapToQuizQuestionList(res.Questions), nil
}

type startAttemptResponse struct {
	Attempt  *quizAttemptResponse `json:"attempt,omitempty"`
	Warnings Warnings             `json:"warnings,omitempty"`
}

func (q *quizAPI) StartAttempt(ctx context.Context, quizID int) (*QuizAttempt, error) {
	res := startAttemptResponse{}
	err := q.callMoodleFunction(
		ctx,
		&res,
		map[string]string{
			"wsfunction": "mod_quiz_start_attempt",
			"quizid":     strconv.Itoa(quizID),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(res.Warnings) > 0 {
		return nil, res.Warnings
	}
	return mapToQuizAttempt(res.Attempt), nil
}

type finishAttemptResponse struct {
	State    string   `json:"state,omitempty"`
	Warnings Warnings `json:"warnings,omitempty"`
}

func (q *quizAPI) FinishAttempt(ctx context.Context, attemptID int, timeUp bool) error {
	res := finishAttemptResponse{}
	err := q.callMoodleFunction(
		ctx,
		&res,
		map[string]string{
			"wsfunction":    "mod_quiz_process_attempt",
			"attemptid":     strconv.Itoa(attemptID),
			"finishattempt": "1",
			"timeup":        mapBoolToBitStr(timeUp),
		},
	)
	if err != nil {
		return err
	}
	if len(res.Warnings) > 0 {
		return res.Warnings
	}
	return nil
}

func mapToQuizList(quizResList []*quizResponse) []*Quiz {
	quizzes := make([]*Quiz, 0, len(quizResList))
	for _, quizRes := range quizResList {
		quizzes = append(quizzes, mapToQuiz(quizRes))
	}
	return quizzes
}

func mapToQuiz(quizRes *quizResponse) *Quiz {
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

func mapToQuizAttemptList(attemptResList []*quizAttemptResponse) []*QuizAttempt {
	attempts := make([]*QuizAttempt, 0, len(attemptResList))
	for _, attemptRes := range attemptResList {
		attempts = append(attempts, mapToQuizAttempt(attemptRes))
	}
	return attempts
}

func mapToQuizAttempt(attemptRes *quizAttemptResponse) *QuizAttempt {
	var timeFinish, timeCheckState *time.Time
	if attemptRes.TimeFinishUnix > 0 {
		t := time.Unix(attemptRes.TimeFinishUnix, 0)
		timeFinish = &t
	}
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
		TimeFinish:          timeFinish,
		TimeModified:        time.Unix(attemptRes.TimeModifiedUnix, 0),
		TimeModifiedOffline: time.Unix(attemptRes.TimeModifiedOfflineUnix, 0),
		TimeCheckState:      timeCheckState,
		SumGrades:           attemptRes.SumGrades,
	}
}

func mapToQuizQuestionList(quizQuestionResList []*quizQuestionResponse) []*QuizQuestion {
	questions := make([]*QuizQuestion, 0, len(quizQuestionResList))
	for _, questionRes := range quizQuestionResList {
		questions = append(questions, mapToQuizQuestion(questionRes))
	}
	return questions
}

func mapToQuizQuestion(quizQuestionRes *quizQuestionResponse) *QuizQuestion {
	return &QuizQuestion{
		Slot:              quizQuestionRes.Slot,
		Type:              quizQuestionRes.Type,
		Page:              quizQuestionRes.Page,
		HtmlRaw:           quizQuestionRes.Html,
		SequenceCheck:     quizQuestionRes.SequenceCheck,
		LastActionTime:    time.Unix(quizQuestionRes.LastActionTimeUnix, 0),
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
