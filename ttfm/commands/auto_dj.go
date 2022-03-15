package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func AutoDjCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Enables/disables auto dj mode. Without args prints current setting",
		Handler:            autoDjCommandHandler,
	}
}

func autoDjCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if len(cmd.Args) == 0 {
		return &ttfm.CommandOutput{Msg: currentAutoDjStatusMsg(b.Config.AutoDj), User: user, ReplyType: cmd.Source}
	}

	switch cmd.Args[0] {
	case "on":
		return &ttfm.CommandOutput{Msg: enableAutoDj(b), User: user, ReplyType: cmd.Source}
	case "off":
		return &ttfm.CommandOutput{Msg: disableAutoDj(b), User: user, ReplyType: cmd.Source}
	default:
		return &ttfm.CommandOutput{Msg: currentAutoDjStatusMsg(b.Config.AutoDj), User: user, ReplyType: cmd.Source}
	}
}

func currentAutoDjStatusMsg(status bool) string {
	if status {
		return "Auto DJ mode is enabled"
	} else {
		return "Auto DJ mode is disabled"
	}
}

func enableAutoDj(b *ttfm.Bot) string {
	if !b.Config.AutoDj {
		b.AutoDj(true)
		return "/me enabled auto dj mode"
	}
	return "/me has already enabled auto dj mode"
}

func disableAutoDj(b *ttfm.Bot) string {
	if b.Config.AutoDj {
		b.AutoDj(false)
		return "/me disabled auto dj mode"
	}
	return "/me has already disabled auto dj mode"

}
