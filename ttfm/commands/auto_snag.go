package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func AutoSnagCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if len(cmd.Args) == 0 {
		return &ttfm.CommandOutput{Msg: currentAutoSnagStatusMsg(b.Config.AutoSnag), User: user, ReplyType: cmd.Source}
	}

	switch cmd.Args[0] {
	case "on":
		return &ttfm.CommandOutput{Msg: enableAutoSnag(b), User: user, ReplyType: cmd.Source}
	case "off":
		return &ttfm.CommandOutput{Msg: disableAutoSnag(b), User: user, ReplyType: cmd.Source}
	default:
		return &ttfm.CommandOutput{Msg: currentAutoSnagStatusMsg(b.Config.AutoSnag), User: user, ReplyType: cmd.Source}
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
	if !b.Config.AutoSnag {
		b.AutoSnag(true)

		return "/me enabled auto snag mode"
	}
	return "/me has already enabled auto snag mode"

}

func disableAutoSnag(b *ttfm.Bot) string {
	if b.Config.AutoSnag {
		b.AutoSnag(false)

		return "/me disabled auto snag mode"
	}
	return "/me has already disabled auto snag mode"
}
