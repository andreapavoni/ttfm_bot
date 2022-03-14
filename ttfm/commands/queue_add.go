package commands

import (
	"errors"
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func QueueAddCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if !b.Config.ModQueue {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("Queue mode is disabled")}
	}

	if err := b.Queue.Add(user.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	msg := fmt.Sprintf("/me put %s in the queue with position #%d", user.Name, b.Queue.IndexOf(user.Id)+1)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: ttfm.MessageTypeRoom}
}
