package utils

import "github.com/youngduc/go-blog/config"

func IsProduction() bool {
	return !config.Cfg.App.Debug
}

func IsDebug() bool {
	return config.Cfg.App.Debug
}
