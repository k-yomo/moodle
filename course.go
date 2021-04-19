package moodle

import (
	"context"
	"fmt"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"net/http"
	"net/url"
)

type CourseClassification string

const (
	CourseClassificationPast       CourseClassification = "past"
	CourseClassificationInProgress CourseClassification = "inprogress"
	CourseClassificationFuture     CourseClassification = "future"
)

type Course struct {
	ID              int    `json:"id"`
	FullName        string `json:"fullname"`
	ShortName       string `json:"shortname"`
	IDNumber        string `json:"idnumber,omitempty"`
	Summary         string `json:",omitempty"`
	SummaryFormat   int    `json:"summaryformat"`
	StartDateUnix   int    `json:"startdate"`
	EndDateUnix     int    `json:"enddate"`
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
type GetEnrolledCoursesByTimelineClassificationResponse struct {
	Courses    []*Course `json:"courses"`
	NextOffset int       `json:"nextoffset"`
}

func getEnrolledCoursesByTimelineClassification(ctx context.Context, client *http.Client, apiURL *url.URL, classification CourseClassification) (*GetEnrolledCoursesByTimelineClassificationResponse, error) {
	u := urlutil.Copy(apiURL)
	urlutil.SetQueries(u, map[string]string{
		"wsfunction":     "core_course_get_enrolled_courses_by_timeline_classification",
		"classification": string(classification),
	})
	fmt.Println(u.String())
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	res := GetEnrolledCoursesByTimelineClassificationResponse{}
	if err := doAndMap(client, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
