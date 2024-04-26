package util

import "poem-bot/global"

func GetOrDefault(configPath string, defaultStr string) string {
	config := global.VP.GetString(configPath)
	if config == "" {
		return defaultStr
	}
	return config
}
