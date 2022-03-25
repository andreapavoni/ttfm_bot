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
		msg := fmt.Sprintf("Availble configs: autodjslots, autowelcome, bot, maxduration, maxsongs, qinviteduration, songstats")
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}

	if len(cmd.Args) == 1 {
		key := cmd.Args[0]
		msg := "Current setting for "

		switch key {
		case "autodjslots":
			msg += fmt.Sprintf("`autodjslots` is: %d", b.Config.AutoDjMinDjs)
			return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
		case "maxduration":
			msg += fmt.Sprintf("`maxduration` is: %d", b.Config.MaxSongDuration)
			return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
		case "maxsongs":
			msg += fmt.Sprintf("`maxsongs` is: %d", b.Config.MaxSongsPerDj)
			return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
		case "songstats":
			msg += fmt.Sprintf("`songstats` is: %v", printBool(b.Config.AutoShowSongStatsEnabled))
			return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
		case "autowelcome":
			msg += fmt.Sprintf("`autowelcome` is: %v", printBool(b.Config.AutoWelcomeEnabled))
			return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
		case "qinviteduration":
			msg += fmt.Sprintf("`qinviteduration` is: %v", b.Config.QueueInviteDuration)
			return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
		case "bot":
			msg += fmt.Sprintf("`bot` is: %v", printBool(b.Config.SetBot))
			return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
		default:
			return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("I can't find the setting you specified")}
		}
	}

	if len(cmd.Args) != 2 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the config key and its value")}
	}

	key := cmd.Args[0]

	var value interface{}
	var err error

	switch key {
	case "autodjslots":
		value, err = setInt(&b.Config.AutoDjMinDjs, cmd.Args[1])
	case "maxduration":
		value, err = setInt(&b.Config.MaxSongDuration, cmd.Args[1])
		b.Actions.EnforceSongDuration()
	case "maxsongs":
		value, err = setInt(&b.Config.MaxSongsPerDj, cmd.Args[1])
	case "songstats":
		value, err = setBool(&b.Config.AutoShowSongStatsEnabled, cmd.Args[1])
	case "autowelcome":
		value, err = setBool(&b.Config.AutoWelcomeEnabled, cmd.Args[1])
	case "qinviteduration":
		value, err = setInt(&b.Config.QueueInviteDuration, cmd.Args[1])
	case "bot":
		value, err = setBool(&b.Config.SetBot, cmd.Args[1])
		if value.(bool) {
			b.Actions.SetBot()
		}

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
		return true, nil
	case "off":
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
