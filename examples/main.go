package main

import (
	"context"
	"github.com/k-yomo/moodle"
	"github.com/k0kubun/pp"
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
			Username: "SXXXXXX",
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
		pp.Println(c)
		quizzes, err := moodleClient.QuizAPI.GetQuizzesByCourse(ctx, c.ID)
		if err != nil {
			panic(err)
		}
		for _, q := range quizzes {
			pp.Println(q)
		}
	}
}
