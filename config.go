package noteofitcli

import (
	"errors"
	"path/filepath"

	"github.com/donatj/appsettings"
	"github.com/shibukawa/configdir"
)

const keyToken = "token"

type Config struct {
	settings *appsettings.AppSettings
}

func NewConfig() (*Config, error) {
	configDirs := configdir.New("donatstudios", "noteofit-cli")

	folders := configDirs.QueryFolders(configdir.Global)
	cache := configDirs.QueryFolders(configdir.Cache)
	if len(folders) < 1 || len(cache) < 1 {
		return nil, errors.New("unable to store config")
	}

	folders[0].MkdirAll()
	cache[0].MkdirAll()

	s, err := appsettings.NewAppSettings(filepath.Join(folders[0].Path, "settings.json"))
	if err != nil {
		return nil, err
	}

	return &Config{
		settings: s,
	}, nil
}

func (c *Config) GetToken() string {
	s, err := c.settings.GetString(keyToken)
	if err != nil {
		return ""
	}

	return s
}

func (c *Config) SetToken(s string) {
	c.settings.SetString(keyToken, s)
	c.settings.Persist()
}
