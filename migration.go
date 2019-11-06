package main

import (
	"os"
	"path/filepath"

	"github.com/videofree/torrents/xbmc"
	"github.com/videofree/torrents/config"
	"github.com/videofree/torrents/repository"
)

func Migrate() bool {
	firstRun := filepath.Join(config.Get().Info.Path, ".firstrun")
	if _, err := os.Stat(firstRun); err == nil {
		return false
	}
	file, _ := os.Create(firstRun)
	defer file.Close()

	log.Info("Preparing for first run...")

	log.Info("Creating Torrents repository add-on...")
	if err := repository.MakeTorrentsRepositoryAddon(); err != nil {
		log.Errorf("Unable to create repository add-on: %s", err)
	} else {
		xbmc.UpdateLocalAddons()
		for _, addon := range xbmc.GetAddons("xbmc.addon.repository", "unknown", "all", []string{"name", "version", "enabled"}).Addons {
			if addon.ID == "repository.torrents" && addon.Enabled == true {
				log.Info("Found enabled Torrents repository add-on")
				return false
			}
		}
		log.Info("Torrents repository not installed, installing...")
		xbmc.InstallAddon("repository.torrents")
		xbmc.SetAddonEnabled("repository.torrents", true)
		xbmc.UpdateLocalAddons()
		xbmc.UpdateAddonRepos()
	}

	return true
}
