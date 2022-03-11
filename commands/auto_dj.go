package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func AutoDjCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: disableAutoDj(b), User: user, ReplyWith: "pm", Err: err}
	}

	if len(args) == 0 {
		return &ttfm.CommandOutput{Msg: currentAutoDjStatusMsg(b.Config.AutoDj), User: user, ReplyWith: "room", Err: nil}
	}

	switch args[0] {
	case "on":
		return &ttfm.CommandOutput{Msg: enableAutoDj(b), User: user, ReplyWith: "action", Err: nil}
	case "off":
		return &ttfm.CommandOutput{Msg: disableAutoDj(b), User: user, ReplyWith: "action", Err: nil}
	default:
		return &ttfm.CommandOutput{Msg: currentAutoDjStatusMsg(b.Config.AutoDj), User: user, ReplyWith: "room", Err: nil}
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
		b.ToggleAutoDj()
		return "/me enabled auto dj mode"
	}
	return "/me has already enabled auto dj mode"
}

func disableAutoDj(b *ttfm.Bot) string {
	if b.Config.AutoDj {
		b.ToggleAutoDj()

		return "/me disabled auto dj mode"
	}
	return "/me has already disabled auto dj mode"

}
