package ttfm

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/alaingilbert/ttapi"
	"github.com/sirupsen/logrus"
)

type CommandHandlerSpeak func(*Bot, ttapi.SpeakEvt, ...string)
type CommandHandlerPm func(*Bot, ttapi.PmmedEvt, ...string)
type CommandHandler func(*Bot, string, ...string) (string, *User, error)

var commands = map[string]CommandHandler{
	"!escortme": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}

		if requireBotModerator(b, userId) {
			msg := "I'm going to escort you after your next song has been played"
			return msg, user, nil
		}
		return "", user, errors.New("I can't proceed because I'm not moderator in this room")
	},
	"!autodj": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}

		if !requireAdmin(b, userId) {
			return "", user, errors.New("I won't obey you because you aren't one of my admins")
		}

		b.AutoDj()
		return "Jumping on stage!", user, nil
	},
	"!autodj+": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}

		if !requireAdmin(b, userId) {
			return "", user, errors.New("I won't obey you because you aren't one of my admins")
		}

		if !b.config.AutoDj {
			b.ToggleAutoDj()
			msg := "I'll jump on stage when possible"
			return msg, user, nil
		} else {
			return "I've already enabled auto DJ mode", user, nil
		}
	},
	"!autodj-": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}
		if !requireAdmin(b, userId) {
			return "", user, errors.New("I won't obey you because you aren't one of my admins")
		}

		if b.config.AutoDj {
			b.ToggleAutoDj()
			msg := "I've disabled auto DJ mode"
			return msg, user, nil
		} else {
			return "I've already disabled auto DJ mode", user, nil
		}
	},
	"!autobop+": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}
		if !requireAdmin(b, userId) {
			return "", user, errors.New("I won't obey you because you aren't one of my admins")
		}

		if !b.config.AutoBop {
			b.ToggleAutoBop()
			msg := "I'm going to bop every song played from now on"
			return msg, user, nil
		} else {
			return "I'm already doing bop for every song played", user, nil
		}

	},
	"!autobop-": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}
		if !requireAdmin(b, userId) {
			return "", user, errors.New("I won't obey you because you aren't one of my admins")
		}

		if b.config.AutoBop {
			b.ToggleAutoBop()
			msg := "I won't bop songs played from now on"
			return msg, user, nil
		} else {
			return "I'm already not doing bop songs played", user, nil
		}
	},
	"!bop": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}
		if !requireAdmin(b, userId) {
			return "", user, errors.New("I won't obey you because you aren't one of my admins")
		}

		b.Bop()
		return "", user, nil
	},
	"!autosnag+": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}
		if !requireAdmin(b, userId) {
			return "", user, errors.New("I won't obey you because you aren't one of my admins")
		}

		if !b.config.AutoSnag {
			b.ToggleAutoSnag()
			msg := "I'm going to snag songs from now on"
			return msg, user, nil
		} else {
			return "I already snag songs", user, nil
		}
	},
	"!autosnag-": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}
		if !requireAdmin(b, userId) {
			return "", user, errors.New("I won't obey you because you aren't one of my admins")
		}

		if b.config.AutoSnag {
			b.ToggleAutoSnag()
			msg := "I won't snag songs anymore"
			return msg, user, nil
		} else {
			return "I already don't snag songs", user, nil
		}
	},
	"!snag": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}
		if !requireAdmin(b, userId) {
			return "", user, errors.New("I won't obey you because you aren't one of my admins")
		}

		if err := b.Snag(b.room.song.id); err == nil {
			return "I did snag this song!", user, nil
		}

		return "", user, errors.New("I've failed to snag this song")
	},
	"!fan": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}
		if !requireAdmin(b, userId) {
			return "", user, errors.New("I won't obey you because you aren't one of my admins")
		}

		if len(args) < 1 {
			return "", user, errors.New("You must specify the username ofr the user you want to become a fan")
		}

		fannedUserName := strings.Join(args, " ")
		fannedUserId, err := b.api.GetUserID(fannedUserName)

		if err != nil {
			return "", user, errors.New("I can't find the user you want to fan: @" + fannedUserName)
		}

		if err := b.api.BecomeFan(fannedUserId); err == nil {
			return "", user, errors.New("I failed to fan @" + fannedUserName)
		}

		msg := fmt.Sprintf("I became a fan of @%s", fannedUserName)
		return msg, user, nil
	},
	"!unfan": func(b *Bot, userId string, args ...string) (string, *User, error) {
		userName := b.room.UserNameFromId(userId)
		user := &User{Id: userId, Name: userName}
		if !requireAdmin(b, userId) {
			return "", user, errors.New("I won't obey you because you aren't one of my admins")
		}

		if len(args) < 1 {
			return "", user, errors.New("You must specify the username ofr the user you want to become a fan")
		}

		fannedUserName := strings.Join(args, " ")
		fannedUserId, err := b.api.GetUserID(fannedUserName)

		if err != nil {
			return "", user, errors.New("I can't find the user you want to unfan")
		}

		if err := b.api.RemoveFan(fannedUserId); err != nil {
			return "", user, errors.New("I failed to unfan @" + fannedUserName)
		}

		msg := fmt.Sprintf("I'm not a fan of @%s anymore", fannedUserName)
		return msg, user, nil
	},
}

func parseCommand(msg string) (string, []string, error) {
	re := regexp.MustCompile(`(?P<cmd>![a-zA-Z+\-!?]+)(?P<args>\s?(.*)?)`)
	matches := re.FindStringSubmatch(msg)

	if cmdIndex := re.SubexpIndex("cmd"); cmdIndex >= 0 && len(matches) > cmdIndex {
		cmd := strings.Trim(matches[cmdIndex], " ")

		if argsIndex := re.SubexpIndex("args"); argsIndex >= 0 {
			argsRaw := strings.Trim(matches[argsIndex], " ")
			args := strings.Split(argsRaw, " ")
			return cmd, args, nil
		} else {
			return cmd, []string{}, nil
		}
	} else {
		return "", []string{}, errors.New("Not a command")
	}
}

func handleCommandSpeak(b Bot, e ttapi.SpeakEvt) {
	if cmd, args, err := parseCommand(e.Text); err == nil {
		handler, err := b.recognizeCommand(cmd)

		if err != nil {
			userName := b.room.UserNameFromId(e.UserID)
			b.RoomMessage("@" + userName + " " + err.Error())
		}

		logrus.WithFields(logrus.Fields{
			"text":     e.Text,
			"cmd":      cmd,
			"args":     args,
			"userId":   e.UserID,
			"userName": e.Name,
		}).Info("MSG:ROOM:CMD")

		msg, user, err := handler(&b, e.UserID, args...)

		if msg != "" && err == nil {
			b.RoomMessage("@" + user.Name + " " + msg)
		}

		if err != nil {
			b.RoomMessage("@" + user.Name + " " + err.Error())
		}

		return
	}

	logrus.WithFields(logrus.Fields{"text": e.Text, "userId": e.UserID, "userName": e.Name}).Info("MSG:ROOM")
}

func handleCommandPm(b Bot, e ttapi.PmmedEvt) {
	userName := b.room.UserNameFromId(e.SenderID)

	if cmd, args, err := parseCommand(e.Text); err == nil {
		handler, err := b.recognizeCommand(cmd)

		if err != nil {
			b.PrivateMessage(userName, err.Error())
		}

		logrus.WithFields(logrus.Fields{
			"text":     e.Text,
			"cmd":      cmd,
			"args":     args,
			"userId":   e.SenderID,
			"userName": userName,
		}).Info("MSG:PM:CMD")

		msg, user, err := handler(&b, e.SenderID, args...)

		if msg != "" && err == nil {
			b.PrivateMessage(user.Id, msg)
		}

		if err != nil {
			b.PrivateMessage(user.Id, err.Error())
		}

		return
	}

	logrus.WithFields(logrus.Fields{"text": e.Text, "userId": e.SenderID, "userName": userName}).Info("MSG:ROOM")
}

func (b *Bot) recognizeCommand(cmd string) (CommandHandler, error) {
	if command, exists := commands[cmd]; exists {
		return command, nil
	}

	return nil, errors.New("Command not recognized")
}

func requireAdmin(b *Bot, userId string) bool {
	if !b.UserIsAdmin(userId) {
		b.PrivateMessage(userId, "Sorry, you aren't in my admins list")
		return false
	}

	return true
}

func requireBotModerator(b *Bot, userId string) bool {
	if !b.room.UserIsModerator(b.config.UserId) {
		b.PrivateMessage(userId, "Sorry, I can't proceed because I'm not a moderator in this room")
		return false
	}

	return true
}
