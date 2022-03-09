package commands

import "github.com/andreapavoni/ttfm_bot/ttfm"

func BopCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	b.Bop()
	return "", user, nil
}
