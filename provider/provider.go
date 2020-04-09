package provider

type Provider interface {
	GetTorrents() (torrents []Torrent, err error)
	SetLocation(torrent Torrent, remoteSharePath string) (err error)
}
