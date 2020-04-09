package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const csrfHeader = "X-Transmission-Session-Id"

type Transmission struct {
	url             string
	user            string
	password        string
	httpC           *http.Client
	sessionID       string
	sessionIDAccess sync.RWMutex
	rnd             *rand.Rand
}

func NewTransmission(url, username, password string) *Transmission {
	return &Transmission{
		url:      url,
		user:     username,
		password: password,
		httpC:    cleanhttp.DefaultPooledClient(),
		rnd:      rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (t *Transmission) GetTorrents() (torrents []Torrent, err error) {
	var result torrentGetResults
	if err = t.rpcCall("torrent-get", torrentGetParams{
		Fields: []string{
			"id",
			"name",
			"percentDone",
			"downloadDir",
			"files",
		},
		IDs: nil,
	}, &result); err != nil {
		err = fmt.Errorf("'torrent-get' rpc method failed: %v", err)
		return
	}

	torrents = make([]Torrent, len(result.Torrents))
	for i, torrent := range result.Torrents {

		var files = make([]TorrentFile, len(torrent.Files))
		for j, file := range torrent.Files {
			files[j] = TorrentFile{
				Name:           file.Name,
				Length:         file.Length,
				BytesCompleted: file.BytesCompleted,
			}
		}

		torrents[i] = Torrent{
			Id:          *torrent.ID,
			Name:        *torrent.Name,
			PercentDone: *torrent.PercentDone,
			Files:       files,
			DownloadDir: *torrent.DownloadDir,
		}
	}
	return
}

func (t *Transmission) SetLocation(torrent Torrent, remoteSharePath string) (err error) {
	if err = t.rpcCall("torrent-set-location", torrentSetLocationPayload{
		IDs:      []int64{torrent.Id},
		Location: remoteSharePath,
		Move:     true,
	}, nil); err != nil {
		err = fmt.Errorf("'torrent-set-location' rpc method failed: %v", err)
	}
	return
}

type torrentGetParams struct {
	Fields []string `json:"fields"`
	IDs    []int64  `json:"ids,omitempty"`
}

type torrentGetResults struct {
	Torrents []*torrent `json:"torrents"`
}

type torrent struct {
	DownloadDir *string        `json:"downloadDir"`
	Files       []*torrentFile `json:"files"`
	ID          *int64         `json:"id"`
	PercentDone *float64       `json:"percentDone"`
	Name        *string        `json:"name"`
}

type torrentFile struct {
	BytesCompleted int64  `json:"bytesCompleted"`
	Length         int64  `json:"length"`
	Name           string `json:"name"`
}

type torrentSetLocationPayload struct {
	IDs      []int64 `json:"ids"`
	Location string  `json:"location"`
	Move     bool    `json:"move"`
}

type requestPayload struct {
	Method    string      `json:"method"`
	Arguments interface{} `json:"arguments,omitempty"`
	Tag       int         `json:"tag,omitempty"`
}

type answerPayload struct {
	Arguments interface{} `json:"arguments"`
	Result    string      `json:"result"`
	Tag       *int        `json:"tag"`
}

func (t *Transmission) rpcCall(method string, arguments interface{}, result interface{}) (err error) {
	return t.request(method, arguments, result, true)
}

func (t *Transmission) request(method string, arguments interface{}, result interface{}, retry bool) (err error) {
	// Let's avoid crashing
	if t.httpC == nil {
		err = errors.New("this controller is not initialized, please use the New() function")
		return
	}
	// Prepare the pipeline between payload generation and request
	pOut, pIn := io.Pipe()
	// Prepare the request
	var req *http.Request
	if req, err = http.NewRequest("POST", t.url, pOut); err != nil {
		err = fmt.Errorf("can't prepare request for '%s' method: %v", method, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(csrfHeader, t.getSessionID())
	req.SetBasicAuth(t.user, t.password)
	// Prepare the marshalling goroutine
	var tag int
	var encErr error
	var mg sync.WaitGroup
	mg.Add(1)
	go func() {
		tag = t.rnd.Int()
		encErr = json.NewEncoder(pIn).Encode(&requestPayload{
			Method:    method,
			Arguments: arguments,
			Tag:       tag,
		})
		pIn.Close()
		mg.Done()
	}()
	// Execute request
	var resp *http.Response
	if resp, err = t.httpC.Do(req); err != nil {
		mg.Wait()
		if encErr != nil {
			err = fmt.Errorf("request error: %v | json payload marshall error: %v", err, encErr)
		} else {
			err = fmt.Errorf("request error: %v", err)
		}
		return
	}
	defer resp.Body.Close()
	// Let's test the enc result, just in case
	mg.Wait()
	if encErr != nil {
		err = fmt.Errorf("request payload JSON marshalling failed: %v", encErr)
		return
	}
	// Is the CRSF token invalid ?
	if resp.StatusCode == http.StatusConflict {
		// Recover new token and save it
		t.updateSessionID(resp.Header.Get(csrfHeader))
		// Retry request if first try
		if retry {
			return t.request(method, arguments, result, false)
		}
		err = errors.New("CSRF token invalid 2 times in a row: stopping to avoid infinite loop")
		return
	}
	// Is request successful ?
	if resp.StatusCode != 200 {
		err = fmt.Errorf("HTTP error %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		return
	}

	answer := answerPayload{
		Arguments: result,
	}
	if err = json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		err = fmt.Errorf("can't unmarshall request answer body: %v", err)
		return
	}
	// fmt.Println("DEBUG >", answer.Result)
	// Final checks
	if answer.Tag == nil {
		err = errors.New("http answer does not have a tag within it's payload")
		return
	}
	if *answer.Tag != tag {
		err = errors.New("http request tag and answer payload tag do not match")
		return
	}
	if answer.Result != "success" {
		err = fmt.Errorf("http request ok but payload does not indicate success: %s", answer.Result)
		return
	}
	// All good
	return
}

func (t *Transmission) getSessionID() string {
	defer t.sessionIDAccess.RUnlock()
	t.sessionIDAccess.RLock()
	return t.sessionID
}

func (t *Transmission) updateSessionID(newID string) {
	defer t.sessionIDAccess.Unlock()
	t.sessionIDAccess.Lock()
	t.sessionID = newID
}
