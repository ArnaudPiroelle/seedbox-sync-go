package notifier

import (
	"fmt"
	"seedbox-sync/model"
	"seedbox-sync/provider"
)

type LoggerNotifier struct {
}

func NewLogger() *LoggerNotifier {
	return &LoggerNotifier{}
}

func (n *LoggerNotifier) StartSynchro() {
	fmt.Println("Start synchronisation...")
}

func (n *LoggerNotifier) EndSynchro() {
	fmt.Println("Synchronisation ended")
}

func (n *LoggerNotifier) StartFolder(folder model.Folder) {
	fmt.Println("Start download folder ", folder.RemoteCompletePath)
}

func (n *LoggerNotifier) EndFolder(folder model.Folder) {

}

func (n *LoggerNotifier) StartTorrent(torrent provider.Torrent) {
	fmt.Println(torrent.Name)
}

func (n *LoggerNotifier) EndTorrent(torrent provider.Torrent) {

}

func (n *LoggerNotifier) StartFile(file provider.TorrentFile) {
	fmt.Println(file.Name)
}

func (n *LoggerNotifier) ProgressFile(file provider.TorrentFile, bytesRead int64, totalBytesRead int64) {

}

func (n *LoggerNotifier) EndFile(file provider.TorrentFile, success bool) {

}
