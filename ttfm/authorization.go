package ttfm

import "errors"

type UserRole int

const (
	UserRoleNone UserRole = iota
	UserRoleAdmin
	UserRoleModerator
	UserRoleBotModerator
)

func (m UserRole) String() string {
	switch m {
	case UserRoleAdmin:
		return "admin"
	case UserRoleModerator:
		return "moderator"
	case UserRoleBotModerator:
		return "bot_moderator"
	default:
		return "none"
	}
}

type AuthFunc func(b *Bot, u *User) error

func CheckAuthorizations(b *Bot, u *User, roles ...UserRole) error {
	for _, r := range roles {
		f := recognizeRole(r)
		if err := f(b, u); err != nil {
			return err
		}
	}
	return nil
}

func recognizeRole(r UserRole) AuthFunc {
	switch r {
	case UserRoleAdmin:
		return requireAdmin
	case UserRoleModerator:
		return requireModerator
	case UserRoleBotModerator:
		return requireBotModerator
	default:
		return requireNoRole
	}
}

func requireNoRole(b *Bot, u *User) error { return nil }

func requireAdmin(b *Bot, u *User) error {
	if !b.UserIsAdmin(u) {
		return errors.New("I won't obey you because you aren't one of my admins or we aren't in the same room")
	}

	return nil
}

func requireModerator(b *Bot, u *User) error {
	if !b.UserIsModerator(u.Id) {
		return errors.New("Sorry, I can't proceed because I'm not a moderator in this room")
	}

	return nil
}

func requireBotModerator(b *Bot, u *User) error {
	if !b.UserIsModerator(b.Config.UserId) {
		return errors.New("Sorry, I can't proceed because I'm not a moderator in this room")
	}

	return nil
}
