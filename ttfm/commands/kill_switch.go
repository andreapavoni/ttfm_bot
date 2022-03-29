package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func KillSwitchCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Kill switch the bot",
		Handler:            killSwitchCommandHandler}
}

func killSwitchCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {

	b.Actions.KillSwitch()
	return &ttfm.CommandOutput{ReplyType: ttfm.MessageTypeNone}
}
