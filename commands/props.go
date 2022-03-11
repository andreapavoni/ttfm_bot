package commands

import (
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PropsCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	msg := fmt.Sprintf("🔥 Hey @%s! @%s is giving you props on the song you're playing! 💣", b.Room.Song.DjName, user.Name)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyWith: "pm", Err: nil}

}
