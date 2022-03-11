package commands

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func DjCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "", Err: err}
	}

	if b.UserIsDj(b.Config.UserId) {
		if b.UserIsCurrentDj(b.Config.UserId) {
			if err := b.AddDjEscorting(b.Config.UserId); err != nil {
				return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "", Err: errors.New("I've had some problem in getting into the escorting list: " + err.Error())}
			}
			return &ttfm.CommandOutput{Msg: "/me will get off the stage at the end of current song", User: user, ReplyWith: "action", Err: nil}
		}

		b.EscortDj(b.Config.UserId)
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "", Err: nil}
	}

	b.AutoDj()
	return &ttfm.CommandOutput{Msg: "/me is going on stage", User: user, ReplyWith: "action", Err: nil}
}
