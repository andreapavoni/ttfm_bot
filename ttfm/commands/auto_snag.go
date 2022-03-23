package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func AutoSnagCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Enables/disables auto snag mode. Without args shows current setting",
		Handler:            autoSnagCommandHandler,
	}
}

func autoSnagCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if len(cmd.Args) == 0 {
		return &ttfm.CommandOutput{Msg: currentAutoSnagStatusMsg(b.Config.AutoSnagEnabled), User: user, ReplyType: cmd.Source}
	}

	switch cmd.Args[0] {
	case "on":
		return &ttfm.CommandOutput{Msg: enableAutoSnag(b), User: user, ReplyType: cmd.Source}
	case "off":
		return &ttfm.CommandOutput{Msg: disableAutoSnag(b), User: user, ReplyType: cmd.Source}
	default:
		return &ttfm.CommandOutput{Msg: currentAutoSnagStatusMsg(b.Config.AutoSnagEnabled), User: user, ReplyType: cmd.Source}
	}
}

func currentAutoSnagStatusMsg(status bool) string {
	if status {
		return "Auto snag mode is enabled"
	} else {
		return "Auto snag mode is disabled"
	}
}

func enableAutoSnag(b *ttfm.Bot) string {
	if !b.Config.AutoSnagEnabled {
		b.Config.EnableAutoSnag(true)

		return "/me enabled auto snag mode"
	}
	return "/me has already enabled auto snag mode"

}

func disableAutoSnag(b *ttfm.Bot) string {
	if b.Config.AutoSnagEnabled {
		b.Config.EnableAutoSnag(false)

		return "/me disabled auto snag mode"
	}
	return "/me has already disabled auto snag mode"
}
