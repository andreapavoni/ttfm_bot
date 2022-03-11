package commands

import (
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func QueueRemoveCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if err := requireBotModerator(b, user); err != nil {
		return "", user, err
	}

	if err := b.Queue.Remove(user.Id); err != nil {
		return "", user, err
	}

	msg := fmt.Sprintf("/me removed %s from the queue", user.Name)
	return msg, user, nil
}
