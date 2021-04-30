package moodle

import (
	"context"
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

type CourseClassification string

const (
	CourseClassificationPast       CourseClassification = "past"
	CourseClassificationInProgress CourseClassification = "inprogress"
	CourseClassificationFuture     CourseClassification = "future"
)

type Course struct {
	ID              int
	FullName        string
	ShortName       string
	Summary         string
	SummaryFormat   int
	StartDate       time.Time
	EndDate         time.Time
	Visible         bool
	FullNameDisplay string
	ViewURL         string
	CourseImage     string
	Progress        int
	HasProgress     bool
	IsSavourite     bool
	Hidden          bool
	ShowShortName   bool
	CourseCategory  string
}

func (c *courseAPI) GetEnrolledCoursesByTimelineClassification(ctx context.Context, classification CourseClassification) ([]*Course, error) {
	res, err := c.getEnrolledCoursesByTimelineClassification(ctx, classification)
	if err != nil {
		return nil, err
	}
	return mapFromCourseListResponse(res.Courses), nil
}

func mapFromCourseListResponse(courseResList []*courseResponse) []*Course {
	courses := make([]*Course, 0, len(courseResList))
	for _, courseRes := range courseResList {
		courses = append(courses, mapFromCourseResponse(courseRes))
	}
	return courses
}

func mapFromCourseResponse(courseRes *courseResponse) *Course {
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
