package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func SetConfigCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Set config values at runtime",
		Handler:            setConfigCommandHandler,
	}
}

func setConfigCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if len(cmd.Args) == 0 {
		msg := fmt.Sprintf("Availble configs: autobop, autodj, autodjslots, autosnag, autowelcome, bot, djstats, maxduration, maxsongs, qinviteduration, queue, songstats")
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}

	if len(cmd.Args) == 1 {
		key := cmd.Args[0]
		msg := "Current setting for "

		switch key {
		case "autobop":
			msg += fmt.Sprintf("`autobop` is: %s", printBool(b.Config.AutoBopEnabled))
		case "autodj":
			msg += fmt.Sprintf("`autodj` is: %s", printBool(b.Config.AutoDjEnabled))
		case "autodjslots":
			msg += fmt.Sprintf("`autodjslots` is: %d", b.Config.AutoDjMinDjs)
		case "autosnag":
			msg += fmt.Sprintf("`autosnag` is: %s", printBool(b.Config.AutoSnagEnabled))
		case "autowelcome":
			msg += fmt.Sprintf("`autowelcome` is: %s", printBool(b.Config.AutoWelcomeEnabled))
		case "bot":
			msg += fmt.Sprintf("`bot` is: %s", printBool(b.Config.SetBot))
		case "djstats":
			msg += fmt.Sprintf("`djstats` is: %s", printBool(b.Config.AutoShowDjStatsEnabled))
		case "maxduration":
			msg += fmt.Sprintf("`maxduration` is: %d", b.Config.MaxSongDuration)
		case "maxsongs":
			msg += fmt.Sprintf("`maxsongs` is: %d", b.Config.MaxSongsPerDj)
		case "qinviteduration":
			msg += fmt.Sprintf("`qinviteduration` is: %d", b.Config.QueueInviteDuration)
		case "queue":
			msg += fmt.Sprintf("`queue` is: %s", printBool(b.Config.QueueEnabled))
		case "songstats":
			msg += fmt.Sprintf("`songstats` is: %s", printBool(b.Config.AutoShowSongStatsEnabled))
		default:
			return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("I can't find the setting you specified")}
		}
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}

	if len(cmd.Args) != 2 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the config key and its value")}
	}

	key := cmd.Args[0]
	var value interface{}
	var err error

	switch key {
	case "autobop":
		value, err = setBool(&b.Config.AutoBopEnabled, cmd.Args[1])
	case "autodj":
		value, err = setBool(&b.Config.AutoDjEnabled, cmd.Args[1])
	case "autodjslots":
		value, err = setInt(&b.Config.AutoDjMinDjs, cmd.Args[1])
	case "autosnag":
		value, err = setBool(&b.Config.AutoSnagEnabled, cmd.Args[1])
	case "autowelcome":
		value, err = setBool(&b.Config.AutoWelcomeEnabled, cmd.Args[1])
	case "bot":
		value, err = setBool(&b.Config.SetBot, cmd.Args[1])
		if value.(bool) {
			b.Actions.SetBot()
		}
	case "djstats":
		value, err = setBool(&b.Config.AutoShowDjStatsEnabled, cmd.Args[1])
		if value.(bool) {
			b.Actions.SetBot()
		}
	case "maxduration":
		value, err = setInt(&b.Config.MaxSongDuration, cmd.Args[1])
		b.Actions.EnforceSongDuration()
	case "maxsongs":
		value, err = setInt(&b.Config.MaxSongsPerDj, cmd.Args[1])
	case "qinviteduration":
		value, err = setInt(&b.Config.QueueInviteDuration, cmd.Args[1])
	case "queue":
		value, err = setBool(&b.Config.QueueEnabled, cmd.Args[1])
	case "songstats":
		value, err = setBool(&b.Config.AutoShowSongStatsEnabled, cmd.Args[1])
	default:
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("I can't find the key you want to set")}
	}

	if err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: err}
	}

	b.Config.Save()
	msg := fmt.Sprintf("/me has set `%s` to: `%v`", key, cmd.Args[1])
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func setInt(cfg *int64, val string) (int64, error) {
	if parsed, err := strconv.ParseInt(val, 10, 32); err == nil {
		*cfg = parsed
		return parsed, nil
	}
	return 0, errors.New("I can't parse the numeric value")
}

func setBool(cfg *bool, val string) (bool, error) {
	switch val {
	case "on":
		*cfg = true
		return true, nil
	case "off":
		*cfg = false
		return false, nil
	default:
		return false, errors.New("I can't parse `on` or `off` value")
	}
}

func printBool(value bool) string {
	if value {
		return "on"
	} else {
		return "off"
	}
}
