# moodle

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)
![Tests Workflow](https://github.com/k-yomo/moodle/workflows/Tests/badge.svg)
[![codecov](https://codecov.io/gh/k-yomo/moodle/branch/main/graph/badge.svg)](https://codecov.io/gh/k-yomo/moodle)
[![Go Report Card](https://goreportcard.com/badge/k-yomo/moodle)](https://goreportcard.com/report/k-yomo/moodle)

Go Moodle API Client

## Example

```go
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
```