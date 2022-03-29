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
	bot   *Bot
}

func NewPlaylist(b *Bot, playlistName string) *Playlist {
	return &Playlist{
		bot:   b,
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

func (p *Playlist) RemoveSong(songId string) error {
	song, idx, err := p.bot.CurrentPlaylist.GetSongById(songId)

	if err != nil {
		return err
	}

	if err := p.bot.api.PlaylistRemove(p.bot.Config.CurrentPlaylist, idx); err != nil {
		return err
	}

	return p.songs.Remove(*song)
}

// LoadPlaylist
func (p *Playlist) LoadPlaylist(playlistName string) error {
	playlist, err := p.bot.api.PlaylistAll(p.bot.Config.CurrentPlaylist)

	if err != nil {
		return err
	}

	for _, s := range playlist.List {
		p.AddSong(&SongItem{
			id:     s.ID,
			title:  s.Metadata.Song,
			artist: s.Metadata.Artist,
			length: s.Metadata.Length,
		})
	}
	return nil
}

// PushSongBottom
func (p *Playlist) PushSongBottom() error {
	if err := p.bot.api.PlaylistReorder(p.bot.Config.CurrentPlaylist, 0, p.songs.Size()-1); err == nil {
		currentSong, _ := p.songs.Shift()
		p.AddSong(&currentSong)
		return nil
	} else {
		return err
	}
}

// Snag current playing song into the current playlist
func (p *Playlist) Snag() error {
	if p.bot.Room.CurrentSong.DjId == p.bot.Identity.Id {
		return errors.New("I'm the current DJ and I already have this song in my playlist...")
	}

	playlist, err := p.bot.api.PlaylistAll(p.bot.Config.CurrentPlaylist)
	if err != nil {
		return err
	}

	p.bot.api.Snag()
	if err = p.bot.api.PlaylistAdd(p.bot.Room.CurrentSong.Id, p.bot.Config.CurrentPlaylist, len(playlist.List)); err != nil {
		return nil
	}

	p.AddSong(&SongItem{
		id:     p.bot.Room.CurrentSong.Id,
		title:  p.bot.Room.CurrentSong.Title,
		artist: p.bot.Room.CurrentSong.Artist,
		length: p.bot.Room.CurrentSong.Length,
	})

	return nil
}

type Playlists struct {
	bot *Bot
	*collections.SmartList[string]
}

func NewPlaylists(b *Bot) *Playlists {
	return &Playlists{bot: b, SmartList: collections.NewSmartList[string]()}
}

// LoadPlaylists from API and cache them in memory
func (p *Playlists) LoadPlaylists() error {
	playlists, err := p.bot.api.PlaylistListAll()
	if err != nil {
		return err
	}

	for _, pl := range playlists.List {
		p.SmartList.Push(pl.Name)
	}
	return nil
}

// AddPlaylist
func (p *Playlists) Add(playlistName string) error {
	if !p.SmartList.HasElement(playlistName) {
		if err := p.bot.api.PlaylistCreate(playlistName); err != nil {
			return err
		}
		p.SmartList.Push(playlistName)
		return nil
	}

	return errors.New("Playlist not found")
}

// RemovePlaylist
func (p *Playlists) Remove(playlistName string) error {
	if p.SmartList.HasElement(playlistName) {
		if err := p.bot.api.PlaylistDelete(playlistName); err != nil {
			return err
		}
		return p.LoadPlaylists()
	}
	return errors.New("Playlist not found")
}

// SwitchPlaylist
func (p *Playlists) Switch(playlistName string) error {
	if p.SmartList.HasElement(playlistName) {
		if err := p.bot.api.PlaylistSwitch(playlistName); err != nil {
			return err
		}
		p.bot.Config.CurrentPlaylist = playlistName
		p.bot.Config.Save()
		return p.bot.CurrentPlaylist.LoadPlaylist(playlistName)
	}

	return errors.New("Playlist not found")
}
