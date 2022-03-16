package commands

import (
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func ReactionCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleNone},
		Help:               "React with a gif. Without args shows available reactions",
		Handler:            reactionCommandHandler,
	}
}

func reactionCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if len(cmd.Args) == 0 {
		msg := fmt.Sprintf("Available reactions: %s", strings.Join(b.Reactions.Availables.List(), ", "))
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}

	if reaction := b.Reactions.Get(cmd.Args[0]); reaction != "" {
		return &ttfm.CommandOutput{Msg: reaction, User: user, ReplyType: ttfm.MessageTypeRoom}
	}
	return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: fmt.Errorf("I didn't find any reaction for %s. Type !r to know the available reactions", cmd.Args[0])}
}
