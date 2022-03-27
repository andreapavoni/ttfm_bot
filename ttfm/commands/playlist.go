package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Handle playlists",
		Handler:            playlistCommandHandler,
	}
}

func playlistCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	if len(cmd.Args) >= 1 && cmd.Args[0] == "list" {
		return playlistListCommandHandler(b, cmd)
	}

	if len(cmd.Args) >= 2 && cmd.Args[0] == "add" {
		return playlistAddCommandHandler(b, cmd)
	}

	if len(cmd.Args) >= 2 && cmd.Args[0] == "switch" {
		return playlistSwitchCommandHandler(b, cmd)
	}

	if len(cmd.Args) >= 2 && cmd.Args[0] == "rm" {
		return playlistDeleteCommandHandler(b, cmd)
	}

	if len(cmd.Args) >= 1 && cmd.Args[0] == "rmsong" {
		return playlistRemoveSongCommandHandler(b, cmd)
	}

	user, _ := b.Users.UserFromId(cmd.UserId)
	msg := "Available playlist commands: add, list, rm, rmsong, switch"
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func playlistAddCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	if len(cmd.Args) < 2 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify a name of the new playlist")}
	}

	playlistName := strings.Join(cmd.Args[1:], " ")
	if err := b.Playlists.Add(playlistName); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I was unable to add the new playlist: %s", err.Error())}
	}

	msg := fmt.Sprintf("/me created playlist `%s`", playlistName)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func playlistListCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	pls := []string{}

	for _, pl := range b.Playlists.List() {
		if pl == b.Config.CurrentPlaylist {
			pl = "[" + pl + "]"
		}
		pls = append(pls, pl)
	}

	msg := "Available playists (the current one is highlighted): " + strings.Join(pls, ", ")
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func playlistSwitchCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	if len(cmd.Args) < 2 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify a name of the playlist you want to switch to")}
	}

	playlistName := strings.Join(cmd.Args[1:], " ")
	if err := b.Playlists.Switch(playlistName); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I was unable to switch playlist: %s", err.Error())}
	}

	msg := fmt.Sprintf("/me switched to playlist `%s`", playlistName)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func playlistDeleteCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if len(cmd.Args) < 2 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the name of the playlist you want to delete")}
	}

	playlistName := strings.Join(cmd.Args[1:], " ")
	if err := b.Playlists.Remove(playlistName); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I was unable to delete the playlist: %s", err.Error())}
	}

	msg := fmt.Sprintf("/me deleted playlist `%s`", playlistName)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func playlistRemoveSongCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	if err := b.CurrentPlaylist.RemoveSong(b.Room.Song.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I was unable to delete the playlist: %s", err.Error())}
	}

	return &ttfm.CommandOutput{Msg: "/me removed the song from the current playlist", User: user, ReplyType: cmd.Source}
}
