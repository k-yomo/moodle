# moodle

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
		&moodle.LoginParams{
			Username: "S000000",
			Password: "password",
		},
	)
	if err != nil {
		panic(err)
	}

	courses, err := moodleClient.CourseAPI.GetEnrolledCoursesByTimelineClassification(
		ctx,
		moodle.CourseClassificationInProgress,
	)
	if err != nil {
		panic(err)
	}

	for _, c := range courses {
		fmt.Printf("%#v", c)
	}
}
```