package util

import (
	"fmt"
	"github.com/videofree/libtorrent-go"
)

var (
	Version   string
)

func UserAgent() string {
	return fmt.Sprintf("Torrents/%s libtorrent/%s", Version[1:len(Version) - 1], libtorrent.Version())
}
