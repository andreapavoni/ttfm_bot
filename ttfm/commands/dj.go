package commands

import (
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func DjCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Jump on/off the stage",
		Handler:            djCommandHandler,
	}
}

func djCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if b.Users.UserIsDj(b.Identity.Id) {
		if b.Users.UserIsCurrentDj(b.Identity.Id) {
			if err := b.Room.AddDjEscorting(b.Identity.Id); err != nil {
				return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I've was unable to prepare myself for escorting : %s", err.Error())}
			}
			return &ttfm.CommandOutput{Msg: "/me will get off the stage at the end of current song", User: user, ReplyType: cmd.Source}
		}

		if err := b.Users.EscortDj(b.Identity.Id); err == nil {
			return &ttfm.CommandOutput{Msg: "/me has left the stage", User: user, ReplyType: cmd.Source}
		} else {
			return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
		}
	}

	b.Actions.AutoDj()
	return &ttfm.CommandOutput{Msg: "/me is going on stage", User: user, ReplyType: cmd.Source}
}
