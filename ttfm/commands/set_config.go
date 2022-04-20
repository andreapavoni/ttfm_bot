package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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

	if len(cmd.Args) >= 1 {
		key := cmd.Args[0]
		var out *ttfm.CommandOutput

		switch key {
		case "autobop":
			out = handleConfigValue(user, cmd, &b.Config.AutoBopEnabled)
		case "autodj":
			out = handleConfigValue(user, cmd, &b.Config.AutoDjEnabled)
			if b.Config.AutoDjEnabled {
				b.Actions.ConsiderStartAutoDj()
			} else {
				b.Actions.ConsiderStopAutoDj()
			}
		case "autodjslots":
			out = handleConfigValue(user, cmd, &b.Config.AutoDjMinDjs)
		case "autosnag":
			out = handleConfigValue(user, cmd, &b.Config.AutoSnagEnabled)
		case "autowelcome":
			out = handleConfigValue(user, cmd, &b.Config.AutoWelcomeEnabled)
		case "bot":
			out = handleConfigValue(user, cmd, &b.Config.SetBot)
			if b.Config.SetBot {
				b.Actions.SetBot()
			}
		case "cmdprefix":
			out = handleConfigStringValue(user, cmd, &b.Config.CmdPrefix)
		case "djstats":
			out = handleConfigValue(user, cmd, &b.Config.AutoShowDjStatsEnabled)
		case "maxduration":
			out = handleConfigValue(user, cmd, &b.Config.MaxSongDuration)
			b.Actions.EnforceSongDuration()
		case "maxsongs":
			out = handleConfigValue(user, cmd, &b.Config.MaxSongsPerDj)
		case "qinviteduration":
			out = handleConfigValue(user, cmd, &b.Config.QueueInviteDuration)
		case "queue":
			out = handleConfigValue(user, cmd, &b.Config.QueueEnabled)
		case "songstats":
			out = handleConfigValue(user, cmd, &b.Config.AutoShowSongStatsEnabled)
		case "theme":
			out = handleConfigStringValue(user, cmd, &b.Config.MusicTheme)
		default:
			return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("I can't find the setting you specified")}
		}
		if err := b.Config.Save(); err != nil {
			return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: err}
		}
		return out
	}
	msg := fmt.Sprintf("Availble configs: autobop, autodj, autodjslots, autosnag, autowelcome, bot, cmdprefix, djstats, maxduration, maxsongs, qinviteduration, queue, songstats, theme")
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}

func handleConfigStringValue(user *ttfm.User, cmd *ttfm.CommandInput, configKey *string) *ttfm.CommandOutput {
	if len(cmd.Args) == 1 {
		msg := fmt.Sprintf("Current setting for `%s` is: %s", cmd.Args[0], *configKey)
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}

	if len(cmd.Args) >= 2 {
		key := cmd.Args[0]
		*configKey = strings.Join(cmd.Args[1:], " ")
		msg := fmt.Sprintf("/me has set `%s` to: %v", key, *configKey)
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}
	return &ttfm.CommandOutput{ReplyType: ttfm.MessageTypeNone}
}

func handleConfigValue[T any](user *ttfm.User, cmd *ttfm.CommandInput, configKey *T) *ttfm.CommandOutput {
	if len(cmd.Args) == 1 {
		msg := fmt.Sprintf("Current setting for `%s` is: %s", cmd.Args[0], configValueToString(*configKey))
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}

	if len(cmd.Args) == 2 {
		key := cmd.Args[0]
		newValue, err := parseConfigValueFromString(*configKey, cmd.Args[1])
		if err != nil {
			return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: err}
		}
		*configKey = newValue.(T)
		msg := fmt.Sprintf("/me has set `%s` to: %v", key, configValueToString(newValue))
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}
	return &ttfm.CommandOutput{ReplyType: ttfm.MessageTypeNone}
}

func parseConfigValueFromString(configValue interface{}, value string) (parsed interface{}, err error) {
	switch configValue.(type) {
	case int64:
		parsed, err = parseInt(value)
	case bool:
		parsed, err = parseBool(value)
	default:
		return nil, errors.New("no valid value")
	}
	return parsed, err
}

func configValueToString(value interface{}) string {
	switch value.(type) {
	case int64:
		return fmt.Sprintf("`%d`", value.(int64))
	case bool:
		return fmt.Sprintf("`%s`", printBool(value.(bool)))
	default:
		return "not recognized"
	}
}

func parseInt(val string) (int64, error) {
	if parsed, err := strconv.ParseInt(val, 10, 32); err == nil {
		return parsed, nil
	}
	return 0, errors.New("I can't parse the numeric value")
}

func parseBool(val string) (bool, error) {
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
