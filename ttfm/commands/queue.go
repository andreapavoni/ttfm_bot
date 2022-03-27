package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func QueueCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin, ttfm.UserRoleBotModerator},
		Help:               "Add, remove yourself from queue, or just check the current status",
		Handler:            queueCommandHandler,
	}
}

func queueCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	if len(cmd.Args) == 0 {
		return &ttfm.CommandOutput{Msg: currentQueueStatusMsg(b), User: user, ReplyType: cmd.Source}
	}

	if len(cmd.Args) > 0 && cmd.Args[0] == "add" {
		return queueAddCommandHandler(b, cmd)
	}

	if len(cmd.Args) > 0 && cmd.Args[0] == "rm" {
		return queueRemoveCommandHandler(b, cmd)
	}

	return &ttfm.CommandOutput{Msg: currentQueueStatusMsg(b), User: user, ReplyType: cmd.Source}
}

func queueAddCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	if !b.Config.QueueEnabled {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("Queue mode is disabled")}
	}

	if err := b.Queue.Add(user.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	msg := fmt.Sprintf("/me put %s in the queue with position #%d", user.Name, b.Queue.IndexOf(user.Id)+1)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: ttfm.MessageTypeRoom}
}

func queueRemoveCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	if !b.Config.QueueEnabled {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("Queue mode is disabled")}
	}

	if err := b.Queue.Remove(user.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	msg := fmt.Sprintf("/me removed %s from the queue", user.Name)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: ttfm.MessageTypeRoom}
}

func currentQueueStatusMsg(b *ttfm.Bot) string {
	if !b.Config.QueueEnabled {
		return "Queue mode is disabled"
	}

	if b.Queue.Size() > 0 {
		return "Queue mode is enabled. People in line: " + strings.Join(b.Queue.List(), ", ")
	} else {
		return "Queue mode is enabled. No people in line"
	}
}
