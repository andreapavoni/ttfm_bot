package commands

import "github.com/andreapavoni/ttfm_bot/ttfm"

func BopCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Bop current song",
		Handler:            bopCommandHandler,
	}
}

func bopCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	b.Room.Song.Bop()
	return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypeNone}
}
