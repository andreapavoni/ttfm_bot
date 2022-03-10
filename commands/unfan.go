package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func UnfanCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if len(args) < 1 {
		return "", user, errors.New("You must specify the username of the user you want to unfan")
	}

	unfannedUserName := strings.Join(args, " ")
	unfannedUser, err := b.UserFromName(unfannedUserName)

	if err != nil {
		return "", user, errors.New("I can't find the user you want to unfan")
	}

	if err := b.Unfan(unfannedUser.Id); err != nil {
		return "", user, errors.New("I was unable to unfan @" + unfannedUserName)
	}

	msg := fmt.Sprintf("I'm not a fan of @%s anymore", unfannedUserName)
	return msg, user, nil
}
