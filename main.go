package main

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
	"github.com/andreapavoni/ttfm_bot/utils"
)

func main() {
	ttfm.Start(ttfm.Config{
		ApiAuth:             utils.GetEnvOrPanic("TTFM_API_AUTH"),
		UserId:              utils.GetEnvOrPanic("TTFM_API_USER_ID"),
		Admins:              utils.StringToSlice(utils.GetEnvOrDefault("TTFM_ADMINS", "pavonz"), ","),
		RoomId:              utils.GetEnvOrPanic("TTFM_API_ROOM_ID"),
		AutoSnag:            utils.GetEnvBoolOrDefault("TTFM_AUTO_SNAG", true),
		AutoBop:             utils.GetEnvBoolOrDefault("TTFM_AUTO_BOP", true),
		AutoDj:              utils.GetEnvBoolOrDefault("TTFM_AUTO_DJ", false),
		AutoQueue:           utils.GetEnvBoolOrDefault("TTFM_AUTO_QUEUE", false),
		AutoQueueMsg:        utils.GetEnvOrDefault("TTFM_AUTO_QUEUE_MSG", "Next in line is: @JumpingMonkey. Time to claim your spot!"),
		AutoShowSongStats:   utils.GetEnvBoolOrDefault("TTFM_AUTO_SHOW_SONG_STATS", false),
		ModAutoWelcome:      utils.GetEnvBoolOrDefault("TTFM_AUTO_WELCOME", false),
		ModQueue:            utils.GetEnvBoolOrDefault("TTFM_MOD_QUEUE", false),
		ModSongsMaxDuration: utils.GetEnvIntOrDefault("TTFM_MOD_SONGS_MAX_DURATION", 10),
		ModSongsMaxPerDj:    utils.GetEnvIntOrDefault("TTFM_MOD_SONGS_MAX_DURATION", 0),
		DefaultPlaylist:     utils.GetEnvOrDefault("TTFM_DEFAULT_PLAYLIST", "default"),
		SetBot:              utils.GetEnvBoolOrDefault("TTFM_SET_BOT", false),
	})
}
