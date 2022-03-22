package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func RoomCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Handle favourite rooms and join one",
		Handler:            roomCommandHandler,
	}
}

func roomCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	if len(cmd.Args) == 1 && cmd.Args[0] == "list" {
		return listRoomsCommandHandler(b, cmd)
	}

	if len(cmd.Args) == 3 {
		switch cmd.Args[0] {
		case "add":
			return addRoomCommandHandler(b, cmd)
		// case "rm":
		// 	return removeRoomCommandHandler(b, cmd)
		case "join":
			return joinRoomCommandHandler(b, cmd)
		}
	}
	user, _ := b.UserFromId(cmd.UserId)
	msg := "Available commands: add, remove, list, join"
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func addRoomCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)
	if len(cmd.Args) < 3 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the slug and Id of the room you want to add")}
	}

	roomSlug := cmd.Args[1]
	roomId := cmd.Args[2]

	b.Rooms.AddFavorite(roomSlug, roomId)
	msg := fmt.Sprintf("/me added `%s` to the list of rooms", roomSlug)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func listRoomsCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)
	rooms := strings.Join(b.Rooms.Keys(), ", ")
	msg := fmt.Sprintf("List of favourite rooms: %s", rooms)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func joinRoomCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)
	if len(cmd.Args) < 2 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the slug of the room you want to join")}
	}
	roomSlug := cmd.Args[1]
	if err := b.Rooms.Join(roomSlug); err == nil {
		msg := fmt.Sprintf("/me is joining `%s`", roomSlug)
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	} else {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: fmt.Errorf("Error joining room `%s`: %s", roomSlug, err.Error())}
	}
}
