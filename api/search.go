package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/videofree/torrents/bittorrent"
	"github.com/videofree/torrents/providers"
	"github.com/videofree/torrents/config"
	"github.com/videofree/torrents/xbmc"
)

var searchLog = logging.MustGetLogger("search")
var searchHistory []string

func Search(btService *bittorrent.BTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		query := ctx.Query("q")

		if len(searchHistory) > 0 && xbmc.DialogConfirm("Torrents", "LOCALIZE[30262]") {
			choice := xbmc.ListDialog("LOCALIZE[30261]", searchHistory...)
			query = searchHistory[choice]
		} else if query == "" {
			query = xbmc.Keyboard("", "LOCALIZE[30209]")
			if query == "" {
				return
			}
			searchHistory = append(searchHistory, query)
		}
		if query == "" {
			return
		}

		existingTorrent := ExistingTorrent(btService, query)
		if existingTorrent != "" && xbmc.DialogConfirm("Torrents", "LOCALIZE[30270]") {
			xbmc.PlayURL(UrlQuery(UrlForXBMC("/play"), "uri", existingTorrent))
			return
		}

		searchLog.Infof("Searching providers for: %s", query)

		searchers := providers.GetSearchers()
		torrents := providers.Search(searchers, query)

		if len(torrents) == 0 {
			xbmc.Notify("Torrents", "LOCALIZE[30205]", config.AddonIcon())
			return
		}

		choices := make([]string, 0, len(torrents))
		for _, torrent := range torrents {
			resolution := ""
			if torrent.Resolution > 0 {
				resolution = fmt.Sprintf("[B][COLOR %s]%s[/COLOR][/B] ", bittorrent.Colors[torrent.Resolution], bittorrent.Resolutions[torrent.Resolution])
			}

			info := make([]string, 0)
			if torrent.Size != "" {
				info = append(info, fmt.Sprintf("[B][%s][/B]", torrent.Size))
			}
			if torrent.RipType > 0 {
				info = append(info, bittorrent.Rips[torrent.RipType])
			}
			if torrent.VideoCodec > 0 {
				info = append(info, bittorrent.Codecs[torrent.VideoCodec])
			}
			if torrent.AudioCodec > 0 {
				info = append(info, bittorrent.Codecs[torrent.AudioCodec])
			}
			if torrent.Provider != "" {
				info = append(info, fmt.Sprintf(" - [B]%s[/B]", torrent.Provider))
			}

			multi := ""
			if torrent.Multi {
				multi = "\nmulti"
			}

			label := fmt.Sprintf("%s(%d / %d) %s\n%s\n%s%s",
				resolution,
				torrent.Seeds,
				torrent.Peers,
				strings.Join(info, " "),
				torrent.Name,
				torrent.Icon,
				multi,
			)
			choices = append(choices, label)
		}

		choice := xbmc.ListDialogLarge("LOCALIZE[30228]", query, choices...)
		if choice >= 0 {
			xbmc.PlayURL(UrlQuery(UrlForXBMC("/play"), "uri", torrents[choice].URI))
		}
	}
}
