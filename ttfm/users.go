package ttfm

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type User struct {
	Id   string
	Name string
}

type Users struct {
	bot *Bot
	*collections.SmartMap[string]
}

func NewUsers(b *Bot) *Users {
	return &Users{bot: b, SmartMap: collections.NewSmartMap[string]()}
}

func (u *Users) Update(users []User) {
	u.SmartMap = collections.NewSmartMap[string]()
	for _, usr := range users {
		u.AddUser(usr.Id, usr.Name)
	}
}

func (u *Users) AddUser(id, name string) {
	u.Set(id, name)
}

func (u *Users) RemoveUser(id string) {
	u.Delete(id)
}

// UserFromId
func (u *Users) UserFromApi(userId string) (*User, error) {
	profile, err := u.bot.api.GetProfile(u.bot.Identity.Id)

	if err != nil {
		return nil, err
	}

	return &User{Id: userId, Name: profile.Name}, nil
}

func (u *Users) UserFromId(userId string) (*User, error) {
	if userName, ok := u.Get(userId); ok {
		return &User{Id: userId, Name: userName}, nil
	}
	return nil, errors.New("User with ID " + userId + " wasn't found")
}

// UserFromName
func (u *Users) UserFromName(userName string) (*User, error) {
	if id, err := u.bot.api.GetUserID(userName); err == nil {
		return &User{Id: id, Name: userName}, nil
	} else {
		return nil, err
	}
}

// UserIsAdmin
func (u *Users) UserIsAdmin(userId string) bool {
	return u.bot.Admins.HasKey(userId)
}

// UserIsDj
func (u *Users) UserIsDj(userId string) bool {
	return u.bot.Room.Djs.HasElement(userId)
}

// UserIsCurrentDj
func (u *Users) UserIsCurrentDj(userId string) bool {
	return u.bot.Room.Song.DjId == userId
}

// UserIsModerator
func (u *Users) UserIsModerator(userId string) bool {
	return u.bot.Room.HasModerator(userId)
}

func (u *Users) CheckAuthorizations(user *User, roles ...UserRole) error {
	for _, r := range roles {
		f := r.recognize()
		if err := f(u.bot, user); err != nil {
			return err
		}
	}
	return nil
}

// BootUser from room
func (u *Users) BootUser(userId, reason string) error {
	return u.bot.api.BootUser(userId, reason)
}

// Fan another user
func (u *Users) Fan(userId string) error {
	return u.bot.api.BecomeFan(userId)
}

// Unfan
func (u *Users) Unfan(userId string) error {
	return u.bot.api.RemoveFan(userId)
}

// EscortDj
func (u *Users) EscortDj(userId string) error {
	if !u.UserIsDj(userId) {
		if userId == u.bot.Identity.Id {
			return errors.New("I'm not on stage!")
		}
		return errors.New("user is not on stage")
	}
	return u.bot.api.RemDj(userId)
}

type UserRole int

const (
	UserRoleNone UserRole = iota
	UserRoleAdmin
	UserRoleModerator
	UserRoleBotModerator
)

func (u *UserRole) String() string {
	switch *u {
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

func (u *UserRole) recognize() AuthFunc {
	switch *u {
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

type AuthFunc func(b *Bot, u *User) error

func requireNoRole(b *Bot, u *User) error { return nil }

func requireAdmin(b *Bot, u *User) error {
	if !b.Users.UserIsAdmin(u.Id) {
		return errors.New("I won't obey you because you aren't one of my admins or we aren't in the same room")
	}

	return nil
}

func requireModerator(b *Bot, u *User) error {
	if !b.Users.UserIsModerator(u.Id) {
		return errors.New("Sorry, I can't proceed because I'm not a moderator in this room")
	}

	return nil
}

func requireBotModerator(b *Bot, u *User) error {
	if !b.Users.UserIsModerator(b.Identity.Id) {
		return errors.New("Sorry, I can't proceed because I'm not a moderator in this room")
	}

	return nil
}
