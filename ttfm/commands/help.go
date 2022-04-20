package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func HelpCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleNone},
		Help:               "Show commands help. Usage `help [command]`. Without arguments show the available commands to user",
		Handler:            helpCommandHandler,
	}
}

func helpCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if len(cmd.Args) == 1 {
		command, err := b.Commands.Get(cmd.Args[0])

		if err != nil {
			return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: err}
		}

		msg := fmt.Sprintf("%s - %s", cmd.Args[0], command.Help)
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}

	availableCmds := []string{}
	for _, c := range b.Commands.List() {
		command, err := b.Commands.Get(c)
		if err != nil {
			continue
		}

		if err := b.Users.CheckAuthorizations(user, command.AuthorizationRoles...); err == nil {
			availableCmds = append(availableCmds, c)
		}
	}

	sort.Strings(availableCmds)
	msg := fmt.Sprintf("@%s according to our respective roles here, you can run the following commands: ", user.Name) + strings.Join(availableCmds, ", ")
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
