package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func AutoSnagCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "room", Err: err}
	}

	if len(args) == 0 {
		return &ttfm.CommandOutput{Msg: currentAutoSnagStatusMsg(b.Config.AutoSnag), User: user, ReplyWith: "room", Err: nil}
	}

	switch args[0] {
	case "on":
		return &ttfm.CommandOutput{Msg: enableAutoSnag(b), User: user, ReplyWith: "action", Err: nil}
	case "off":
		return &ttfm.CommandOutput{Msg: disableAutoSnag(b), User: user, ReplyWith: "action", Err: nil}
	default:
		return &ttfm.CommandOutput{Msg: currentAutoSnagStatusMsg(b.Config.AutoSnag), User: user, ReplyWith: "room", Err: nil}
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
		b.ToggleAutoSnag()

		return "/me enabled auto snag mode"
	}
	return "/me has already enabled auto snag mode"

}

func disableAutoSnag(b *ttfm.Bot) string {
	if b.Config.AutoSnag {
		b.ToggleAutoSnag()

		return "/me disabled auto snag mode"
	}
	return "/me has already disabled auto snag mode"
}
