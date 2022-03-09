package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func FanCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if len(args) < 1 {
		return "", user, errors.New("You must specify the username of the user you want to become a fan")
	}

	fannedUserName := strings.Join(args, " ")
	fannedUser, err := b.UserFromName(fannedUserName)

	if err != nil {
		return "", user, errors.New("I can't find the user you want to fan: @" + fannedUserName)
	}

	if err := b.Fan(fannedUser.Id); err != nil {
		return "", user, errors.New("I failed to fan @" + fannedUserName)
	}

	msg := fmt.Sprintf("I became a fan of @%s", fannedUserName)
	return msg, user, nil
}
