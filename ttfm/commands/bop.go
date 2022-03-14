package commands

import "github.com/andreapavoni/ttfm_bot/ttfm"

func BopCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	b.Bop()
	return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypeNone}
}
