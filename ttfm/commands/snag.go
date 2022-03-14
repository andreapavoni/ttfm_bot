package commands

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func SnagCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if err := b.Snag(); err == nil {
		return &ttfm.CommandOutput{Msg: "/me snagged this song", User: user, ReplyType: cmd.Source}
	}

	return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I've failed to snag this song")}
}
