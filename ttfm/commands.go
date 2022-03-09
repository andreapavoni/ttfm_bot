package ttfm

import (
	"errors"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type CommandHandler func(*Bot, string, []string) (string, *User, error)

func handleCommandSpeak(b *Bot, userId string, message string) {
	user, _ := b.UserFromId(userId)

	if cmd, args, err := parseCommand(message); err == nil {
		handler, err := b.recognizeCommand(cmd)

		if err != nil {
			b.RoomMessage("@" + user.Name + " " + err.Error())
			logrus.WithFields(logrus.Fields{
				"text":     message,
				"cmd":      cmd,
				"args":     args,
				"userId":   user.Id,
				"userName": user.Name,
			}).Info("MSG:ROOM:CMD:ERR")
			return
		}

		logrus.WithFields(logrus.Fields{
			"text":     message,
			"cmd":      cmd,
			"args":     args,
			"userId":   userId,
			"userName": user.Name,
		}).Info("MSG:ROOM:CMD")

		msg, _, err := handler(b, userId, args)

		if msg != "" && err == nil {
			// b.RoomMessage("@" + user.Name + " " + msg)
			b.RoomMessage(msg)
		}

		if err != nil {
			b.RoomMessage("@" + user.Name + " " + err.Error())
		}

		return
	}

	logrus.WithFields(logrus.Fields{"text": message, "userId": userId, "userName": user.Name}).Info("MSG:ROOM")
}

func handleCommandPm(b *Bot, userId string, message string) {
	user, _ := b.UserFromId(userId)

	if cmd, args, err := parseCommand(message); err == nil {
		handler, err := b.recognizeCommand(cmd)

		if err != nil {
			b.PrivateMessage(userId, err.Error())
			logrus.WithFields(logrus.Fields{
				"text":     message,
				"cmd":      cmd,
				"args":     args,
				"userId":   userId,
				"userName": user.Name,
			}).Info("MSG:PM:CMD:ERR")
			return
		}

		logrus.WithFields(logrus.Fields{
			"text":     message,
			"cmd":      cmd,
			"args":     args,
			"userId":   userId,
			"userName": user.Name,
		}).Info("MSG:PM:CMD")

		msg, _, err := handler(b, userId, args)

		if msg != "" && err == nil {
			b.PrivateMessage(userId, msg)
		}

		if err != nil {
			b.PrivateMessage(userId, err.Error())
		}

		return
	}

	logrus.WithFields(logrus.Fields{"text": message, "userId": userId, "userName": user.Name}).Info("MSG:ROOM")
}

func (b *Bot) recognizeCommand(cmd string) (CommandHandler, error) {
	if command, ok := b.commands.Get(cmd); ok {
		return command, nil
	}
	return nil, errors.New("Command not recognized")
}

func parseCommand(msg string) (string, []string, error) {
	re := regexp.MustCompile(`(?P<cmd>^![a-zA-Z+\-!?]+)(?P<args>\s?(.*)?)`)
	matches := re.FindStringSubmatch(msg)

	if cmdIndex := re.SubexpIndex("cmd"); cmdIndex >= 0 && len(matches) > cmdIndex {
		cmd := strings.Trim(matches[cmdIndex], " ")

		argsIndex := re.SubexpIndex("args")
		if argsIndex >= 0 && len(strings.Trim(matches[argsIndex], " ")) > 0 {
			argsRaw := strings.Trim(matches[argsIndex], " ")
			args := strings.Split(argsRaw, " ")
			return cmd, args, nil
		} else {
			return cmd, nil, nil
		}
	} else {
		return "", nil, errors.New("Not a command")
	}
}
