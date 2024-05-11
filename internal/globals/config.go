package globals

import (
	"path/filepath"
	"webchat/internal/utils"
)

var Config utils.Config

var ConfigFilepath, _ = filepath.Abs("./config/main.yaml")

func init() {
	Config = utils.ParseConfig(ConfigFilepath)
}
