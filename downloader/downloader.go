package downloader

import (
	"io"
	"seedbox-sync/provider"
)

type Downloader interface {
	Connect()
	Disconnect()
	GetFile(file string, resumeAt uint64) (io.Reader, error)
	GetRemoteSize(file provider.TorrentFile, remoteCompletePath string) (size int64, err error)
	GetRoot() string
}
