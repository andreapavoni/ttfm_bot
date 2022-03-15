package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func SetConfigCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if len(cmd.Args) == 1 {
		key := cmd.Args[0]
		msg := "Current setting for "

		switch key {
		case "autodjslots":
			msg += fmt.Sprintf("`autodjslots` is: %d", b.Config.AutoDjCountTrigger)
			return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
		case "maxduration":
			msg += fmt.Sprintf("`maxduration` is: %d", b.Config.ModSongsMaxDuration)
			return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
		case "maxsongs":
			msg += fmt.Sprintf("`maxsongs` is: %d", b.Config.ModSongsMaxPerDj)
			return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
		default:
			return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("I can't find the setting you specified")}
		}
	}

	if len(cmd.Args) != 2 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the config key and its value")}
	}

	key := cmd.Args[0]
	value, err := parseInt(cmd.Args[1])

	if err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: err}
	}

	switch key {
	case "autodjslots":
		b.Config.AutoDjCountTrigger = value
	case "maxduration":
		b.Config.ModSongsMaxDuration = value
	case "maxsongs":
		b.Config.ModSongsMaxPerDj = value
	default:
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("I can't find the key you want to set")}
	}

	msg := fmt.Sprintf("/me has set `%s` to: `%v`", key, value)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func parseInt(val string) (int64, error) {
	if value, err := strconv.ParseInt(val, 10, 32); err == nil {
		return value, nil
	}

	return 0, errors.New("I can't parse the numeric value")
}
