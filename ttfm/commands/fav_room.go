package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func FavRoomCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Handle favourite rooms and join one. Usage: `room <list | current | join slug | add slug room_id>`",
		Handler:            favRoomCommandHandler,
	}
}

func favRoomCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	if len(cmd.Args) == 1 && cmd.Args[0] == "list" {
		return listFavRoomsCommandHandler(b, cmd)
	}

	if len(cmd.Args) == 1 && cmd.Args[0] == "current" {
		return showCurrentRoomCommandHandler(b, cmd)
	}

	if len(cmd.Args) == 2 && cmd.Args[0] == "join" {
		return joinFavRoomCommandHandler(b, cmd)
	}

	if len(cmd.Args) == 3 && cmd.Args[0] == "add" {
		return addFavRoomCommandHandler(b, cmd)
	}

	user, _ := b.Users.UserFromId(cmd.UserId)
	msg := "Available favorite room commands: add, list, join"
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func addFavRoomCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	if len(cmd.Args) < 3 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the slug and Id of the room you want to add")}
	}

	roomSlug := cmd.Args[1]
	roomId := cmd.Args[2]

	b.FavRooms.AddFavorite(roomSlug, roomId)
	msg := fmt.Sprintf("/me added `%s` to the list of rooms", roomSlug)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func listFavRoomsCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	rooms := strings.Join(b.FavRooms.Keys(), ", ")
	msg := fmt.Sprintf("List of favourite rooms: %s", rooms)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func showCurrentRoomCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	msg := fmt.Sprintf("Current room: %s (ID: %s)", b.Room.Name, b.Room.Id)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func joinFavRoomCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	if len(cmd.Args) < 2 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the slug of the room you want to join")}
	}
	roomSlug := cmd.Args[1]
	if err := b.FavRooms.Join(roomSlug); err == nil {
		msg := fmt.Sprintf("/me is joining `%s`", roomSlug)
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	} else {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: fmt.Errorf("Error joining room `%s`: %s", roomSlug, err.Error())}
	}
}
