package commands

import (
	"errors"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func SayCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Say something in the room",
		Handler:            sayCommandHandler,
	}
}

func sayCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("You must specify the message you want me to say")}
	}
	return &ttfm.CommandOutput{Msg: strings.Join(cmd.Args, " "), User: user, ReplyType: ttfm.MessageTypeRoom}
}
