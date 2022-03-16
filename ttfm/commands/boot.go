package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func BootCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin, ttfm.UserRoleBotModerator},
		Help:               "Boots a user out of the room",
		Handler:            bootCommandHandler,
	}
}

func bootCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the username of the user you kick")}
	}

	kickedUser, err := b.UserFromName(strings.Join(cmd.Args, " "))
	reason := fmt.Sprintf("Ask @%s for details", user.Name)

	if err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I can't find the user you want to kick")}
	}

	if err := b.BootUser(kickedUser.Id, reason); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I wasn't able to boot @%s", kickedUser.Name)}
	}

	return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypeNone}
}
