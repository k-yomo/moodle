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

func Test_gradeAPI_GetGradeItems(t *testing.T) {
	type args struct {
		ctx      context.Context
		userID   int
		courseID int
	}
	tests := []struct {
		name     string
		args     args
		response string
		want     []*UserGrade
		wantErr  bool
	}{
		{
			name: "Successful response",
			args: args{ctx: context.Background(), courseID: 1111},
			response: `{
  "usergrades": [
    {
      "courseid": 1111,
      "userid": 123456,
      "userfullname": "Test User",
      "maxdepth": 3,
      "gradeitems": [
        {
          "id": 121212,
          "itemname": "Test Quiz",
          "itemtype": "mod",
          "itemmodule": "workshop",
          "iteminstance": 333333,
          "itemnumber": 0,
          "categoryid": 444444,
          "outcomeid": null,
          "scaleid": null,
          "locked": null,
          "cmid": 555555,
          "graderaw": 81,
          "gradedatesubmitted": 1577836800,
          "gradedategraded": 1577837100,
          "gradehiddenbydate": false,
          "gradeneedsupdate": false,
          "gradeishidden": false,
          "gradeislocked": null,
          "gradeisoverridden": null,
          "gradeformatted": "81.00",
          "grademin": 0,
          "grademax": 90,
          "rangeformatted": "0&ndash;90",
          "feedback": "",
          "feedbackformat": 1
        }
      ]
	}
  ]
}`,
			want: []*UserGrade{{
				CourseID:     1111,
				UserID:       123456,
				UserFullname: "Test User",
				MaxDepth:     3,
				GradeItems: []*GradeItem{
					{
						ID:                 121212,
						ItemName:           "Test Quiz",
						ItemType:           "mod",
						ItemModule:         func() *string { s := "workshop"; return &s }(),
						ItemInstance:       333333,
						ItemNumber:         func() *int { i := 0; return &i }(),
						CategoryID:         func() *int { i := 444444; return &i }(),
						CmID:               555555,
						GradeRaw:           81,
						GradeDateSubmitted: func() *time.Time { t := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC); return &t }(),
						GradeDateGraded:    func() *time.Time { t := time.Date(2020, 1, 1, 0, 5, 0, 0, time.UTC); return &t }(),
						GradeHiddenByDate:  false,
						GradeNeedsUpdate:   false,
						GradeIsHidden:      false,
						GradeIsLocked:      nil,
						GradeIsOverridden:  nil,
						GradeFormatted:     "81.00",
						GradeMin:           0,
						GradeMax:           90,
						RangeFormatted:     "0&ndash;90",
						Feedback:           "",
						FeedbackFormat:     1,
					},
				},
			}},
		},
		{
			name:     "Warning response",
			args:     args{ctx: context.Background(), courseID: 1111},
			response: `{"usergrades":[],"warnings":[{"item":"grade","itemid":1111,"warningcode":"1","message":"test warning"}]}`,
			wantErr:  true,
		},
		{
			name:     "Error response",
			args:     args{ctx: context.Background(), courseID: 0000},
			response: `{"exception":"dml_missing_record_exception","errorcode":"invalidrecord","message":"Can't find data record in database table course."}`,
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

			g := mockGradeAPI(t, tt.response)
			got, err := g.GetGradeItems(tt.args.ctx, tt.args.userID, tt.args.courseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGradeItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("GetGradeItems() (-got, +want)\n%s", diff)
			}
		})
	}
}

func mockGradeAPI(t *testing.T, response string) *gradeAPI {
	t.Helper()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	})
	s := httptest.NewServer(h)
	apiURL, _ := url.Parse(s.URL)
	return &gradeAPI{
		httpClient: &http.Client{},
		apiURL:     apiURL,
	}
}
