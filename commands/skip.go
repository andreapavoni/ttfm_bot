package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func SkipCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	b.SkipSong()
	return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "none", Err: nil}
}
