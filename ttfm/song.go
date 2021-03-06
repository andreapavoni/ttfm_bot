package ttfm

import (
	"time"

	"github.com/andreapavoni/ttfm_bot/utils"
	"github.com/sirupsen/logrus"
)

type Song struct {
	Id        string
	DjName    string
	DjId      string
	Title     string
	Artist    string
	Length    int
	up        int
	down      int
	snag      int
	bot       *Bot
	skipTimer *time.Timer
}

func NewSong(bot *Bot) *Song {
	s := &Song{bot: bot}
	s.ResetSkipTimer()
	return s
}

func (s *Song) UpdateStats(up, down, snag int) {
	s.up = up
	s.down = down
	s.snag = snag
}

func (s *Song) Reset(id, title, artist string, length int, djName, djId string) {
	s.Id = id
	s.DjName = djName
	s.DjId = djId
	s.Title = title
	s.Artist = artist
	s.Length = length
	s.up = 0
	s.down = 0
	s.snag = 0
	s.ResetSkipTimer()
}

func (s *Song) UnpackVotelog(votelog [][]string) (userId, vote string) {
	if len(votelog) >= 1 && len(votelog[0]) >= 2 {
		return votelog[0][0], votelog[0][1]
	}
	logrus.WithField("votelog", votelog).Warn("Cannot parse Votelog")
	return "", ""
}

// Bop
func (s *Song) Bop() {
	if s.bot.Room.CurrentSong.DjId != s.bot.Identity.Id {
		s.bot.api.Bop()
	}
}

// Downvote
func (s *Song) Downvote() {
	if s.bot.Room.CurrentSong.DjId != s.bot.Identity.Id {
		s.bot.api.VoteDown()
	}
}

// SkipSong current song (must be moderator to skip others songs)
func (s *Song) Skip() {
	s.bot.api.Skip()
}

func (s *Song) StopSkipTimer() {
	if s.skipTimer != nil {
		s.skipTimer.Stop()
	}
}

func (s *Song) ResetSkipTimer() {
	if s.skipTimer != nil {
		s.skipTimer.Stop()
	}

	duration := utils.MinutesToDuration(int(s.bot.Config.MaxSongDuration))
	s.skipTimer = utils.ExecuteDelayed(duration, func() {
		s.Skip()
	})
}
