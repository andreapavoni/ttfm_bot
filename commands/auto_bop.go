package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func AutoBopCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {

	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if len(args) == 0 {
		return currentAutoBopStatusMsg(b.Config.AutoBop), user, nil
	}

	switch args[0] {
	case "on":
		return enableAutoBop(b), user, nil
	case "off":
		return disableAutoBop(b), user, nil
	default:
		return currentAutoBopStatusMsg(b.Config.AutoBop), user, nil
	}
}

func currentAutoBopStatusMsg(status bool) string {
	if status {
		return "Auto bop mode is enabled"
	} else {
		return "Auto bop mode is disabled"
	}
}

func enableAutoBop(b *ttfm.Bot) string {
	if !b.Config.AutoBop {
		b.ToggleAutoBop()
		return "I'm going to bop every song played from now on"
	}
	return "I'm already doing bop for every song played"

}

func disableAutoBop(b *ttfm.Bot) string {
	if b.Config.AutoBop {
		b.ToggleAutoBop()

		return "I won't bop songs played from now on"
	}
	return "I'm already not doing bop songs played"

}
