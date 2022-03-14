package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func SkipCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	b.SkipSong()
	return &ttfm.CommandOutput{Msg: "/me skipped song", User: user, ReplyType: cmd.Source}
}
