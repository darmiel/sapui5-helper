package s5

import (
	"encoding/json"
	"os"
	"path"
)

type ManifestApp struct {
}

type Manifest struct {
	App struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		I18n       string `json:"i18n"`
		AppVersion struct {
			Version string `json:"version"`
		} `json:"applicationVersion"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Resources   string `json:"resources"`
		Ach         string `json:"ach"`
	} `json:"sap.app"`
	UI struct {
		Technology  string `json:"technology"`
		DeviceTypes struct {
			Desktop bool `json:"desktop"`
			Tablet  bool `json:"tablet"`
			Phone   bool `json:"phone"`
		} `json:"deviceTypes"`
	} `json:"sap.ui"`
}

func ReadManifestFromFile(file string) (resp *Manifest, err error) {
	var data []byte
	if data, err = os.ReadFile(file); err != nil {
		return
	}
	err = json.Unmarshal(data, &resp)
	return
}

func ReadManifest() (*Manifest, error) {
	if ifo, err := os.Stat("webapp"); err == nil && ifo != nil && ifo.IsDir() {
		return ReadManifestFromFile(path.Join("webapp", "manifest.json"))
	}
	return ReadManifestFromFile(path.Join("webapp", "manifest.json"))
}
