package config

type StartUpAndroid struct {
	LoadingIndex int `json:"loadingIndex"`
	GameSwitch   int `json:"gameSwitch"`
	ImageSize    struct {
		SrThumb struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"srThumb"`
		MainThumb struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"mainThumb"`
		ArticleThumb struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"articleThumb"`
		Sr struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"sr"`
		BarThumb struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"barThumb"`
		GroupThumb struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"groupThumb"`
		ZjThumb struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"zjThumb"`
		Bar struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"bar"`
		Group struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"group"`
	} `json:"imageSize"`
	Version struct {
		Url     string `json:"url"`
		Force   bool   `json:"force"`
		Code    int    `json:"code"`
		Version string `json:"version"`
		Desc    string `json:"desc"`
	} `json:"version"`
	QyWx           string `json:"qyWx"`
	HomeTaskSwitch int    `json:"homeTaskSwitch"`
	SrToastImg     string `json:"srToastImg"`
	SrToastJumpUrl string `json:"srToastJumpUrl"`
	PdB            string `json:"pdB"`
	PdWx           string `json:"pdWx"`
}
