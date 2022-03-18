package ttfm

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/utils"
)

type Config struct {
	ApiAuth                string
	UserId                 string
	RoomId                 string
	Admins                 []string
	AutoSnag               bool
	AutoBop                bool
	AutoDj                 bool
	AutoDjCountTrigger     int64
	AutoShowSongStats      bool
	ModAutoWelcome         bool
	ModQueue               bool
	ModQueueInviteDuration int64
	ModSongsMaxDuration    int64
	ModSongsMaxPerDj       int64
	ModDjRestDuration      int64
	CurrentPlaylist        string
	SetBot                 bool
	brain                  *Brain
}

func NewConfig(b *Brain) *Config {
	var cfg *Config
	var err error

	cfg, err = loadConfigFromDb(b)
	if err != nil {
		cfg = loadConfigFromEnvs()
	}
	cfg.brain = b

	return cfg
}

func (c *Config) Save() error {
	return c.brain.Put("config", "config", c)
}

func loadConfigFromDb(b *Brain) (*Config, error) {
	c := Config{}
	if err := b.Get("config", "config", &c); err != nil {
		return nil, errors.New("config not found")
	}

	return &c, nil
}

func loadConfigFromEnvs() *Config {
	return &Config{
		ApiAuth:                utils.GetEnvOrPanic("TTFM_API_AUTH"),
		UserId:                 utils.GetEnvOrPanic("TTFM_API_USER_ID"),
		Admins:                 utils.StringToSlice(utils.GetEnvOrDefault("TTFM_ADMINS", "pavonz"), ","),
		RoomId:                 utils.GetEnvOrPanic("TTFM_API_ROOM_ID"),
		AutoSnag:               utils.GetEnvBoolOrDefault("TTFM_AUTO_SNAG", true),
		AutoBop:                utils.GetEnvBoolOrDefault("TTFM_AUTO_BOP", true),
		AutoDj:                 utils.GetEnvBoolOrDefault("TTFM_AUTO_DJ", false),
		AutoDjCountTrigger:     utils.GetEnvIntOrDefault("TTFM_AUTO_DJ_COUNT_TRIGGER", 0),
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
