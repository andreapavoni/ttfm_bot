package commands

import (
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func QueueCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin, ttfm.UserRoleBotModerator},
		Help:               "Enable/disable queue mode. Without args prints current setting",
		Handler:            queueCommandHandler,
	}
}

func queueCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)
	if len(cmd.Args) == 0 {
		return &ttfm.CommandOutput{Msg: currentQueueStatusMsg(b, b.Config.ModQueue), User: user, ReplyType: cmd.Source}
	}

	switch cmd.Args[0] {
	case "on":
		return &ttfm.CommandOutput{Msg: enableQueue(b), User: user, ReplyType: cmd.Source}
	case "off":
		return &ttfm.CommandOutput{Msg: disableQueue(b), User: user, ReplyType: cmd.Source}
	default:
		return &ttfm.CommandOutput{Msg: currentQueueStatusMsg(b, b.Config.ModQueue), User: user, ReplyType: cmd.Source}
	}
}

func currentQueueStatusMsg(b *ttfm.Bot, status bool) string {
	if !status {
		return "Queue mode is disabled"
	}

	if b.Queue.Size() > 0 {
		return "Current queue: " + strings.Join(b.Queue.List(), ", ")
	} else {
		return "Queue mode is enabled, but empty"
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
