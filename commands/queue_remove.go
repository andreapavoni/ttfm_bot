package commands

import (
	"errors"
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func QueueRemoveCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if !b.Config.ModQueue {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("Queue mode is disabled")}
	}

	if err := b.Queue.Remove(user.Id); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	msg := fmt.Sprintf("/me removed %s from the queue", user.Name)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyWith: "pm", Err: nil}
}
