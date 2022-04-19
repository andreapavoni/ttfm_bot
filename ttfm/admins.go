package ttfm

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type Admins struct {
	*collections.SmartMap[string]
	brain *Brain
}

func NewAdmins(brain *Brain) *Admins {
	a := &Admins{SmartMap: collections.NewSmartMap[string](), brain: brain}
	a.loadAdmins()

	return a
}

func (a *Admins) Put(userId, userName string) error {
	if !a.SmartMap.HasKey(userId) {
		a.SmartMap.Set(userId, userName)
		return a.Save()
	}
	return errors.New("Already present in my admins list")
}


func (a *Admins) Delete(userId string) error {
	if a.SmartMap.HasKey(userId) {
		a.SmartMap.Delete(userId)
		return a.Save()
	}
	return errors.New("Unable to find Id in admins list")
}

func (a *Admins) Get(userId string) (*User, error) {
	if name, ok := a.SmartMap.Get(userId); ok {
		return &User{Id: userId, Name: name}, nil
	}
	return &User{}, errors.New("Admin with ID " + userId + " wasn't found")
}

func (a *Admins) Save() error {
	users := []User{}
	for i := range a.Iter() {
		users = append(users, User{Id: i.Key, Name: i.Value})
	}

	return a.brain.Put("admins", &users)
}

func (a *Admins) loadAdmins() error {
	admins := []User{}
	err := a.brain.Get("admins", &admins)
	if err != nil {
		return err
	}

	for _, u := range admins {
		a.Put(u.Id, u.Name)
	}
	return nil
}
