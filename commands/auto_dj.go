package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func AutoDjCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {

	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if len(args) == 0 {
		return currentAutoDjStatusMsg(b.Config.AutoDj), user, nil
	}

	switch args[0] {
	case "on":
		return enableAutoDj(b), user, nil
	case "off":
		return disableAutoDj(b), user, nil
	default:
		return currentAutoDjStatusMsg(b.Config.AutoDj), user, nil
	}
}

func currentAutoDjStatusMsg(status bool) string {
	if status {
		return "Auto DJ mode is enabled"
	} else {
		return "Auto DJ mode is disabled"
	}
}

func enableAutoDj(b *ttfm.Bot) string {
	if !b.Config.AutoDj {
		b.ToggleAutoDj()
		return "I'll jump on stage when possible"
	}
	return "I've already enabled auto DJ mode"
}

func disableAutoDj(b *ttfm.Bot) string {
	if b.Config.AutoDj {
		b.ToggleAutoDj()

		return "I've disabled auto DJ mode"
	}
	return "I've already disabled auto DJ mode"

}
