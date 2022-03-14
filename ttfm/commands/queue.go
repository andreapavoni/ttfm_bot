package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func QueueCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if len(cmd.Args) == 0 {
		return &ttfm.CommandOutput{Msg: currentQueueStatusMsg(b.Config.ModQueue), User: user, ReplyType: cmd.Source}
	}

	switch cmd.Args[0] {
	case "on":
		return &ttfm.CommandOutput{Msg: enableQueue(b), User: user, ReplyType: cmd.Source}
	case "off":
		return &ttfm.CommandOutput{Msg: disableQueue(b), User: user, ReplyType: cmd.Source}
	default:
		return &ttfm.CommandOutput{Msg: currentQueueStatusMsg(b.Config.ModQueue), User: user, ReplyType: cmd.Source}
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
		b.ModQueue(true)
		return "/me has enabled queue mode"
	}
	return "I've already enabled queue mode"
}

func disableQueue(b *ttfm.Bot) string {
	if b.Config.ModQueue {
		b.ModQueue(false)

		return "/me has disabled queue mode"
	}
	return "I've already disabled queue mode"

}
