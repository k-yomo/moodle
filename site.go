package moodle

type SiteInfo struct {
	SiteName              string
	Username              string
	Firstname             string
	Lastname              string
	Fullname              string
	Lang                  string
	UserID                int
	SiteURL               string
	UserPictureURL        string
	Functions             []*SiteFunctionVersion
	DownloadFiles         bool
	UploadFiles           bool
	Release               string
	Version               string
	MobileCSSURL          string
	AdvancedFeatures      []*AdvancedFeatureEnabled
	UserCanManageOwnFiles bool
	UserQuota             int
	UserMaxUploadFileSize int
	UserHomePage          int
	SiteID                int
	SiteCalendarType      string
	UserCalendarType      string
	Theme                 string
}

type SiteFunctionVersion struct {
	Name    string
	Version string
}

type AdvancedFeatureEnabled struct {
	Name    string
	Enabled bool
}
