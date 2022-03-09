package commands

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func DjCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if b.UserIsDj(b.Config.UserId) {
		if b.UserIsCurrentDj(b.Config.UserId) {
			if err := b.AddDjEscorting(b.Config.UserId); err != nil {
				return "", user, errors.New("I've had some problem in getting into the escorting list: " + err.Error())
			}
			return "Ok, I'll get off the stage at the end of the current song!", user, nil
		}

		b.EscortDj(b.Config.UserId)
		return "", user, nil
	}

	b.AutoDj()
	return "Ok, I'm going to spin some tracks on stage!", user, nil
}
