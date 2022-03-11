package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func QueueCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if len(args) == 0 {
		return &ttfm.CommandOutput{Msg: currentQueueStatusMsg(b.Config.ModQueue), User: user, ReplyWith: "room", Err: nil}
	}

	switch args[0] {
	case "on":
		return &ttfm.CommandOutput{Msg: enableQueue(b), User: user, ReplyWith: "pm", Err: nil}
	case "off":
		return &ttfm.CommandOutput{Msg: disableQueue(b), User: user, ReplyWith: "pm", Err: nil}
	default:
		return &ttfm.CommandOutput{Msg: currentQueueStatusMsg(b.Config.ModQueue), User: user, ReplyWith: "room", Err: nil}
	}
}

func currentQueueStatusMsg(status bool) string {
	if status {
		return "Queue mode is enabled"
	} else {
		return "Queue mode is disabled"
	}
}

func enableQueue(b *ttfm.Bot) string {
	if !b.Config.ModQueue {
		b.ToggleModQueue()
		return "/me has enabled queue mode"
	}
	return "I've already enabled queue mode"
}

func disableQueue(b *ttfm.Bot) string {
	if b.Config.ModQueue {
		b.ToggleModQueue()

		return "/me has disabled queue mode"
	}
	return "I've already disabled queue mode"

}
