package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func UnfanCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if len(args) < 1 {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("You must specify the username of the user you want to unfan")}
	}

	unfannedUser, err := b.UserFromName(strings.Join(args, " "))

	if err != nil {

		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("I can't find the user you want to unfan")}
	}

	if err := b.Unfan(unfannedUser.Id); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("I was unable to unfan @" + unfannedUser.Name)}
	}

	msg := fmt.Sprintf("I'm not a fan of @%s anymore", unfannedUser.Name)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyWith: "pm", Err: nil}
}
