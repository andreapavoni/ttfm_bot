package commands

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func SnagCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "none", Err: err}
	}

	if err := b.Snag(); err == nil {
		return &ttfm.CommandOutput{Msg: "/me snagged this song", User: user, ReplyWith: "action", Err: nil}
	}

	return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("I've failed to snag this song")}
}
