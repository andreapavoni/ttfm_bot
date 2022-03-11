package commands

import (
	"errors"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func EscortCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if len(args) < 1 {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("You must specify the username of the user you want to escort")}
	}

	escortedUser, err := b.UserFromName(strings.Join(args, " "))
	if err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("I can't find the user you want to escort: @" + escortedUser.Name)}
	}

	if err := b.EscortDj(escortedUser.Id); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("I was unable to escort @" + escortedUser.Name)}

	}

	return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "none", Err: nil}
}
