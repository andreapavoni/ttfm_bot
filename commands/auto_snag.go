package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func AutoSnagCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {

	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if len(args) == 0 {
		return currentAutoSnagStatusMsg(b.Config.AutoSnag), user, nil
	}

	switch args[0] {
	case "on":
		return enableAutoSnag(b), user, nil
	case "off":
		return disableAutoSnag(b), user, nil
	default:
		return currentAutoSnagStatusMsg(b.Config.AutoSnag), user, nil
	}
}

func currentAutoSnagStatusMsg(status bool) string {
	if status {
		return "Auto snag mode is enabled"
	} else {
		return "Auto snag mode is disabled"
	}
}

func enableAutoSnag(b *ttfm.Bot) string {
	if !b.Config.AutoSnag {
		b.ToggleAutoSnag()

		return "I'm going to snag songs from now on"
	}
	return "I'm already snagging songs"

}

func disableAutoSnag(b *ttfm.Bot) string {
	if b.Config.AutoSnag {
		b.ToggleAutoSnag()

		return "I won't snag songs anymore"
	}
	return "I'm already not snagging songs"
}
