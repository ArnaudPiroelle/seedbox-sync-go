package notifier

import (
	"github.com/cheggaaa/pb/v3"
	"seedbox-sync/model"
	"seedbox-sync/provider"
)

type ConsoleNotifier struct {
	bars map[string]*pb.ProgressBar
}

func NewConsole() *ConsoleNotifier {
	return &ConsoleNotifier{
		bars: map[string]*pb.ProgressBar{},
	}
}

func (n *ConsoleNotifier) StartSynchro() {

}

func (n *ConsoleNotifier) EndSynchro() {

}

func (n *ConsoleNotifier) StartFolder(folder model.Folder) {

}

func (n *ConsoleNotifier) EndFolder(folder model.Folder) {

}

func (n *ConsoleNotifier) StartTorrent(torrent provider.Torrent) {

}

func (n *ConsoleNotifier) EndTorrent(torrent provider.Torrent) {

}

func (n *ConsoleNotifier) StartFile(file provider.TorrentFile) {
	bar := pb.Full.Start64(file.Length)
	bar.Set(pb.Bytes, true)
	bar.Start()
	n.bars[file.Name] = bar
}

func (n *ConsoleNotifier) ProgressFile(file provider.TorrentFile, bytesRead int64, totalBytesRead int64) {
	bar := n.bars[file.Name]
	if bar != nil {
		bar.SetCurrent(bytesRead)
	}
}

func (n *ConsoleNotifier) EndFile(file provider.TorrentFile, success bool) {
	bar := n.bars[file.Name]
	bar.SetCurrent(file.Length)
	if bar != nil {
		bar.Finish()
	}
}
