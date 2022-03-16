package ttfm

import (
	"github.com/andreapavoni/ttfm_bot/utils"
	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type JsonReaction struct {
	Name string
	Imgs []string
}

type Reactions struct {
	filePath string
	*collections.SmartMap[[]string]
	Availables *collections.SmartList[string]
}

func NewReactions(filePath string) *Reactions {
	r := &Reactions{SmartMap: collections.NewSmartMap[[]string](), filePath: filePath, Availables: collections.NewSmartList[string]()}
	r.loadReactions()
	return r
}

func (r *Reactions) Get(reactionName string) string {
	if imgs, ok := r.SmartMap.Get(reactionName); ok {
		i := utils.RandomInteger(0, len(imgs)-1)
		return imgs[i]
	}
	return ""
}

func (r *Reactions) loadReactions() error {
	data := []JsonReaction{}
	if err := utils.ReadJson(r.filePath, &data); err != nil {
		return err
	}

	for _, i := range data {
		r.Set(i.Name, i.Imgs)
		r.Availables.Push(i.Name)
	}
	return nil
}

func (r *Reactions) dumpReactions() error {
	data := []JsonReaction{}

	for i := range r.Iter() {
		data = append(data, JsonReaction{Name: i.Key, Imgs: i.Value})
	}

	if err := utils.WriteJson(r.filePath, &data); err != nil {
		return err
	}

	return nil
}
