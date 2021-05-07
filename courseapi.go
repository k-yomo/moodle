package moodle

import (
	"context"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"net/http"
	"net/url"
	"time"
)

type CourseAPI interface {
	GetEnrolledCoursesByTimelineClassification(ctx context.Context, classification CourseClassification) ([]*Course, error)
}

type courseAPI struct {
	httpClient *http.Client
	apiURL     *url.URL
}

func newCourseAPI(httpClient *http.Client, apiURL *url.URL) *courseAPI {
	return &courseAPI{
		httpClient: httpClient,
		apiURL:     apiURL,
	}
}

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

func (c *courseAPI) GetEnrolledCoursesByTimelineClassification(ctx context.Context, classification CourseClassification) ([]*Course, error) {
	u := urlutil.CopyWithQueries(c.apiURL, map[string]string{
		"wsfunction":     "core_course_get_enrolled_courses_by_timeline_classification",
		"classification": string(classification),
	})
	res := getEnrolledCoursesByTimelineClassificationResponse{}
	if err := getAndUnmarshal(ctx, c.httpClient, u, &res); err != nil {
		return nil, err
	}

	return mapToCourseList(res.Courses), nil
}

func mapToCourseList(courseResList []*courseResponse) []*Course {
	courses := make([]*Course, 0, len(courseResList))
	for _, courseRes := range courseResList {
		courses = append(courses, mapToCourse(courseRes))
	}
	return courses
}

func mapToCourse(courseRes *courseResponse) *Course {
	return &Course{
		ID:              courseRes.ID,
		FullName:        courseRes.FullName,
		ShortName:       courseRes.ShortName,
		Summary:         courseRes.Summary,
		SummaryFormat:   courseRes.SummaryFormat,
		StartDate:       time.Unix(courseRes.StartDateUnix, 0),
		EndDate:         time.Unix(courseRes.StartDateUnix, 0),
		Visible:         courseRes.Visible,
		FullNameDisplay: courseRes.FullName,
		ViewURL:         courseRes.ViewURL,
		CourseImage:     courseRes.CourseImage,
		Progress:        courseRes.Progress,
		HasProgress:     courseRes.HasProgress,
		IsSavourite:     courseRes.IsSavourite,
		Hidden:          courseRes.Hidden,
		ShowShortName:   courseRes.ShowShortName,
		CourseCategory:  courseRes.CourseCategory,
	}
}
