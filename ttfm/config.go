package ttfm

import "github.com/andreapavoni/ttfm_bot/utils"

type Config struct {
	ApiAuth                string
	UserId                 string
	RoomId                 string
	Admins                 []string
	AutoSnag               bool
	AutoBop                bool
	AutoDj                 bool
	AutoShowSongStats      bool
	ModAutoWelcome         bool
	ModQueue               bool
	ModQueueInviteDuration int64
	ModSongsMaxDuration    int64
	ModSongsMaxPerDj       int64
	ModDjRestDuration      int64
	CurrentPlaylist        string
	SetBot                 bool
}

func LoadConfigFromEnvs() *Config {
	return &Config{
		ApiAuth:                utils.GetEnvOrPanic("TTFM_API_AUTH"),
		UserId:                 utils.GetEnvOrPanic("TTFM_API_USER_ID"),
		Admins:                 utils.StringToSlice(utils.GetEnvOrDefault("TTFM_ADMINS", "pavonz"), ","),
		RoomId:                 utils.GetEnvOrPanic("TTFM_API_ROOM_ID"),
		AutoSnag:               utils.GetEnvBoolOrDefault("TTFM_AUTO_SNAG", true),
		AutoBop:                utils.GetEnvBoolOrDefault("TTFM_AUTO_BOP", true),
		AutoDj:                 utils.GetEnvBoolOrDefault("TTFM_AUTO_DJ", false),
		AutoShowSongStats:      utils.GetEnvBoolOrDefault("TTFM_AUTO_SHOW_SONG_STATS", false),
		ModAutoWelcome:         utils.GetEnvBoolOrDefault("TTFM_AUTO_WELCOME", false),
		ModQueue:               utils.GetEnvBoolOrDefault("TTFM_MOD_QUEUE", false),
		ModQueueInviteDuration: utils.GetEnvIntOrDefault("TTFM_MOD_QUEUE_INVITE_DURATION", 1),
		ModSongsMaxDuration:    utils.GetEnvIntOrDefault("TTFM_MOD_SONGS_MAX_DURATION", 10),
		ModSongsMaxPerDj:       utils.GetEnvIntOrDefault("TTFM_MOD_SONGS_MAX_PER_DJ", 0),
		ModDjRestDuration:      utils.GetEnvIntOrDefault("TTFM_MOD_DJ_REST_DURATION", 0),
		CurrentPlaylist:        utils.GetEnvOrDefault("TTFM_DEFAULT_PLAYLIST", "default"),
		SetBot:                 utils.GetEnvBoolOrDefault("TTFM_SET_BOT", false),
	}
}
