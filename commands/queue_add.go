package commands

import (
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func QueueAddCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if err := requireBotModerator(b, user); err != nil {
		return "", user, err
	}

	if err := b.Queue.Add(user.Id); err != nil {
		return "", user, err
	}

	msg := fmt.Sprintf("/me put %s in the queue with position #%d", user.Name, b.Queue.IndexOf(user.Id)+1)
	return msg, user, nil
}
