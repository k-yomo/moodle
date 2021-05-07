package main

import (
	"context"
	"fmt"
	"github.com/k-yomo/moodle"
	"net/url"
)

func main() {
	ctx := context.Background()
	serviceURL, err := url.Parse("https://my.uopeople.edu")
	if err != nil {
		panic(err)
	}
	moodleClient, err := moodle.NewClientWithLogin(
		ctx,
		serviceURL,
		"SXXXXXX",
		"password",
	)
	if err != nil {
		panic(err)
	}

	siteInfo, err := moodleClient.SiteAPI.GetSiteInfo(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", siteInfo)

	courses, err := moodleClient.CourseAPI.GetEnrolledCoursesByTimelineClassification(
		ctx,
		moodle.CourseClassificationInProgress,
	)
	if err != nil {
		panic(err)
	}

	for _, c := range courses {
		fmt.Printf("%#v\n", c)
	}
}
