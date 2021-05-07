package moodle

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_siteAPI_GetSiteInfo(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		response string
		want     *SiteInfo
		wantErr  bool
	}{
		{
			name: "Successful response",
			args: args{ctx: context.Background()},
			response: `{
  "sitename": "Test Site",
  "username": "s111111",
  "firstname": "Test",
  "lastname": "User",
  "fullname": "Test user",
  "lang": "en",
  "userid": 111111,
  "siteurl": "https:\/\/test.edu",
  "userpictureurl": "https:\/\/test.edu\/pluginfile.php\/111111\/user\/icon\/lambda\/f1",
  "functions": [
    {
      "name": "core_files_get_files",
      "version": "2019052009"
    }
  ],
  "downloadfiles": 1,
  "uploadfiles": 1,
  "release": "3.7.9 (Build: 20201109)",
  "version": "2019052009",
  "mobilecssurl": "",
  "advancedfeatures": [
    {
      "name": "usecomments",
      "value": 1
    },
    {
      "name": "enableblogs",
      "value": 0
    }
  ],
  "usercanmanageownfiles": true,
  "userquota": 104857600,
  "usermaxuploadfilesize": 104857600,
  "userhomepage": 0,
  "siteid": 1,
  "sitecalendartype": "gregorian",
  "usercalendartype": "gregorian",
  "theme": "lambda"
}`,
			want: &SiteInfo{
				SiteName:              "Test Site",
				Username:              "s111111",
				Firstname:             "Test",
				Lastname:              "User",
				Fullname:              "Test user",
				Lang:                  "en",
				UserID:                111111,
				SiteURL:               "https://test.edu",
				UserPictureURL:        "https://test.edu/pluginfile.php/111111/user/icon/lambda/f1",
				Functions:             []*SiteFunctionVersion{{Name: "core_files_get_files", Version: "2019052009"}},
				DownloadFiles:         true,
				UploadFiles:           true,
				Release:               "3.7.9 (Build: 20201109)",
				Version:               "2019052009",
				AdvancedFeatures:      []*AdvancedFeatureEnabled{{Name: "usecomments", Enabled: true}, {Name: "enableblogs", Enabled: false}},
				UserCanManageOwnFiles: true,
				UserQuota:             104857600,
				UserMaxUploadFileSize: 104857600,
				UserHomePage:          0,
				SiteID:                1,
				SiteCalendarType:      "gregorian",
				UserCalendarType:      "gregorian",
				Theme:                 "lambda",
			},
		},
		{
			name:     "Error response",
			args:     args{ctx: context.Background()},
			response: `{"errorcode": "invalidtoken"}`,
			wantErr:  true,
		},
		{
			name:     "Invalid json response",
			args:     args{ctx: context.Background()},
			response: "{",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := mockSiteAPI(t, tt.response)
			got, err := s.GetSiteInfo(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSiteInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("GetSiteInfo() (-got, +want)\n%s", diff)
			}
		})
	}
}

func mockSiteAPI(t *testing.T, response string) *siteAPI {
	t.Helper()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	})
	s := httptest.NewServer(h)
	apiURL, _ := url.Parse(s.URL)
	return &siteAPI{
		&apiClient{
			httpClient: http.DefaultClient,
			apiURL:     apiURL,
		},
	}
}
