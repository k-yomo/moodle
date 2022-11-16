package moodle

import "context"

type UserAPI interface {
	CreateUsers(ctx context.Context)
}

type userAPI struct {
	*apiClient
}

func newUserAPI(apiClient *apiClient) *userAPI {
	return &userAPI{apiClient}
}

type userResponse struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
}
type UserRequest struct {
	createpassword    int         `json:"createpassword"`
	username          string      `json:"username"`
	auth              string      `json:"auth"`
	password          string      `json:"passowrd"`
	firstname         string      `json:"firstname"`
	lastname          string      `json:"lastname"`
	email             string      `json:"email"`
	maildisplay       int         `json:"maildisplay"`
	city              string      `json:"city"`
	country           string      `json:"country"`
	timezone          string      `json:"timezone"`
	description       string      `json:"description"`
	firstnamephonetic string      `json:"firstnamephonetic"`
	lastnamephonetic  string      `json:"lastnamephonetic"`
	middlename        string      `json:"middlename"`
	alternatename     string      `json:"alternatename"`
	interests         string      `json:"interests"`
	idnumber          string      `json:"idnumber"`
	institution       string      `json:"institution"`
	department        string      `json:"department"`
	phone1            string      `json:"phone1"`
	phone2            string      `json:"phone2"`
	address           string      `json:"address"`
	lang              string      `json:"lang"`
	calendartype      string      `json:"calendartype"`
	theme             string      `json:"theme"`
	mailformat        int         `json:"mailformat"`
	customfields      interface{} `json:"customfields"`
}

func (u *userAPI) CreateUsers(ctx context.Context, param []UserRequest) ([]userResponse, error) {
	res := []userResponse{}
	err := u.callMoodleFunctionPost(ctx, &res, param, map[string]string{
		"wsfunction": "core_user_create_users",
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
