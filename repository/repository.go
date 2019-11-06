package repository

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/videofree/torrents/config"
	"github.com/videofree/torrents/util"
	"github.com/videofree/torrents/xbmc"
)

func copyFile(from string, to string) error {
	input, err := os.Open(from)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(to)
	if err != nil {
		return err
	}
	defer output.Close()
	io.Copy(output, input)
	return nil
}

func MakeTorrentsRepositoryAddon() error {
	addonId := "repository.torrents"
	addonName := "Torrents Repository"

	quasarHost := fmt.Sprintf("http://localhost:%d", config.ListenPort)
	addon := &xbmc.Addon{
		Id:           addonId,
		Name:         addonName,
		Version:      util.Version[2:len(util.Version) - 1],
		ProviderName: config.Get().Info.Author,
		Extensions: []*xbmc.AddonExtension{
			&xbmc.AddonExtension{
				Point: "xbmc.addon.repository",
				Name:  addonName,
				Info: &xbmc.AddonRepositoryInfo{
					Text:       quasarHost + "/repository/videofree/plugin.video.torrents/addons.xml",
					Compressed: false,
				},
				Checksum: quasarHost + "/repository/videofree/plugin.video.torrents/addons.xml.md5",
				Datadir: &xbmc.AddonRepositoryDataDir{
					Text: quasarHost + "/repository/videofree/",
					Zip:  true,
				},
			},
			&xbmc.AddonExtension{
				Point: "xbmc.addon.metadata",
				Summaries: []*xbmc.AddonText{
					&xbmc.AddonText{
						Text: "GitHub repository for Torrents updates",
						Lang: "en",
					},
				},
				Platform: "all",
			},
		},
	}

	addonPath := filepath.Clean(filepath.Join(config.Get().Info.Path, "..", addonId))
	if err := os.MkdirAll(addonPath, 0777); err != nil {
		return err
	}

	if err := copyFile(filepath.Join(config.Get().Info.Path, "icon.png"), filepath.Join(addonPath, "icon.png")); err != nil {
		return err
	}

	if err := copyFile(filepath.Join(config.Get().Info.Path, "fanart.jpg"), filepath.Join(addonPath, "fanart.jpg")); err != nil {
		return err
	}

	addonXmlFile, err := os.Create(filepath.Join(addonPath, "addon.xml"))
	if err != nil {
		return err
	}
	defer addonXmlFile.Close()
	return xml.NewEncoder(addonXmlFile).Encode(addon)
}
