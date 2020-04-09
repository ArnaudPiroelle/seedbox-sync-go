package downloader

import (
	"io"
	"log"
	"seedbox-sync/provider"
	"strconv"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTP struct {
	host     string
	port     int
	username string
	password string
	root     string
	client   *ftp.ServerConn
}

func NewFtp(host string, port int, username, password, root string) *FTP {
	return &FTP{
		host:     host,
		port:     port,
		username: username,
		password: password,
		root:     root,
	}
}

func (f *FTP) Connect() {
	c, err := ftp.Dial(f.host+":"+strconv.Itoa(f.port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	f.client = c
	err = c.Login(f.username, f.password)
	if err != nil {
		log.Fatal(err)
	}
}

func (f *FTP) Disconnect() {
	err := f.client.Quit()
	if err != nil {
		log.Fatal(err)
	}
}

func (f *FTP) GetFile(file string, resumeAt uint64) (io.Reader, error) {
	resp, err := f.client.RetrFrom(file, resumeAt)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (f *FTP) GetRemoteSize(file provider.TorrentFile, remoteCompletePath string) (size int64, err error) {
	root := f.GetRoot()
	remotePath := strings.Replace(remoteCompletePath, root, "", 1)
	remoteFile := remotePath + "/" + file.Name
	return f.client.FileSize(remoteFile)
}

func (f *FTP) GetRoot() string {
	return f.root
}
