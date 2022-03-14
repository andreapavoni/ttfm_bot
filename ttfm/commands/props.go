package commands

import (
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PropsCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	msg := fmt.Sprintf("🔥 Hey @%s! @%s is enjoying the song you're playing! 🚀", b.Room.Song.DjName, user.Name)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: ttfm.MessageTypeRoom}

}
