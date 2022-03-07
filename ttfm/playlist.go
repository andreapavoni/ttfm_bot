package ttfm

import (
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

func (p *Playlist) List() *[]Playlist {
	return &[]Playlist{}
}

func (p *Playlist) AddSong(s *SongItem) {
	p.songs.Push(*s)
}

func (p *Playlist) RemoveSong(s *SongItem) error {
	return p.songs.Remove(*s)
}
