package api

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/videofree/torrents/cloudhole"
	"github.com/videofree/torrents/config"
	"github.com/videofree/torrents/xbmc"
)

var cmdLog = logging.MustGetLogger("cmd")

func ClearCache(ctx *gin.Context) {
	clearPageCache(ctx)
	xbmc.Notify("Torrents", "LOCALIZE[30200]", config.AddonIcon())
}

func ClearPageCache(ctx *gin.Context) {
	clearPageCache(ctx)
}

func ResetClearances(ctx *gin.Context) {
	cloudhole.ResetClearances()
	xbmc.Notify("Torrents", "LOCALIZE[30264]", config.AddonIcon())
}

func SetViewMode(ctx *gin.Context) {
	content_type := ctx.Params.ByName("content_type")
	viewName := xbmc.InfoLabel("Container.Viewmode")
	viewMode := xbmc.GetCurrentView()
	cmdLog.Noticef("ViewMode: %s (%s)", viewName, viewMode)
	if viewMode != "0" {
		xbmc.SetSetting("viewmode_" + content_type, viewMode)
	}
	ctx.String(200, "")
}
