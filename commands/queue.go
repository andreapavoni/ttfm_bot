package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func QueueCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if err := requireBotModerator(b, user); err != nil {
		return "", user, err
	}

	switch args[0] {
	case "on":
		return enableQueue(b), user, nil
	case "off":
		return disableQueue(b), user, nil
	default:
		return currentQueueStatusMsg(b.Config.ModQueue), user, nil
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
