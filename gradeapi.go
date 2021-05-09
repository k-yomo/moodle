package moodle

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"time"
)

type GradeAPI interface {
	GetGradeItems(ctx context.Context, userID int, courseID int) ([]*UserGrade, error)
	GetGradesTable(ctx context.Context, userID int, courseID int) ([]*GradeTable, error)
}

type gradeAPI struct {
	*apiClient
}

func newGradeAPI(apiClient *apiClient) *gradeAPI {
	return &gradeAPI{apiClient}
}

type userGradeResponse struct {
	CourseID     int                  `json:"courseid"`
	UserID       int                  `json:"userid"`
	UserFullname string               `json:"userfullname"`
	MaxDepth     int                  `json:"maxdepth"`
	GradeItems   []*gradeItemResponse `json:"gradeitems"`
}

type gradeItemResponse struct {
	ID                     int     `json:"id"`
	ItemName               string  `json:"itemname"`
	ItemType               string  `json:"itemtype"`
	ItemModule             *string `json:"itemmodule"`
	ItemInstance           int     `json:"iteminstance"`
	ItemNumber             *int    `json:"itemnumber"`
	CategoryID             *int    `json:"categoryid"`
	OutcomeID              *int    `json:"outcomeid"`
	ScaleID                *int    `json:"scaleid"`
	Locked                 *bool   `json:"locked"`
	CmID                   int     `json:"cmid"`
	GradeRaw               float64 `json:"graderaw"`
	GradeDateSubmittedUnix *int64  `json:"gradedatesubmitted"`
	GradeDateGradedUnix    int64   `json:"gradedategraded"`
	GradeHiddenByDate      bool    `json:"gradehiddenbydate"`
	GradeNeedsUpdate       bool    `json:"gradeneedsupdate"`
	GradeIsHidden          bool    `json:"gradeishidden"`
	GradeIsLocked          *bool   `json:"gradeislocked"`
	GradeIsOverridden      *bool   `json:"gradeisoverridden"`
	GradeFormatted         string  `json:"gradeformatted"`
	GradeMin               int     `json:"grademin"`
	GradeMax               int     `json:"grademax"`
	RangeFormatted         string  `json:"rangeformatted"`
	Feedback               string  `json:"feedback"`
	FeedbackFormat         int     `json:"feedbackformat"`
}

type getGradeItemsResponse struct {
	UserGrades []*userGradeResponse `json:"usergrades"`
	Warnings   Warnings             `json:"warnings"`
}

func (g *gradeAPI) GetGradeItems(ctx context.Context, userID int, courseID int) ([]*UserGrade, error) {
	res := getGradeItemsResponse{}
	err := g.callMoodleFunction(ctx, &res, map[string]string{
		"wsfunction": "gradereport_user_get_grade_items",
		"userid":     strconv.Itoa(userID),
		"courseid":   strconv.Itoa(courseID),
	})
	if err != nil {
		return nil, err
	}
	if len(res.Warnings) > 0 {
		return nil, res.Warnings
	}

	return mapToUserGradeList(res.UserGrades), nil
}

type getGradesTableResponse struct {
	Tables []*struct {
		CourseID     int                      `json:"courseid"`
		UserID       int                      `json:"userid"`
		UserFullName string                   `json:"userfullname"`
		MaxDepth     int                      `json:"maxdepth"`
		TableDataRaw []interface{}            `json:"tabledata"` // this can be `[]`(empty array) or `tableDataItemResponse`
		TableData    []*tableDataItemResponse `json:"-"`
	} `json:"tables"`
	Warnings Warnings `json:"warnings"`
}

type tableDataItemResponse struct {
	ItemName *struct {
		ID       string `json:"id"`
		Class    string `json:"class"`
		ColSpan  int    `json:"colspan"`
		Content  string `json:"content"`
		CellType string `json:"celltype"`
	} `json:"itemname"`
	Leader *struct {
		Class   string `json:"class"`
		RowSpan int    `json:"rowspan"`
	} `json:"leader,omitempty"`
	Grade *struct {
		Class   string `json:"class"`
		Content string `json:"content"`
		Headers string `json:"headers"`
	} `json:"grade,omitempty"`
	Range *struct {
		Class   string `json:"class"`
		Content string `json:"content"`
		Headers string `json:"headers"`
	} `json:"range"`
	Feedback *struct {
		Class   string `json:"class"`
		Content string `json:"content"`
		Headers string `json:"headers"`
	} `json:"feedback"`
	ContributionToCourseTotal *struct {
		Class   string `json:"class"`
		Content string `json:"content"`
		Headers string `json:"headers"`
	} `json:"contributiontocoursetotal"`
}

func (g *gradeAPI) GetGradesTable(ctx context.Context, userID int, courseID int) ([]*GradeTable, error) {
	res := getGradesTableResponse{}
	err := g.callMoodleFunction(
		ctx,
		&res,
		map[string]string{
			"wsfunction": "gradereport_user_get_grades_table",
			"userid":     strconv.Itoa(userID),
			"courseid":   strconv.Itoa(courseID),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(res.Warnings) > 0 {
		return nil, res.Warnings
	}

	for i, t := range res.Tables {
		for _, td := range t.TableDataRaw {
			switch td.(type) {
			case []interface{}:
				continue
			case map[string]interface{}:
				tdJSON, err := json.Marshal(td)
				if err != nil {
					return nil, err
				}
				tdItem := tableDataItemResponse{}
				if err := json.Unmarshal(tdJSON, &tdItem); err != nil {
					return nil, err
				}
				res.Tables[i].TableData = append(res.Tables[i].TableData, &tdItem)
			}
		}
	}

	return mapToGradeTableList(&res)
}

func mapToUserGradeList(userGradeResList []*userGradeResponse) []*UserGrade {
	userGrades := make([]*UserGrade, 0, len(userGradeResList))
	for _, gradeItemRes := range userGradeResList {
		userGrades = append(userGrades, mapToUserGrade(gradeItemRes))
	}
	return userGrades
}

func mapToUserGrade(userGradeRes *userGradeResponse) *UserGrade {
	return &UserGrade{
		CourseID:     userGradeRes.CourseID,
		UserID:       userGradeRes.UserID,
		UserFullname: userGradeRes.UserFullname,
		MaxDepth:     userGradeRes.MaxDepth,
		GradeItems:   mapToGradeItemList(userGradeRes.GradeItems),
	}
}

func mapToGradeItemList(gradeItemResList []*gradeItemResponse) []*GradeItem {
	gradeItems := make([]*GradeItem, 0, len(gradeItemResList))
	for _, gradeItemRes := range gradeItemResList {
		gradeItems = append(gradeItems, mapToGradeItem(gradeItemRes))
	}
	return gradeItems
}

func mapToGradeItem(gradeItemRes *gradeItemResponse) *GradeItem {
	var gradeDateSubmitted, gradeDateGraded *time.Time
	if gradeItemRes.GradeDateSubmittedUnix != nil {
		t := time.Unix(*gradeItemRes.GradeDateSubmittedUnix, 0)
		gradeDateSubmitted = &t
	}
	if gradeItemRes.GradeDateGradedUnix > 0 {
		t := time.Unix(gradeItemRes.GradeDateGradedUnix, 0)
		gradeDateGraded = &t
	}
	return &GradeItem{
		ID:                 gradeItemRes.ID,
		ItemName:           gradeItemRes.ItemName,
		ItemType:           gradeItemRes.ItemType,
		ItemModule:         gradeItemRes.ItemModule,
		ItemInstance:       gradeItemRes.ItemInstance,
		ItemNumber:         gradeItemRes.ItemNumber,
		CategoryID:         gradeItemRes.CategoryID,
		OutcomeID:          gradeItemRes.OutcomeID,
		ScaleID:            gradeItemRes.ScaleID,
		Locked:             gradeItemRes.Locked,
		CmID:               gradeItemRes.CmID,
		GradeRaw:           gradeItemRes.GradeRaw,
		GradeDateSubmitted: gradeDateSubmitted,
		GradeDateGraded:    gradeDateGraded,
		GradeHiddenByDate:  gradeItemRes.GradeHiddenByDate,
		GradeNeedsUpdate:   gradeItemRes.GradeNeedsUpdate,
		GradeIsHidden:      gradeItemRes.GradeIsHidden,
		GradeIsLocked:      gradeItemRes.GradeIsLocked,
		GradeIsOverridden:  gradeItemRes.GradeIsOverridden,
		GradeFormatted:     gradeItemRes.GradeFormatted,
		GradeMin:           gradeItemRes.GradeMin,
		GradeMax:           gradeItemRes.GradeMax,
		RangeFormatted:     gradeItemRes.RangeFormatted,
		Feedback:           gradeItemRes.Feedback,
		FeedbackFormat:     gradeItemRes.FeedbackFormat,
	}
}

func mapToGradeTableList(res *getGradesTableResponse) ([]*GradeTable, error) {
	gradeTables := make([]*GradeTable, 0, len(res.Tables))
	for _, table := range res.Tables {
		tableItemGroups, err := mapToGradeTableItemGroupList(table.TableData)
		if err != nil {
			return nil, err
		}
		gradeTables = append(gradeTables, &GradeTable{
			CourseID:     table.CourseID,
			UserID:       table.UserID,
			UserFullname: table.UserFullName,
			MaxDepth:     table.MaxDepth,
			ItemGroups:   tableItemGroups,
		})
	}
	return gradeTables, nil
}

func mapToGradeTableItemGroupList(tableDataItemResList []*tableDataItemResponse) ([]*GradeTableItemGroup, error) {
	gradeTableItemGroups := make([]*GradeTableItemGroup, 0)

	var groupName string
	gradeTableItems := make([]*GradeTableItem, 0)
	for _, gradeItemRes := range tableDataItemResList {
		isLabelItem := gradeItemRes.Grade == nil
		if isLabelItem {
			if len(gradeTableItems) > 0 {
				gradeTableItemGroups = append(gradeTableItemGroups, &GradeTableItemGroup{
					Name:  groupName,
					Items: gradeTableItems,
				})
				gradeTableItems = []*GradeTableItem{}
			}
			itemNameDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(gradeItemRes.ItemName.Content)))
			if err != nil {
				return nil, err
			}
			groupName = itemNameDoc.Text()
		} else {
			// Exclude non graded item
			if gradeItemRes.ContributionToCourseTotal == nil || gradeItemRes.ContributionToCourseTotal.Content == "-" {
				continue
			}
			gradeTableItem, err := mapToGradeTableItem(gradeItemRes)
			if err != nil {
				return nil, err
			}
			gradeTableItems = append(gradeTableItems, gradeTableItem)
		}
	}
	gradeTableItemGroups = append(gradeTableItemGroups, &GradeTableItemGroup{
		Name:  groupName,
		Items: gradeTableItems,
	})
	return gradeTableItemGroups, nil
}

var floatValueRegex = regexp.MustCompile("([0-9]*[.])?[0-9]+")

func mapToGradeTableItem(tableDataItemRes *tableDataItemResponse) (*GradeTableItem, error) {
	itemNameDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(tableDataItemRes.ItemName.Content)))
	if err != nil {
		return nil, err
	}

	var itemURL *string
	if href, itemURLExist := itemNameDoc.Find("a").Attr("href"); itemURLExist {
		itemURL = &href
	}

	isGraded := tableDataItemRes.Grade.Content != "-"
	var grade float64
	if isGraded {
		grade, err = strconv.ParseFloat(tableDataItemRes.Grade.Content, 64)
		if err != nil {
			return nil, err
		}
	}

	// the content is something like "0-100"
	gradeRangeMatch := floatValueRegex.FindAllString(tableDataItemRes.Range.Content, 2)
	gradeRangeMin, err := strconv.ParseFloat(gradeRangeMatch[0], 64)
	if err != nil {
		return nil, err
	}
	gradeRangeMax, err := strconv.ParseFloat(gradeRangeMatch[1], 64)
	if err != nil {
		return nil, err
	}

	feedbackDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(tableDataItemRes.Feedback.Content)))
	if err != nil {
		return nil, err
	}

	// the content is something like "10.00 %"
	contributionMatch := floatValueRegex.FindStringSubmatch(tableDataItemRes.ContributionToCourseTotal.Content)
	contributionToCourseTotal, err := strconv.ParseFloat(contributionMatch[0], 64)
	if err != nil {
		return nil, err
	}

	return &GradeTableItem{
		ItemName:                  itemNameDoc.Text(),
		ItemNameRawHTML:           tableDataItemRes.ItemName.Content,
		ItemURL:                   itemURL,
		IsGraded:                  isGraded,
		Grade:                     grade,
		GradeRangeMin:             gradeRangeMin,
		GradeRangeMax:             gradeRangeMax,
		Feedback:                  feedbackDoc.Text(),
		FeedBackRawHTML:           tableDataItemRes.Feedback.Content,
		ContributionToCourseTotal: contributionToCourseTotal,
	}, nil
}
