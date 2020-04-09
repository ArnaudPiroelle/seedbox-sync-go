package notifier

import (
	"seedbox-sync/model"
	"seedbox-sync/provider"
)

type ComposeNotifier struct {
	notifiers []Notifier
}

func NewCompose(notifiers ...Notifier) *ComposeNotifier {
	return &ComposeNotifier{notifiers: notifiers}
}

func (n *ComposeNotifier) StartSynchro() {
	for _, notifier := range n.notifiers {
		notifier.StartSynchro()
	}
}

func (n *ComposeNotifier) EndSynchro() {
	for _, notifier := range n.notifiers {
		notifier.EndSynchro()
	}
}

func (n *ComposeNotifier) StartFolder(folder model.Folder) {
	for _, notifier := range n.notifiers {
		notifier.StartFolder(folder)
	}
}

func (n *ComposeNotifier) EndFolder(folder model.Folder) {
	for _, notifier := range n.notifiers {
		notifier.EndFolder(folder)
	}
}

func (n *ComposeNotifier) StartTorrent(torrent provider.Torrent) {
	for _, notifier := range n.notifiers {
		notifier.StartTorrent(torrent)
	}
}

func (n *ComposeNotifier) EndTorrent(torrent provider.Torrent) {
	for _, notifier := range n.notifiers {
		notifier.EndTorrent(torrent)
	}
}

func (n *ComposeNotifier) StartFile(file provider.TorrentFile) {
	for _, notifier := range n.notifiers {
		notifier.StartFile(file)
	}
}

func (n *ComposeNotifier) ProgressFile(file provider.TorrentFile, bytesRead int64, totalBytesRead int64) {
	for _, notifier := range n.notifiers {
		notifier.ProgressFile(file, bytesRead, totalBytesRead)
	}
}

func (n *ComposeNotifier) EndFile(file provider.TorrentFile, success bool) {
	for _, notifier := range n.notifiers {
		notifier.EndFile(file, success)
	}
}
