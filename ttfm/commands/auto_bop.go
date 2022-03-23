package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func AutoBopCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Enables/disables auto bop mode. Without args shows current setting",
		Handler:            autoBopCommandHandler,
	}
}

func autoBopCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if len(cmd.Args) == 0 {
		return &ttfm.CommandOutput{Msg: currentAutoBopStatusMsg(b.Config.AutoBopEnabled), User: user, ReplyType: cmd.Source}
	}

	switch cmd.Args[0] {
	case "on":
		return &ttfm.CommandOutput{Msg: enableAutoBop(b), User: user, ReplyType: cmd.Source}
	case "off":
		return &ttfm.CommandOutput{Msg: disableAutoBop(b), User: user, ReplyType: cmd.Source}
	default:
		return &ttfm.CommandOutput{Msg: currentAutoBopStatusMsg(b.Config.AutoBopEnabled), User: user, ReplyType: cmd.Source}
	}
}

func currentAutoBopStatusMsg(status bool) string {
	if status {
		return "Auto bop mode is enabled"
	} else {
		return "Auto bop mode is disabled"
	}
}

func enableAutoBop(b *ttfm.Bot) string {
	if !b.Config.AutoBopEnabled {
		b.Config.EnableAutoBop(true)
		return "/me enabled auto bop mode"
	}
	return "/me has already enabled auto bop mode"

}

func disableAutoBop(b *ttfm.Bot) string {
	if b.Config.AutoBopEnabled {
		b.Config.EnableAutoBop(false)

		return "/me disabled auto bop mode"
	}
	return "/me has already disabled auto bop mode"
}
