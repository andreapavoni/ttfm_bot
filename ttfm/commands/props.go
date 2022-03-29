package commands

import (
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PropsCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{},
		Help:               "Give props to current dj",
		Handler:            propsCommandHandler,
	}
}

func propsCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	msg := fmt.Sprintf("🔥 Hey @%s! @%s is enjoying the song you're playing! 🚀", b.Room.CurrentSong.DjName, user.Name)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: ttfm.MessageTypeRoom}

}
