package moodle

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func Test_quizAPI_GetQuizzesByCourse(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx      context.Context
		courseID int
	}
	tests := []struct {
		name     string
		args     args
		response string
		want     []*Quiz
		wantErr  bool
	}{
		{
			name: "Successful response",
			args: args{ctx: context.Background(), courseID: 1111},
			response: `{
  "quizzes": [
    {
      "id": 2222,
      "course": 1111,
      "coursemodule": 123456,
      "name": "Quiz 1",
      "intro": "<p>This is a test quiz.<\/p>",
      "introformat": 1,
      "introfiles": [],
      "timeopen": 1577836800,
      "timeclose": 1590969600,
      "timelimit": 0,
      "preferredbehaviour": "deferredfeedback",
      "attempts": 0,
      "grademethod": 1,
      "decimalpoints": 2,
      "questiondecimalpoints": -1,
      "sumgrades": 5,
      "grade": 100,
      "hasfeedback": 0,
      "section": 1,
      "visible": 1,
      "groupmode": 1,
      "groupingid": 0
    }
  ],
  "warnings": []
}`,
			want: []*Quiz{
				{
					ID:                    2222,
					CourseID:              1111,
					CourseModuleID:        123456,
					Name:                  "Quiz 1",
					Intro:                 "<p>This is a test quiz.</p>",
					TimeOpen:              time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					TimeClose:             time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
					PreferredBehaviour:    "deferredfeedback",
					GradeMethod:           1,
					DecimalPoints:         2,
					QuestionDecimalPoints: -1,
					SumGrades:             5,
					Grade:                 1,
					Section:               1,
					Visible:               1,
					GroupMode:             1,
				},
			},
		},
		{
			name:     "Error response",
			args:     args{ctx: context.Background(), courseID: 0000},
			response: `{"errorcode": "invalidtoken"}`,
			wantErr:  true,
		},
		{
			name:     "Invalid json response",
			args:     args{ctx: context.Background(), courseID: 0000},
			response: "{",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q := mockQuizAPI(t, tt.response)
			got, err := q.GetQuizzesByCourse(tt.args.ctx, tt.args.courseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQuizzesByCourse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("GetQuizzesByCourse() (-got, +want)\n%s", diff)
			}
		})
	}
}

func Test_quizAPI_GetUserAttempts(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		quizID int
	}
	tests := []struct {
		name     string
		args     args
		response string
		want     []*QuizAttempt
		wantErr  bool
	}{
		{
			name: "Successful response",
			args: args{ctx: context.Background(), quizID: 1111},
			response: `{
  "attempts": [
    {
      "id": 2222,
      "quiz": 1111,
      "userid": 3333,
      "attempt": 1,
      "uniqueid": 123456,
      "layout": "1,2,3,4,5,0",
      "currentpage": 0,
      "preview": 0,
      "state": "finished",
      "timestart": 1577836800,
      "timefinish": 1577837100,
      "timemodified": 1577837400,
      "timemodifiedoffline": 1577837700,
      "timecheckstate": null,
      "sumgrades": 0
    }
  ],
  "warnings": []
}`,
			want: []*QuizAttempt{
				{
					ID:                  2222,
					QuizID:              1111,
					UserID:              3333,
					Attempt:             1,
					UniqueID:            123456,
					State:               "finished",
					TimeStart:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					TimeFinish:          func() *time.Time { t := time.Date(2020, 1, 1, 0, 5, 0, 0, time.UTC); return &t }(),
					TimeModified:        time.Date(2020, 1, 1, 0, 10, 0, 0, time.UTC),
					TimeModifiedOffline: time.Date(2020, 1, 1, 0, 15, 0, 0, time.UTC),
					TimeCheckState:      nil,
					SumGrades:           0,
				},
			},
		},
		{
			name:     "Error response",
			args:     args{ctx: context.Background(), quizID: 0000},
			response: `{"errorcode": "invalidtoken"}`,
			wantErr:  true,
		},
		{
			name:     "Invalid json response",
			args:     args{ctx: context.Background(), quizID: 0000},
			response: "{",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q := mockQuizAPI(t, tt.response)
			got, err := q.GetUserAttempts(tt.args.ctx, tt.args.quizID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserAttempts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("GetUserAttempts() (-got, +want)\n%s", diff)
			}
		})
	}
}

func Test_quizAPI_GetAttemptReview(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		attemptID int
	}
	tests := []struct {
		name     string
		args     args
		response string
		want     *QuizAttempt
		want1    []*QuizQuestion
		wantErr  bool
	}{
		{
			name: "Successful response",
			args: args{ctx: context.Background(), attemptID: 2222},
			response: `{
  "grade": 0,
  "attempt": {
    "id": 2222,
    "quiz": 1111,
    "userid": 3333,
    "attempt": 1,
    "uniqueid": 123456,
    "layout": "1,2,3,4,5,0",
    "currentpage": 0,
    "preview": 0,
    "state": "finished",
    "timestart": 1577836800,
    "timefinish": 1577837100,
    "timemodified": 1577837400,
    "timemodifiedoffline": 1577837700,
    "timecheckstate": null,
    "sumgrades": 0
  },
  "additionaldata": [],
  "questions": [
    {
      "slot": 1,
      "type": "multichoice",
      "page": 0,
      "html": "<div id=\"question-4494491-5\">question body</div>\n",
      "sequencecheck": 2,
      "lastactiontime": 1577836800,
      "hasautosavedstep": false,
      "flagged": false,
      "number": 5,
      "state": "gaveup",
      "status": "Not answered",
      "blockedbyprevious": false,
      "mark": "",
      "maxmark": 1
    }
  ],
  "warnings": []
}`,
			want: &QuizAttempt{
				ID:                  2222,
				QuizID:              1111,
				UserID:              3333,
				Attempt:             1,
				UniqueID:            123456,
				State:               "finished",
				TimeStart:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				TimeFinish:          func() *time.Time { t := time.Date(2020, 1, 1, 0, 5, 0, 0, time.UTC); return &t }(),
				TimeModified:        time.Date(2020, 1, 1, 0, 10, 0, 0, time.UTC),
				TimeModifiedOffline: time.Date(2020, 1, 1, 0, 15, 0, 0, time.UTC),
				TimeCheckState:      nil,
				SumGrades:           0,
			},
			want1: []*QuizQuestion{
				{
					Slot:           1,
					Type:           "multichoice",
					HtmlRaw:        "<div id=\"question-4494491-5\">question body</div>\n",
					SequenceCheck:  2,
					LastActionTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					Number:         5,
					State:          "gaveup",
					Status:         "Not answered",
					MaxMark:        1},
			},
		},
		{
			name:     "Error response",
			args:     args{ctx: context.Background(), attemptID: 0000},
			response: `{"errorcode": "invalidtoken"}`,
			wantErr:  true,
		},
		{
			name:     "Invalid json response",
			args:     args{ctx: context.Background(), attemptID: 0000},
			response: "{",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q := mockQuizAPI(t, tt.response)
			got, got1, err := q.GetAttemptReview(tt.args.ctx, tt.args.attemptID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAttemptReview() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("GetAttemptReview() (-got, +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("GetAttemptReview() (-got1, +want)\n%s", diff)
			}
		})
	}
}

func Test_quizAPI_StartAttempt(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		quizID int
	}
	tests := []struct {
		name     string
		args     args
		response string
		want     *QuizAttempt
		wantErr  bool
	}{
		{
			name: "Successful response",
			args: args{ctx: context.Background(), quizID: 1111},
			response: `{
  "attempt": {
    "id": 2222,
    "quiz": 1111,
    "userid": 3333,
    "attempt": 1,
    "uniqueid": 123456,
    "layout": "1,2,3,4,5,0",
    "currentpage": 0,
    "preview": 0,
    "state": "inprogress",
    "timestart": 1577836800,
    "timefinish": 0,
    "timemodified": 1577836800,
    "timemodifiedoffline": 1577836800,
    "timecheckstate": null,
    "sumgrades": 0
  },
  "warnings": []
}`,
			want: &QuizAttempt{
				ID:                  2222,
				QuizID:              1111,
				UserID:              3333,
				Attempt:             1,
				UniqueID:            123456,
				State:               "inprogress",
				TimeStart:           time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				TimeFinish:          nil,
				TimeModified:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				TimeModifiedOffline: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				TimeCheckState:      nil,
				SumGrades:           0,
			},
		},
		{
			name:     "Warning response",
			args:     args{ctx: context.Background(), quizID: 1111},
			response: `{"attempt":{},"warnings":[{"item":"quiz","itemid":1111,"warningcode":"1","message":"This quiz is not currently available"}]}`,
			wantErr:  true,
		},
		{
			name:     "Error response",
			args:     args{ctx: context.Background(), quizID: 0000},
			response: `{"exception":"dml_missing_record_exception","errorcode":"invalidrecord","message":"Can't find data record in database table quiz."}`,
			wantErr:  true,
		},
		{
			name:     "Invalid json response",
			args:     args{ctx: context.Background(), quizID: 0000},
			response: "{",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q := mockQuizAPI(t, tt.response)
			got, err := q.StartAttempt(tt.args.ctx, tt.args.quizID)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartAttempt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("StartAttempt() (-got, +want)\n%s", diff)
			}
		})
	}
}

func Test_quizAPI_FinishAttempt(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		attemptID int
		timeUp    bool
	}
	tests := []struct {
		name     string
		args     args
		response string
		wantErr  bool
	}{
		{
			name:     "Successful response",
			args:     args{ctx: context.Background(), attemptID: 1111, timeUp: false},
			response: `{"state":"finished","warnings":[]}`,
			wantErr:  false,
		},
		{
			name:     "Warning response",
			args:     args{ctx: context.Background(), attemptID: 1111, timeUp: false},
			response: `{"state":"inprogress","warnings":[{"item":"quiz","itemid":1111,"warningcode":"1","message":"Test message"}]}`,
			wantErr:  true,
		},
		{
			name:     "Error response",
			args:     args{ctx: context.Background(), attemptID: 0000},
			response: `{"exception":"moodle_quiz_exception","errorcode":"attemptalreadyclosed","message":"This attempt has already been finished."}`,
			wantErr:  true,
		},
		{
			name:     "Invalid json response",
			args:     args{ctx: context.Background(), attemptID: 0000},
			response: "{",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			q := mockQuizAPI(t, tt.response)
			if err := q.FinishAttempt(tt.args.ctx, tt.args.attemptID, tt.args.timeUp); (err != nil) != tt.wantErr {
				t.Errorf("FinishAttempt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func mockQuizAPI(t *testing.T, response string) *quizAPI {
	t.Helper()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	})
	s := httptest.NewServer(h)
	apiURL, _ := url.Parse(s.URL)
	return &quizAPI{
		&apiClient{
			httpClient: http.DefaultClient,
			apiURL:     apiURL,
		},
	}
}
