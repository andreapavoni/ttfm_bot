package ttfm

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/utils"
)

type Config struct {
	ApiAuth                  string
	UserId                   string
	RoomId                   string
	MainAdminId              string
	AutoSnagEnabled          bool
	AutoBopEnabled           bool
	AutoDjEnabled            bool
	AutoDjMinDjs             int64
	AutoShowSongStatsEnabled bool
	AutoShowDjStatsEnabled   bool
	AutoWelcomeEnabled       bool
	QueueEnabled             bool
	QueueInviteDuration      int64
	MaxSongDuration          int64
	MaxSongsPerDj            int64
	CurrentPlaylist          string
	SetBot                   bool
	brain                    *Brain
	MusicTheme               string
}

func NewConfig(b *Brain) *Config {
	c := &Config{brain: b}
	if err := c.loadConfigFromDb(); err != nil {
		c.loadDefaultConfig()
	}
	return c
}

func (c *Config) Save() error {
	return c.brain.Put("config", c)
}

// EnableAutoSnag each song
func (c *Config) EnableAutoSnag(status bool) bool {
	c.AutoSnagEnabled = status
	c.Save()
	return status
}

// EnableAutoBop each song
func (c *Config) EnableAutoBop(status bool) bool {
	c.AutoBopEnabled = status
	c.Save()
	return status
}

// EnableAutoDj enabled/disabled. If bot is djing, it will be escorted when song is finished
func (c *Config) EnableAutoDj(status bool) bool {
	c.AutoDjEnabled = status
	c.Save()
	return status
}

// EnableQueue enabled/disabled
func (c *Config) EnableQueue(status bool) bool {
	c.QueueEnabled = status
	c.Save()
	return status
}

func (c *Config) loadConfigFromDb() error {
	if err := c.brain.Get("config", &c); err != nil {
		return errors.New("config not found")
	}

	return nil
}

func (c *Config) loadDefaultConfig() {
	c.ApiAuth = utils.GetEnvOrPanic("TTFM_API_AUTH")
	c.UserId = utils.GetEnvOrPanic("TTFM_API_USER_ID")
	c.MainAdminId = utils.GetEnvOrPanic("TTFM_MAIN_ADMIN_ID")
	c.RoomId = utils.GetEnvOrPanic("TTFM_API_ROOM_ID")
	c.AutoSnagEnabled = false
	c.AutoBopEnabled = true
	c.AutoDjEnabled = false
	c.AutoDjMinDjs = 0
	c.AutoShowSongStatsEnabled = false
	c.AutoShowDjStatsEnabled = false
	c.AutoWelcomeEnabled = false
	c.QueueEnabled = false
	c.QueueInviteDuration = 1
	c.MaxSongDuration = 10
	c.MaxSongsPerDj = 0
	c.CurrentPlaylist = "default"
	c.SetBot = false
	c.MusicTheme = "FREE PLAY"
}
