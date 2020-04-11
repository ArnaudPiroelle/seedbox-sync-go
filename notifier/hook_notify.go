package notifier

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"io"
	"net/http"
	"seedbox-sync/model"
	"seedbox-sync/provider"
)

type HookNotifier struct {
	hooks []model.Hook
	httpC *http.Client
}

func NewHookNotifier(hooks []model.Hook) *HookNotifier {
	return &HookNotifier{
		hooks: hooks,
		httpC: cleanhttp.DefaultPooledClient(),
	}
}

func (n *HookNotifier) StartSynchro() {
	n.call("sync/pre")
}

func (n *HookNotifier) EndSynchro() {
	n.call("sync/post")
}

func (n *HookNotifier) StartFolder(folder model.Folder) {
	n.call("folder/pre")
}

func (n *HookNotifier) EndFolder(folder model.Folder) {
	n.call("folder/post")
}

func (n *HookNotifier) StartTorrent(torrent provider.Torrent) {
	n.call("download/pre")

}

func (n *HookNotifier) EndTorrent(torrent provider.Torrent) {
	n.call("download/post")

}

func (n *HookNotifier) StartFile(file provider.TorrentFile) {

}

func (n *HookNotifier) ProgressFile(file provider.TorrentFile, bytesRead int64, totalBytesRead int64) {

}

func (n *HookNotifier) EndFile(file provider.TorrentFile, success bool) {

}

func (n *HookNotifier) call(event string) {
	for _, hook := range n.hooks {
		if hook.Event == event {
			_ = n.request(hook)
		}
	}
}

func (n *HookNotifier) request(hook model.Hook) (err error) {
	if n.httpC == nil {
		err = errors.New("this controller is not initialized, please use the New() function")
		return
	}

	// Prepare the pipeline between payload generation and request
	pOut, _ := io.Pipe()
	// Prepare the request
	var req *http.Request
	if req, err = http.NewRequest(hook.Method, hook.Url, pOut); err != nil {
		err = fmt.Errorf("can't prepare request for '%s' method: %v", hook.Method, err)
		return
	}

	_, _ = n.httpC.Do(req)

	return
}
