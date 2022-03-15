package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func SkipCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin, ttfm.UserRoleBotModerator},
		Help:               "Skip current song",
		Handler:            skipCommandHandler,
	}
}

func skipCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)
	b.SkipSong()
	return &ttfm.CommandOutput{Msg: "/me skipped song", User: user, ReplyType: cmd.Source}
}
