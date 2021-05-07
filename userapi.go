package moodle

type UserAPI interface {
}

type userAPI struct {
	*apiClient
}

func newUserAPI(apiClient *apiClient) *userAPI {
	return &userAPI{apiClient}
}
