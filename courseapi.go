package moodle

import (
	"context"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"net/http"
)

type courseResponse struct {
	ID              int    `json:"id"`
	FullName        string `json:"fullname"`
	ShortName       string `json:"shortname"`
	Summary         string `json:",omitempty"`
	SummaryFormat   int    `json:"summaryformat"`
	StartDateUnix   int64  `json:"startdate"`
	EndDateUnix     int64  `json:"enddate"`
	Visible         bool   `json:"visible"`
	FullNameDisplay string `json:"fullnamedisplay"`
	ViewURL         string `json:"viewurl"`
	CourseImage     string `json:"courseimage"`
	Progress        int    `json:"progress"`
	HasProgress     bool   `json:"hasprogress"`
	IsSavourite     bool   `json:"isfavourite"`
	Hidden          bool   `json:"hidden"`
	ShowShortName   bool   `json:"showshortname"`
	CourseCategory  string `json:"coursecategory"`
}

type getEnrolledCoursesByTimelineClassificationResponse struct {
	Courses    []*courseResponse `json:"courses"`
	NextOffset int               `json:"nextoffset"`
}

func (c *courseAPI) getEnrolledCoursesByTimelineClassification(ctx context.Context, classification CourseClassification) (*getEnrolledCoursesByTimelineClassificationResponse, error) {
	u := urlutil.CopyWithQueries(c.apiURL, map[string]string{
		"wsfunction":     "core_course_get_enrolled_courses_by_timeline_classification",
		"classification": string(classification),
	})
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	res := getEnrolledCoursesByTimelineClassificationResponse{}
	if err := doAndMap(c.httpClient, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
