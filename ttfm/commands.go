package ttfm

import (
	"errors"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Command struct {
	Help               string
	Handler            CommandHandler
	AuthorizationRoles []UserRole
}

func (c *Command) Run(b *Bot, u *User, i *CommandInput) *CommandOutput {
	if err := CheckAuthorizations(b, u, c.AuthorizationRoles...); err != nil {
		return &CommandOutput{User: u, ReplyType: MessageTypePm, Err: err}
	}

	return c.Handler(b, i)
}

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
	command, err := b.recognizeCommand(cmd)
	logFields := logrus.Fields{
		"text":     i.Text,
		"cmd":      cmd,
		"args":     args,
		"userId":   user.Id,
		"userName": user.Name,
	}

	if err != nil {
		logrus.WithFields(logFields).Info(logTag + ":CMD:ERR")

		switch i.Source {
		case MessageTypePm:
			b.PrivateMessage(user.Id, err.Error())
			return
		case MessageTypeRoom:
			b.RoomMessage("@" + user.Name + " " + err.Error())
			return
		default:
			return
		}
	}

	logrus.WithFields(logFields).Info(logTag + ":CMD")

	out := command.Run(b, user, &CommandInput{UserId: user.Id, Args: args, Source: i.Source})
	msg := commandOutputMessage(out)

	if msg != "" {
		switch out.ReplyType {
		case MessageTypePm:
			b.PrivateMessage(user.Id, msg)
			return
		case MessageTypeRoom:
			b.RoomMessage(msg)
			return
		default:
			return
		}
	}
}

func (b *Bot) recognizeCommand(cmd string) (*Command, error) {
	if command, ok := b.commands.Get(cmd); ok {
		return command, nil
	}
	return nil, errors.New("command not recognized. Type !help to know available commands")
}

func parseCommand(msg string) (string, []string, bool) {
	re := regexp.MustCompile(`^!(?P<cmd>[a-zA-Z+\-!?]+)(?P<args>\s?(.*)?)`)
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

func commandOutputMessage(out *CommandOutput) string {
	if out.ReplyType == MessageTypeNone {
		return ""
	}

	if out.Msg != "" && out.Err == nil {
		return out.Msg
	}

	if out.Err != nil {
		return out.Err.Error()
	}

	return ""
}

func commandLogTag(src MessageType) string {
	switch src {
	case MessageTypePm:
		return "MSG:PM"
	default:
		return "MSG:ROOM"
	}
}
