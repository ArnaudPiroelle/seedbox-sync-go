package model

type Configuration struct {
	Provider   ProviderConfiguration   `json:"provider"`
	Downloader DownloaderConfiguration `json:"downloader"`
	Folders    []Folder                `json:"folders"`
	Hooks      []Hook                  `json:"hooks"`
	Scheduler  SchedulerConfiguration  `json:"scheduler"`
}

type ProviderConfiguration struct {
	Type     string `json:"type"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type DownloaderConfiguration struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Root     string `json:"root"`
}

type Folder struct {
	RemoteCompletePath      string `json:"remoteCompletePath"`
	RemoteSharePath         string `json:"remoteSharePath"`
	LocalTempPath           string `json:"localTempPath"`
	LocalPostProcessingPath string `json:"localPostProcessingPath"`
}

type Hook struct {
	Event  string `json:"event"`
	Method string `json:"method"`
	Url    string `json:"url"`
}

type SchedulerConfiguration struct {
	Cron string `json:"cron"`
}
