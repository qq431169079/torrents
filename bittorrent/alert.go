package bittorrent

import "github.com/videofree/libtorrent-go"

type ltAlert struct {
	libtorrent.Alert
}

type Alert struct {
	Type     int
	Category int
	What     string
	Message  string
	Pointer  uintptr
	Name     string
	Entry    libtorrent.Entry
	InfoHash string
}
