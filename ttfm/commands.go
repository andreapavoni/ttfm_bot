package ttfm

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/andreapavoni/ttfm_bot/utils/collections"
	"github.com/sirupsen/logrus"
)

type Command struct {
	Help               string
	Handler            CommandHandler
	AuthorizationRoles []UserRole
}

func (c *Command) Run(b *Bot, u *User, i *CommandInput) *CommandOutput {
	if err := b.Users.CheckAuthorizations(u, c.AuthorizationRoles...); err != nil {
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

func (out *CommandOutput) String() string {
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

func (out *CommandOutput) SendReply(b *Bot, userId string) {
	if out.String() != "" {
		switch out.ReplyType {
		case MessageTypePm:
			b.PrivateMessage(userId, out.String())
			return
		case MessageTypeRoom:
			b.RoomMessage(out.String())
			return
		default:
			return
		}
	}
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

type MessageInput struct {
	Text   string
	UserId string
	Source MessageType
}

func (i *MessageInput) HandleCommand(b *Bot) {
	logTag := i.commandLogTag()

	cmd, args, ok := i.parseCommand(b.Config.CmdPrefix)
	if !ok {
		user, _ := b.Users.UserFromId(i.UserId)
		logrus.WithFields(logrus.Fields{"text": i.Text, "userId": i.UserId, "userName": user.Name}).Info(logTag)
		return
	}

	user, err := i.userFromCommand(b)
	if err != nil {
		logFields := logrus.Fields{
			"text":   i.Text,
			"userId": i.UserId,
		}
		logrus.WithFields(logFields).Info(logTag + ":CMD:ERR " + err.Error())
		return
	}

	command, err := b.Commands.RecognizeCommand(cmd)
	logFields := logrus.Fields{
		"text":     i.Text,
		"cmd":      cmd,
		"args":     args,
		"userId":   user.Id,
		"userName": user.Name,
	}

	if err != nil {
		logrus.WithFields(logFields).Info(logTag + ":CMD:ERR " + err.Error())

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
	out.SendReply(b, user.Id)
}

func (i *MessageInput) parseCommand(cmdPrefix string) (string, []string, bool) {
	pattern := fmt.Sprintf(`^\%s(?P<cmd>[a-zA-Z+\-!?]+)(?P<args>\s?(.*)?)`, cmdPrefix)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(i.Text)

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

func (i *MessageInput) userFromCommand(b *Bot) (*User, error) {
	if user, err := b.Users.UserFromId(i.UserId); err != nil && b.Users.UserIsAdmin(i.UserId) {
		return b.Admins.Get(i.UserId)
	} else {
		return user, nil
	}
}

func (i *MessageInput) commandLogTag() string {
	switch i.Source {
	case MessageTypePm:
		return "MSG:PM"
	default:
		return "MSG:ROOM"
	}
}

type Commands struct {
	*collections.SmartMap[*Command]
}

func NewCommands() *Commands {
	return &Commands{SmartMap: collections.NewSmartMap[*Command]()}
}

func (c *Commands) RecognizeCommand(cmd string) (*Command, error) {
	if command, ok := c.SmartMap.Get(cmd); ok {
		return command, nil
	}
	return nil, errors.New("command not recognized. Use help command to know available commands")
}

// ListCommands available
func (c *Commands) List() []string {
	return c.Keys()
}

// Get command if exists
func (c *Commands) Get(name string) (*Command, error) {
	if cmd, ok := c.SmartMap.Get(name); ok {
		return cmd, nil
	}
	return nil, errors.New("command not found")
}

// Add command with given alias
func (c *Commands) Add(alias string, cmd *Command) {
	c.SmartMap.Set(alias, cmd)
}
