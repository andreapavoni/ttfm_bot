package ttfm

import (
	"errors"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type CommandHandler func(*Bot, *CommandInput) *CommandOutput

type CommandInput struct {
	UserId string
	Args   []string
	Source MessageType
}

type CommandOutput struct {
	Msg       string
	User      *User
	ReplyType MessageType
	Err       error
}

type MessageInput struct {
	Text   string
	UserId string
	Source MessageType
}

type MessageType int

const (
	MessageTypeNone MessageType = iota
	MessageTypePm
	MessageTypeRoom
)

func (m MessageType) String() string {
	switch m {
	case MessageTypePm:
		return "pm"
	case MessageTypeRoom:
		return "room"
	default:
		return "none"
	}
}

func handleCommand(b *Bot, i *MessageInput) {
	user, _ := b.UserFromId(i.UserId)
	logTag := commandLogTag(i.Source)

	cmd, args, ok := parseCommand(i.Text)
	if !ok {
		logrus.WithFields(logrus.Fields{"text": i.Text, "userId": user.Id, "userName": user.Name}).Info(logTag)
		return
	}
	handler, err := b.recognizeCommand(cmd)
	logFields := logrus.Fields{
		"text":     i.Text,
		"cmd":      cmd,
		"args":     args,
		"userId":   user.Id,
		"userName": user.Name,
	}

	if err != nil {
		b.RoomMessage("@" + user.Name + " " + err.Error())
		logrus.WithFields(logFields).Info(logTag + ":CMD:ERR")
		return
	}

	logrus.WithFields(logFields).Info(logTag + ":CMD")

	out := handler(b, &CommandInput{UserId: user.Id, Args: args, Source: i.Source})

	if out.Msg != "" && out.Err == nil {
		b.RoomMessage(out.Msg)

		switch out.ReplyType {
		case MessageTypePm:
			b.PrivateMessage(user.Id, out.Msg)
			return
		case MessageTypeRoom:
			b.RoomMessage(out.Msg)
			return
		default:
			return
		}
	}

	if err != nil {
		b.PrivateMessage(user.Id, err.Error())
	}
}

func (b *Bot) recognizeCommand(cmd string) (CommandHandler, error) {
	if command, ok := b.commands.Get(cmd); ok {
		return command, nil
	}
	return nil, errors.New("command not recognized")
}

func parseCommand(msg string) (string, []string, bool) {
	re := regexp.MustCompile(`(?P<cmd>^![a-zA-Z+\-!?]+)(?P<args>\s?(.*)?)`)
	matches := re.FindStringSubmatch(msg)

	cmdIndex := re.SubexpIndex("cmd")
	if !(cmdIndex >= 0 && len(matches) > cmdIndex) {
		return "", nil, false
	}
	cmd := strings.Trim(matches[cmdIndex], " ")

	argsIndex := re.SubexpIndex("args")
	if argsIndex >= 0 && len(strings.Trim(matches[argsIndex], " ")) > 0 {
		argsRaw := strings.Trim(matches[argsIndex], " ")
		args := strings.Split(argsRaw, " ")
		return cmd, args, true
	} else {
		return cmd, nil, true
	}
}

func commandLogTag(src MessageType) string {
	switch src {
	case MessageTypePm:
		return "MSG:PM"
	default:
		return "MSG:ROOM"
	}
}
