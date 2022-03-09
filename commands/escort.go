package commands

import (
	"errors"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func EscortCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if err := requireBotModerator(b, user); err != nil {
		return "", user, err
	}

	if len(args) < 1 {
		return "", user, errors.New("You must specify the username of the user you want to escort")
	}

	escortedUserName := strings.Join(args, " ")
	escortedUser, err := b.UserFromName(escortedUserName)

	if err != nil {
		return "", user, errors.New("I can't find the user you want to escort: @" + escortedUserName)
	}

	if err := b.EscortDj(escortedUser.Id); err != nil {
		return "", user, errors.New("I failed to escort @" + escortedUserName)
	}

	return "", user, nil
}
