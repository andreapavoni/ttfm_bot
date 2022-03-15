package commands

import (
	"errors"
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func QueueAddCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleBotModerator},
		Help:               "Ask to get added into queue",
		Handler:            queueAddCommandHandler,
	}
}

func queueAddCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)
	if !b.Config.ModQueue {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("Queue mode is disabled")}
	}

	if err := b.Queue.Add(user.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	msg := fmt.Sprintf("/me put %s in the queue with position #%d", user.Name, b.Queue.IndexOf(user.Id)+1)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: ttfm.MessageTypeRoom}
}
