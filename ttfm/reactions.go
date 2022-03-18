package ttfm

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/andreapavoni/ttfm_bot/utils"
	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type Reactions struct {
	*collections.SmartMap[[]string]
	availables *collections.SmartList[string]
	brain      *Brain
	bucket     string
}

func NewReactions(brain *Brain, bucket string) *Reactions {
	r := &Reactions{SmartMap: collections.NewSmartMap[[]string](), brain: brain, availables: collections.NewSmartList[string](), bucket: bucket}
	r.loadReactions()
	return r
}

func (r *Reactions) Put(reactionName, imgUrl string) error {
	imgs, ok := r.SmartMap.Get(reactionName)
	if !ok {
		r.SmartMap.Set(reactionName, []string{imgUrl})
		return nil
	}

	if utils.IndexOf(imgUrl, imgs) >= 0 {
		return errors.New("already present")
	}

	if !r.availables.HasElement(reactionName) {
		r.availables.Push(reactionName)
	}

	imgs = append(imgs, imgUrl)
	r.SmartMap.Set(reactionName, imgs)
	r.brain.Put(r.bucket, reactionName, imgs)
	return nil
}

func (r *Reactions) Get(reactionName string) string {
	if imgs, ok := r.SmartMap.Get(reactionName); ok {
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

func (r *Reactions) loadReactions() error {
	records, err := r.brain.Db.ReadAll(r.bucket)
	if err != nil {
		return err
	}

	for key, value := range records {
		imgs := []string{}
		for _, img := range value {
			if err := json.Unmarshal([]byte(img), &imgs); err != nil {
				return err
			}
		}
		r.Set(key, imgs)
		r.availables.Push(key)
	}
	return nil
}
