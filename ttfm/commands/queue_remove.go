package commands

import (
	"errors"
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func QueueRemoveCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleBotModerator},
		Help:               "Ask to be removed from the queue",
		Handler:            queueRemoveCommandHandler,
	}
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
