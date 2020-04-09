package notifier

import (
	"seedbox-sync/model"
	"seedbox-sync/provider"
)

type Notifier interface {
	StartSynchro()
	EndSynchro()
	StartFolder(folder model.Folder)
	EndFolder(folder model.Folder)
	StartTorrent(torrent provider.Torrent)
	EndTorrent(torrent provider.Torrent)
	StartFile(file provider.TorrentFile)
	ProgressFile(file provider.TorrentFile, bytesRead int64, totalBytesRead int64)
	EndFile(file provider.TorrentFile, success bool)
}
