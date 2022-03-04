package ttfm

import "github.com/sirupsen/logrus"

type Song struct {
	id     string
	djName string
	djId   string
	title  string
	artist string
	length int
	up     int
	down   int
	snag   int
}

func (s *Song) UpdateStats(up, down, snag int) {
	s.up = up
	s.down = down
	s.snag = snag
}

func (s *Song) Reset(id, title, artist string, length int, djName, djId string) {
	s.id = id
	s.djName = djName
	s.djId = djId
	s.title = title
	s.artist = artist
	s.length = length
	s.up = 0
	s.down = 0
	s.snag = 0
}

func (s *Song) UnpackVotelog(votelog [][]string) (string, string) {
	if len(votelog) < 1 && len(votelog[0]) < 2 {
		logrus.WithField("votelog", votelog).Warn("Cannot parse Votelog")
		return "", ""
	}
	// if len(votelog[0]) < 2 {
	// 	logrus.WithField("votelog", votelog).Warn("Cannot parse Votelog")
	// 	return "", ""
	// }

	userID := votelog[0][0]
	vote := votelog[0][1]

	return userID, vote
}
