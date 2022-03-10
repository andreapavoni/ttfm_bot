package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func BootCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if err := requireBotModerator(b, user); err != nil {
		return "", user, err
	}

	if len(args) < 1 {
		return "", user, errors.New("You must specify the username of the user you kick")
	}

	kickedUser, err := b.UserFromName(strings.Join(args, " "))
	reason := fmt.Sprintf("Requested by mod - ask @%s for details", user.Name)

	if err != nil {
		return "", user, errors.New("I can't find the user you want to kick")
	}

	if err := b.BootUser(kickedUser.Id, reason); err != nil {
		return "", user, errors.New("I wasn't able to boot @" + kickedUser.Name)
	}

	return "", user, nil
}
