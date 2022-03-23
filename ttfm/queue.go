package ttfm

import (
	"errors"
	"time"

	"github.com/andreapavoni/ttfm_bot/utils"
	"github.com/andreapavoni/ttfm_bot/utils/collections"
)

type Queue struct {
	reservationDjId      string
	reservationExpiresAt time.Time
	*collections.SmartList[string]
}

func NewQueue() *Queue {
	return &Queue{"", time.Now(), collections.NewSmartList[string]()}
}

func (q *Queue) Shift(inviteExpirationMinutes int) (string, error) {
	next, err := q.SmartList.Shift()
	if err != nil {
		return "", err
	}

	now := time.Now()
	waitDuration := time.Duration(int64(inviteExpirationMinutes)) * time.Minute
	q.reservationDjId = next
	q.reservationExpiresAt = now.Add(waitDuration)

	defer utils.ExecuteDelayed(waitDuration, q.resetReservation)

	return next, nil
}

func (q *Queue) Add(djId string) error {
	if q.SmartList.HasElement(djId) {
		return errors.New("You're already in queue")
	}
	q.SmartList.Push(djId)

	return nil
}

func (q *Queue) Remove(djId string) error {
	if err := q.SmartList.Remove(djId); err != nil {
		return err
	}

	if q.reservationDjId == djId {
		q.resetReservation()
	}
	return nil
}

func (q *Queue) Reset() {
	q.resetReservation()
	q.SmartList.Empty()
}

func (q *Queue) CheckReservation(djId string) bool {
	return q.reservationDjId == djId && time.Now().Unix() <= q.reservationExpiresAt.Unix()
}

func (q *Queue) resetReservation() {
	q.reservationDjId = ""
	q.reservationExpiresAt = time.Now()
}
