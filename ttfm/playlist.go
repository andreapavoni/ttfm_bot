package ttfm

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type SongItem struct {
	id     string
	title  string
	artist string
	length int
}

type Playlist struct {
	Name  string
	songs *collections.SmartList[SongItem]
}

func NewPlaylist(playlistName string) *Playlist {
	return &Playlist{
		Name:  playlistName,
		songs: collections.NewSmartList[SongItem](),
	}
}

func (p *Playlist) GetSongById(songId string) (*SongItem, int, error) {
	s, idx := p.songs.Find(func(song *SongItem) bool { return songId == song.id })
	if idx >= 0 {
		return s, idx, nil
	}
	return s, idx, errors.New("I can't find the song in the current plalist")
}

func (p *Playlist) AddSong(s *SongItem) {
	p.songs.Push(*s)
}

func (p *Playlist) RemoveSong(s *SongItem) error {
	return p.songs.Remove(*s)
}
