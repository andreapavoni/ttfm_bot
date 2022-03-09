package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func EscortMeCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireBotModerator(b, user); err != nil {
		return "", user, err
	}

	b.AddDjEscorting(userId)
	return "I'm going to escort you after your next song has been played", user, nil
}
