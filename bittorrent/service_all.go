// +build !arm

package bittorrent

import "github.com/videofree/libtorrent-go"

// Nothing to do on regular devices
func setPlatformSpecificSettings(settings libtorrent.SettingsPack) {
}
