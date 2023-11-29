package config

type UpdateVersion struct {
	Ios struct {
		Url     string `json:"url"`
		Force   bool   `json:"force"`
		Code    int    `json:"code"`
		Version string `json:"version"`
		Desc    string `json:"desc"`
	} `json:"ios"`
	Android struct {
		Url     string `json:"url"`
		Force   bool   `json:"force"`
		Code    int    `json:"code"`
		Version string `json:"version"`
		Desc    string `json:"desc"`
	} `json:"android"`
}
