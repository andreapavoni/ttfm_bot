package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func EscortMeCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	b.AddDjEscorting(userId)
	return &ttfm.CommandOutput{Msg: "I'm going to escort you after your next song has been played", User: user, ReplyWith: "pm", Err: nil}
}
