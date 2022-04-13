package ttfm

import (
	"errors"
	"sort"

	"github.com/andreapavoni/ttfm_bot/utils"
	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type reaction struct {
	Name string
	Imgs []string
}

type Reactions struct {
	*collections.SmartMap[[]string]
	availables *collections.SmartList[string]
	brain      *Brain
}

func NewReactions(brain *Brain) *Reactions {
	r := &Reactions{SmartMap: collections.NewSmartMap[[]string](), brain: brain, availables: collections.NewSmartList[string]()}
	r.loadReactions()
	return r
}

func (r *Reactions) Put(reactionName, imgUrl string) error {
	imgs, ok := r.SmartMap.Get(reactionName)
	if !ok {
		r.SmartMap.Set(reactionName, []string{imgUrl})
	}

	if utils.IndexOf(imgUrl, imgs) >= 0 {
		return errors.New("already present")
	}

	if !r.availables.HasElement(reactionName) {
		r.availables.Push(reactionName)
	}

	imgs = append(imgs, imgUrl)
	r.SmartMap.Set(reactionName, imgs)
	return r.Save()
}

func (r *Reactions) Get(reactionName string) string {
	if imgs, ok := r.SmartMap.Get(reactionName); ok && len(imgs) > 0 {
		i := utils.RandomInteger(0, len(imgs)-1)
		return imgs[i]
	}
	return ""
}

func (r *Reactions) Availables() []string {
	availableReactions := r.availables.List()
	reactionsList := make([]string, len(availableReactions))
	copy(reactionsList, availableReactions)
	sort.Strings(reactionsList)
	return reactionsList
}

func (r *Reactions) Save() error {
	reactions := []reaction{}
	for i := range r.SmartMap.Iter() {
		reactions = append(reactions, reaction{Name: i.Key, Imgs: i.Value})
	}
	return r.brain.Put("reactions", &reactions)
}

func (r *Reactions) loadReactions() error {
	reactions := []reaction{}
	err := r.brain.Get("reactions", &reactions)
	if err != nil {
		return err
	}

	for _, react := range reactions {
		r.Set(react.Name, react.Imgs)
		r.availables.Push(react.Name)
	}
	return nil
}
