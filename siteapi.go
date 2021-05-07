package moodle

import (
	"context"
	"github.com/k-yomo/moodle/pkg/urlutil"
	"net/http"
	"net/url"
)

type SiteAPI interface {
	GetSiteInfo(ctx context.Context) (*SiteInfo, error)
}

type siteAPI struct {
	httpClient *http.Client
	apiURL     *url.URL
}

func newSiteAPI(httpClient *http.Client, apiURL *url.URL) *siteAPI {
	return &siteAPI{
		httpClient: httpClient,
		apiURL:     apiURL,
	}
}

type siteInfoResponse struct {
	SiteName       string `json:"sitename"`
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Fullname       string `json:"fullname"`
	Lang           string `json:"lang"`
	UserID         int    `json:"userid"`
	SiteURL        string `json:"siteurl"`
	UserPictureURL string `json:"userpictureurl"`
	Functions      []*struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"functions"`
	DownloadFiles    int    `json:"downloadfiles"`
	UploadFiles      int    `json:"uploadfiles"`
	Release          string `json:"release"`
	Version          string `json:"version"`
	MobileCSSURL     string `json:"mobilecssurl"`
	AdvancedFeatures []*struct {
		Name    string `json:"name"`
		Enabled int    `json:"value"`
	} `json:"advancedfeatures"`
	UserCanManageOwnFiles bool   `json:"usercanmanageownfiles"`
	UserQuota             int    `json:"userquota"`
	UserMaxUploadFileSize int    `json:"usermaxuploadfilesize"`
	UserHomePage          int    `json:"userhomepage"`
	SiteID                int    `json:"siteid"`
	SiteCalendarType      string `json:"sitecalendartype"`
	UserCalendarType      string `json:"usercalendartype"`
	Theme                 string `json:"theme"`
}

func (u *siteAPI) GetSiteInfo(ctx context.Context) (*SiteInfo, error) {
	url := urlutil.CopyWithQueries(u.apiURL, map[string]string{
		"wsfunction": "core_webservice_get_site_info",
	})
	res := siteInfoResponse{}
	if err := getAndUnmarshal(ctx, u.httpClient, url, &res); err != nil {
		return nil, err
	}

	return mapToSiteInfo(&res), nil
}

func mapToSiteInfo(siteInfoRes *siteInfoResponse) *SiteInfo {
	functions := make([]*SiteFunctionVersion, 0, len(siteInfoRes.Functions))
	for _, f := range siteInfoRes.Functions {
		functions = append(functions, &SiteFunctionVersion{
			Name:    f.Name,
			Version: f.Version,
		})
	}
	advancedFeatures := make([]*AdvancedFeatureEnabled, 0, len(siteInfoRes.AdvancedFeatures))
	for _, a := range siteInfoRes.AdvancedFeatures {
		advancedFeatures = append(advancedFeatures, &AdvancedFeatureEnabled{
			Name:    a.Name,
			Enabled: mapBitToBool(a.Enabled),
		})
	}
	return &SiteInfo{
		SiteName:              siteInfoRes.SiteName,
		Username:              siteInfoRes.Username,
		Firstname:             siteInfoRes.Firstname,
		Lastname:              siteInfoRes.Lastname,
		Fullname:              siteInfoRes.Fullname,
		Lang:                  siteInfoRes.Lang,
		UserID:                siteInfoRes.UserID,
		SiteURL:               siteInfoRes.SiteURL,
		UserPictureURL:        siteInfoRes.UserPictureURL,
		Functions:             functions,
		DownloadFiles:         mapBitToBool(siteInfoRes.DownloadFiles),
		UploadFiles:           mapBitToBool(siteInfoRes.UploadFiles),
		Release:               siteInfoRes.Release,
		Version:               siteInfoRes.Version,
		MobileCSSURL:          siteInfoRes.MobileCSSURL,
		AdvancedFeatures:      advancedFeatures,
		UserCanManageOwnFiles: siteInfoRes.UserCanManageOwnFiles,
		UserQuota:             siteInfoRes.UserQuota,
		UserMaxUploadFileSize: siteInfoRes.UserMaxUploadFileSize,
		UserHomePage:          siteInfoRes.UserHomePage,
		SiteID:                siteInfoRes.SiteID,
		SiteCalendarType:      siteInfoRes.SiteCalendarType,
		UserCalendarType:      siteInfoRes.UserCalendarType,
		Theme:                 siteInfoRes.Theme,
	}
}
