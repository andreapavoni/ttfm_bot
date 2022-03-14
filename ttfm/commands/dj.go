package commands

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func DjCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if b.UserIsDj(b.Config.UserId) {
		if b.UserIsCurrentDj(b.Config.UserId) {
			if err := b.AddDjEscorting(b.Config.UserId); err != nil {
				return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I've was unable to prepare myself for escorting : " + err.Error())}
			}
			return &ttfm.CommandOutput{Msg: "/me will get off the stage at the end of current song", User: user, ReplyType: cmd.Source}
		}

		if err := b.EscortDj(b.Config.UserId); err == nil {
			return &ttfm.CommandOutput{Msg: "/me has left the stage", User: user, ReplyType: cmd.Source}
		} else {
			return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
		}
	}

	b.Dj()
	return &ttfm.CommandOutput{Msg: "/me is going on stage", User: user, ReplyType: cmd.Source}
}
