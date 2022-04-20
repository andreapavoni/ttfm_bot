package commands

import (
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func ReactionCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleNone},
		Help:               "React with a gif. Without args shows available reactions. Use `r add <img url> <reaction name>` to add a reaction",
		Handler:            reactionCommandHandler,
	}
}

func reactionCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if len(cmd.Args) == 0 {
		msg := fmt.Sprintf("Available reactions: %s", strings.Join(b.Reactions.Availables(), ", "))
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}

	if len(cmd.Args) == 3 && cmd.Args[0] == "add" {
		if err := b.Reactions.Put(cmd.Args[1], cmd.Args[2]); err != nil {
			return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: err}
		}
		msg := fmt.Sprintf("Added %s to `%s` reaction", cmd.Args[2], cmd.Args[1])
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}

	if reaction := b.Reactions.Get(cmd.Args[0]); len(cmd.Args) == 1 && reaction != "" {
		return &ttfm.CommandOutput{Msg: reaction, User: user, ReplyType: ttfm.MessageTypeRoom}
	}

	err := fmt.Errorf("I can't find any reaction for `%s`. Type `%sr` to know the available reactions", cmd.Args[0], b.Config.CmdPrefix)
	return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: err}
}
