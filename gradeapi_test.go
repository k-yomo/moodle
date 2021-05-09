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

func Test_gradeAPI_GetGradesTable(t *testing.T) {
	type args struct {
		ctx      context.Context
		userID   int
		courseID int
	}
	tests := []struct {
		name     string
		args     args
		response string
		want     []*GradeTable
		wantErr  bool
	}{
		{
			name: "Successful response",
			args: args{ctx: context.Background(), courseID: 1111},
			response: `{
  "tables": [
    {
      "courseid": 1111,
      "userid": 222222,
      "userfullname": "Test User",
      "maxdepth": 3,
      "tabledata": [
        {
          "itemname": {
            "class": "level1 levelodd oddd1 b1b b1t column-itemname",
            "colspan": 6,
            "content": "<img class=\"icon icon itemicon\" alt=\"Category\" title=\"Category\" src=\"https:\/\/test.edu\/theme\/image.php\/lambda\/core\/1620139498\/i\/folder\" \/>Test Course",
            "celltype": "th",
            "id": "cat_25590_114374"
          },
          "leader": {
            "class": "level1 levelodd oddd1 b1t b2b b1l column-leader",
            "rowspan": 52
          }
        },
        {
          "itemname": {
            "class": "level2 leveleven oddd2 b1b b1t column-itemname",
            "colspan": 5,
            "content": "<img class=\"icon icon itemicon\" alt=\"Category\" title=\"Category\" src=\"https:\/\/test.edu\/theme\/image.php\/lambda\/core\/1620139498\/i\/folder\" \/>Assignments",
            "celltype": "th",
            "id": "cat_25595_114374"
          },
          "leader": {
            "class": "level2 leveleven oddd2 b1t b2b b1l column-leader",
            "rowspan": 14
          }
        },
        {
          "itemname": {
            "class": "level3 levelodd item b1b column-itemname",
            "colspan": 1,
            "content": "<a title=\"Assignment 1\" class=\"gradeitemheader\" href=\"https:\/\/test.edu\/mod\/workshop\/view.php?id=111111\"><img class=\"icon itemicon\" src=\"https:\/\/test.edu\/theme\/image.php\/lambda\/workshop\/1620139498\/icon\" alt=\"Workshop\" \/>Assignment 1<\/a>",
            "celltype": "th",
            "id": "row_179752_114374"
          },
          "grade": {
            "class": "level3 levelodd item b1b itemcenter  column-grade",
            "content": "92.00",
            "headers": "cat_25595_114374 row_179752_114374 grade"
          },
          "range": {
            "class": "level3 levelodd item b1b itemcenter  column-range",
            "content": "0&ndash;100",
            "headers": "cat_25595_114374 row_179752_114374 range"
          },
          "feedback": {
            "class": "level3 levelodd item b1b feedbacktext column-feedback",
            "content": "&nbsp;",
            "headers": "cat_25595_114374 row_179752_114374 feedback"
          },
          "contributiontocoursetotal": {
            "class": "level3 levelodd item b1b itemcenter  column-contributiontocoursetotal",
            "content": "2.70 %",
            "headers": "cat_25595_114374 row_179752_114374 contributiontocoursetotal"
          }
        },
        {
          "itemname": {
            "class": "level1 levelodd oddd1 baggt b2b column-itemname",
            "colspan": 2,
            "content": "<span class=\"gradeitemheader\" title=\"Course total\" tabindex=\"0\"><img class=\"icon icon itemicon\" alt=\"Weighted mean of grades\" title=\"Weighted mean of grades\" src=\"https:\/\/test.edu\/theme\/image.php\/lambda\/core\/1620139498\/i\/agg_mean\" \/>Course total<\/span><div class=\"gradeitemdescription\">Weighted mean of grades. Include empty grades.<\/div><div class=\"gradeitemdescriptionfiller\"><\/div>",
            "celltype": "th",
            "id": "row_179750_114374"
          },
          "grade": {
            "class": "level1 levelodd oddd1 baggt b2b itemcenter  column-grade",
            "content": "30.79",
            "headers": "cat_25590_114374 row_179750_114374 grade"
          },
          "range": {
            "class": "level1 levelodd oddd1 baggt b2b itemcenter  column-range",
            "content": "0&ndash;100",
            "headers": "cat_25590_114374 row_179750_114374 range"
          },
          "feedback": {
            "class": "level1 levelodd oddd1 baggt b2b feedbacktext column-feedback",
            "content": "&nbsp;",
            "headers": "cat_25590_114374 row_179750_114374 feedback"
          },
          "contributiontocoursetotal": {
            "class": "level1 levelodd oddd1 baggt b2b itemcenter  column-contributiontocoursetotal",
            "content": "-",
            "headers": "cat_25590_114374 row_179750_114374 contributiontocoursetotal"
          }
        }
      ]
    }
  ]
}`,
			want: []*GradeTable{
				{
					CourseID:     1111,
					UserID:       222222,
					UserFullname: "Test User",
					MaxDepth:     3,
					ItemGroups: []*GradeTableItemGroup{
						{
							Name: "Assignments",
							Items: []*GradeTableItem{
								{
									ItemName:                  "Assignment 1",
									ItemNameRawHTML:           `<a title="Assignment 1" class="gradeitemheader" href="https://test.edu/mod/workshop/view.php?id=111111"><img class="icon itemicon" src="https://test.edu/theme/image.php/lambda/workshop/1620139498/icon" alt="Workshop" />Assignment 1</a>`,
									ItemURL:                   func() *string { s := "https://test.edu/mod/workshop/view.php?id=111111"; return &s }(),
									IsGraded:                  true,
									Grade:                     92,
									GradeRangeMax:             100,
									Feedback:                  "\u00a0",
									FeedBackRawHTML:           "&nbsp;",
									ContributionToCourseTotal: 2.7,
								},
							},
						},
					},
				},
			},
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
		t.Run(tt.name, func(t *testing.T) {
			g := mockGradeAPI(t, tt.response)
			got, err := g.GetGradesTable(tt.args.ctx, tt.args.userID, tt.args.courseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGradesTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("GetGradesTable() (-got, +want)\n%s", diff)
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
		&apiClient{
			httpClient: http.DefaultClient,
			apiURL:     apiURL,
		},
	}
}
