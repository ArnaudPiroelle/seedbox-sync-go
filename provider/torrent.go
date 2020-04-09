package provider

type Torrent struct {
	Id          int64
	Name        string
	PercentDone float64
	Files       []TorrentFile
	DownloadDir string
}

type TorrentFile struct {
	Name           string
	Length         int64
	BytesCompleted int64
}

func (t TorrentFile) IsCompleted() bool {
	return t.Length == t.BytesCompleted
}
