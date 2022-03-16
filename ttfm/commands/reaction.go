package commands

import (
	"fmt"
	"sort"
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
		availableReactions := b.Reactions.Availables.List()
		reactionsList := make([]string, len(availableReactions))
		copy(reactionsList, availableReactions)
		sort.Strings(reactionsList)

		msg := fmt.Sprintf("Available reactions: %s", strings.Join(reactionsList, ", "))
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}

	if reaction := b.Reactions.Get(cmd.Args[0]); reaction != "" {
		return &ttfm.CommandOutput{Msg: reaction, User: user, ReplyType: ttfm.MessageTypeRoom}
	}

	err := fmt.Errorf("I can't find any reaction for `%s`. Type !r to know the available reactions", cmd.Args[0])
	return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: err}
}
