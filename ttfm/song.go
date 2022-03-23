package ttfm

import (
	"github.com/sirupsen/logrus"
)

type Song struct {
	Id     string
	DjName string
	DjId   string
	Title  string
	Artist string
	Length int
	up     int
	down   int
	snag   int
	bot    *Bot
}

func NewSong(bot *Bot) *Song {
	return &Song{bot: bot}
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
}

func (s *Song) UnpackVotelog(votelog [][]string) (userId, vote string) {
	if len(votelog) >= 1 && len(votelog[0]) >= 2 {
		logrus.WithField("votelog", votelog).Warn("Cannot parse Votelog")
		return "", ""
	}
	return votelog[0][0], votelog[0][1]
}

// Bop
func (s *Song) Bop() {
	if s.bot.Room.Song.DjId != s.bot.Identity.Id {
		s.bot.api.Bop()
	}
}

// Downvote
func (s *Song) Downvote() {
	if s.bot.Room.Song.DjId != s.bot.Identity.Id {
		s.bot.api.VoteDown()
	}
}

// SkipSong current song (must be moderator to skip others songs)
func (s *Song) Skip() {
	s.bot.api.Skip()
}
