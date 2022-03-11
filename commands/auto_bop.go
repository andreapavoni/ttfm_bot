package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func AutoBopCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {

	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if len(args) == 0 {
		return &ttfm.CommandOutput{Msg: currentAutoBopStatusMsg(b.Config.AutoBop), User: user, ReplyWith: "room", Err: nil}
	}

	switch args[0] {
	case "on":
		return &ttfm.CommandOutput{Msg: enableAutoBop(b), User: user, ReplyWith: "action", Err: nil}
	case "off":
		return &ttfm.CommandOutput{Msg: disableAutoBop(b), User: user, ReplyWith: "action", Err: nil}
	default:
		return &ttfm.CommandOutput{Msg: currentAutoBopStatusMsg(b.Config.AutoBop), User: user, ReplyWith: "room", Err: nil}
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
	if !b.Config.AutoBop {
		b.ToggleAutoBop()
		return "/me enabled auto bop mode"
	}
	return "/me has already enabled auto bop mode"

}

func disableAutoBop(b *ttfm.Bot) string {
	if b.Config.AutoBop {
		b.ToggleAutoBop()

		return "/me disabled auto bop mode"
	}
	return "/me has already disabled auto bop mode"
}
