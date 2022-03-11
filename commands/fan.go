package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func FanCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if len(args) < 1 {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("You must specify the username of the user you want to become a fan")}
	}

	fannedUser, err := b.UserFromName(strings.Join(args, " "))

	if err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("I can't find the user you want to fan: @" + fannedUser.Name)}
	}

	if err := b.Fan(fannedUser.Id); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("I was unable to fan @" + fannedUser.Name)}
	}

	msg := fmt.Sprintf("I became a fan of @%s", fannedUser.Name)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyWith: "pm", Err: nil}
}
