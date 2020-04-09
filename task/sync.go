package task

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"seedbox-sync/downloader"
	"seedbox-sync/model"
	"seedbox-sync/notifier"
	"seedbox-sync/provider"
	"strings"
)

type Sync struct {
	folders    []model.Folder
	provider   provider.Provider
	downloader downloader.Downloader
	notifier   notifier.Notifier
}

func NewSync(folders []model.Folder, provider provider.Provider, downloader downloader.Downloader, notifier notifier.Notifier) *Sync {
	return &Sync{
		folders:    folders,
		provider:   provider,
		downloader: downloader,
		notifier:   notifier,
	}
}

func (s *Sync) Execute() {
	s.notifier.StartSynchro()

	torrents, err := s.provider.GetTorrents()
	if err != nil {
		fmt.Println("unable to retrieve the torrents")
		return
	}

	for _, folder := range s.folders {
		s.notifier.StartFolder(folder)
		for _, torrent := range torrents {
			if torrent.DownloadDir == folder.RemoteCompletePath && torrent.PercentDone == 1 {
				s.downloadTorrent(folder, torrent)
			}
		}
		s.notifier.EndFolder(folder)
		//TODO clean temp empty folders
	}

	s.notifier.EndSynchro()
}

func (s *Sync) downloadTorrent(folder model.Folder, torrent provider.Torrent) {
	s.notifier.StartTorrent(torrent)

	hasError := false
	for _, file := range torrent.Files {
		if file.IsCompleted() {
			e := s.downloadFile(folder, file)
			hasError = hasError || e != nil
		}
	}

	if !hasError {
		for _, file := range torrent.Files {
			s.moveFile(folder, file)
		}
		_ = s.provider.SetLocation(torrent, folder.RemoteSharePath)
		//TODO revert move if setlocation failed
	}

	s.notifier.EndTorrent(torrent)
}

func (s *Sync) downloadFile(folder model.Folder, file provider.TorrentFile) error {
	s.notifier.StartFile(file)

	d := s.downloader
	root := d.GetRoot()
	remotePath := strings.Replace(folder.RemoteCompletePath, root, "", 1)
	remoteFile := remotePath + "/" + file.Name
	localFile := folder.LocalTempPath + "/" + file.Name

	fi, err := os.Stat(localFile)
	var localSize int64 = 0
	if err == nil {
		localSize = fi.Size()
	}

	if localSize == 0 {
		parent := filepath.Dir(localFile)
		_ = os.MkdirAll(parent, 0755)
	}

	remoteSize, err := d.GetRemoteSize(file, folder.RemoteCompletePath)
	if err != nil {
		return err
	}
	if localSize < remoteSize {
		reader, err := s.downloader.GetFile(remoteFile, uint64(localSize))
		if err != nil {
			return err
		}
		proxyReader := &ProxyReader{
			value:    localSize,
			file:     file,
			reader:   reader,
			notifier: s.notifier,
		}

		open, err := os.OpenFile(localFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			return err
		}

		_, err = io.Copy(open, proxyReader)

		defer proxyReader.Close()
		defer open.Close()

		if err != nil {
			return err
		}
	}
	s.notifier.EndFile(file, true)
	return nil
}

func (s *Sync) moveFile(folder model.Folder, file provider.TorrentFile) {
	//TODO check available space before move a file
	oldName := folder.LocalTempPath + "/" + file.Name
	newName := folder.LocalPostProcessingPath + "/" + file.Name

	parent := filepath.Dir(newName)
	_ = os.MkdirAll(parent, 0755)
	_ = moveFile(oldName, newName)
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}

type ProxyReader struct {
	value    int64
	file     provider.TorrentFile
	reader   io.Reader
	notifier notifier.Notifier
}

func (r *ProxyReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	r.value += int64(n)
	r.notifier.ProgressFile(r.file, r.value, r.file.Length)
	return
}

func (r *ProxyReader) Close() (err error) {
	if closer, ok := r.reader.(io.Closer); ok {
		return closer.Close()
	}
	return
}
