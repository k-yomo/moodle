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

func Test_courseAPI_GetEnrolledCoursesByTimelineClassification(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx            context.Context
		classification CourseClassification
	}
	tests := []struct {
		name     string
		args     args
		response string
		want     []*Course
		wantErr  bool
	}{
		{
			name: "Successful response",
			args: args{ctx: context.Background(), classification: CourseClassificationInProgress},
			response: `{
  "courses": [
    {
      "id": 1111,
      "fullname": "MATH 1111 Introduction to Math",
      "shortname": "MATH 1111",
      "idnumber": "",
      "summary": "<p>This course presents students with basic concepts in mathematics<\/p>",
      "summaryformat": 1,
      "startdate": 1577836800,
      "enddate": 0,
      "visible": true,
      "fullnamedisplay": "MATH 1111 Introduction to MATH",
      "viewurl": "https:\/\/test.edu\/course\/view.php?id=1111",
      "courseimage": "https:\/\/test.edu\/pluginfile.php\/000000\/course\/overviewfiles\/MATH1111.jpg",
      "progress": 32,
      "hasprogress": true,
      "isfavourite": false,
      "hidden": false,
      "showshortname": false,
      "coursecategory": "Current Term"
    }
  ],
  "nextoffset": 2
}`,
			want: []*Course{
				{
					ID:              1111,
					FullName:        "MATH 1111 Introduction to Math",
					ShortName:       "MATH 1111",
					Summary:         "<p>This course presents students with basic concepts in mathematics</p>",
					SummaryFormat:   1,
					StartDate:       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					Visible:         true,
					FullNameDisplay: "MATH 1111 Introduction to Math",
					ViewURL:         "https://test.edu/course/view.php?id=1111",
					CourseImage:     "https://test.edu/pluginfile.php/000000/course/overviewfiles/MATH1111.jpg",
					Progress:        32,
					HasProgress:     true,
					CourseCategory:  "Current Term",
				},
			},
		},
		{
			name:     "Error response",
			args:     args{ctx: context.Background(), classification: CourseClassificationInProgress},
			response: `{"errorcode": "invalidtoken"}`,
			wantErr:  true,
		},
		{
			name:     "Invalid json response",
			args:     args{ctx: context.Background(), classification: CourseClassificationInProgress},
			response: "{",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, tt.response)
			})
			s := httptest.NewServer(h)
			apiURL, _ := url.Parse(s.URL)
			c := &courseAPI{
				&apiClient{
					httpClient: http.DefaultClient,
					apiURL:     apiURL,
				},
			}
			got, err := c.GetEnrolledCoursesByTimelineClassification(tt.args.ctx, tt.args.classification)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnrolledCoursesByTimelineClassification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("GetEnrolledCoursesByTimelineClassification() (-got, +want)\n%s", diff)
			}
		})
	}
}
