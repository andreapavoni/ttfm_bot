package ttfm

import (
	"fmt"
	"time"

	"github.com/alaingilbert/ttapi"
	"github.com/andreapavoni/ttfm_bot/utils"
)

type Actions struct {
	bot *Bot
}

func (a *Actions) AutoBop() {
	if a.bot.Config.AutoBop {
		utils.ExecuteDelayedRandom(30, a.bot.Bop)
	}
}

func (a *Actions) AutoSnag() {
	if a.bot.Config.AutoSnag {
		utils.ExecuteDelayedRandom(30, func() { a.bot.Snag() })
	}
}

func (a *Actions) AutoDj() {
	if a.bot.Config.AutoDj && a.bot.Room.Djs.Size() <= int(a.bot.Config.AutoDjCountTrigger) {
		a.bot.Dj()
	}
}

func (a *Actions) ForwardQueue() {
	if !(a.bot.UserIsModerator(a.bot.Config.UserId) && a.bot.Config.ModQueue && a.bot.Queue.Size() > 0) {
		return
	}

	// stage is full, escort last dj if it's not djing (ex. true on newSong, but false on remDj or deregisterUser)
	if (a.bot.Room.MaxDjs-a.bot.Room.Djs.Size()) == 0 && !a.bot.UserIsCurrentDj(a.bot.Room.Song.DjId) {
		if err := a.bot.EscortDj(a.bot.Room.Song.DjId); err == nil {
			a.bot.Queue.Push(a.bot.Room.Song.DjId)
			a.bot.PrivateMessage(a.bot.Room.Song.DjId, "Thank you for your awesome set. You've been temporarily removed from the stage and automatically added to the queue. I'll let you know when it'll be your turn again. If you want to opt-out, just type !qrm and you'll be removed.")
		}
	}

	// stage is available and queue is > 0, so shift queue and invite next in line
	if (a.bot.Room.MaxDjs - a.bot.Room.Djs.Size()) > 0 {
		if newDjId, err := a.bot.Queue.Shift(int(a.bot.Config.ModQueueInviteDuration)); err == nil {
			var msg string
			newDj, _ := a.bot.UserFromId(newDjId)

			if a.bot.Queue.Size() > 0 {
				msg = fmt.Sprintf("Hey @%s! you can now jump on stage :) Your reserved slot will last %d minute(s) from now, grab it till you can!", newDj.Name, a.bot.Config.ModQueueInviteDuration)
			} else {
				msg = fmt.Sprintf("Hey @%s! you can now jump on stage :)", newDj.Name)
				a.bot.ModQueue(false)
				a.bot.RoomMessage("/me has disabled queue mode")
			}

			a.bot.RoomMessage(msg)
		}
	}
}

func (a *Actions) EnforceQueueStageReservation(userId string) {
	user, _ := a.bot.UserFromId(userId)
	if a.bot.UserIsModerator(a.bot.Config.UserId) && a.bot.Config.ModQueue && !a.bot.Queue.CheckReservation(userId) {
		a.bot.EscortDj(userId)
		msg := fmt.Sprintf("Sorry to remove you, @%s. But queue is active and the available slot is reserved. You can join the queue by typing !q+", user.Name)
		a.bot.RoomMessage(msg)
		return
	}
}

func (a *Actions) ConsiderQueueActivation() {
	stageIsFull := (a.bot.Room.MaxDjs - a.bot.Room.Djs.Size()) == 0

	if a.bot.UserIsModerator(a.bot.Config.UserId) && !a.bot.Config.ModQueue && stageIsFull {
		a.bot.ModQueue(true)
		a.bot.RoomMessage("/me has enabled queue mode")
	}

	if stageIsFull && a.bot.UserIsDj(a.bot.Config.UserId) {
		if a.bot.UserIsCurrentDj(a.bot.Config.UserId) {
			a.bot.AddDjEscorting(a.bot.Config.UserId)
			a.bot.RoomMessage("/me will leave the stage at the end of this song to free a slot for humans")
		} else {
			a.bot.EscortDj(a.bot.Config.UserId)
			a.bot.RoomMessage("/me leaves the stage to free a slot for humans")
		}
	}
}

func (a *Actions) EscortDjs() {
	if a.bot.UserIsModerator(a.bot.Config.UserId) && a.bot.Room.escorting.Size() > 0 {
		utils.ExecuteDelayed(time.Duration(10)*time.Millisecond, func() {
			for _, u := range a.bot.Room.escorting.List() {
				if a.bot.Room.Djs.HasElement(u) {
					a.bot.EscortDj(u)
					a.bot.RemoveDjEscorting(u)
				}
			}
		})
	}
}

func (a *Actions) InitPlaylists() {
	a.bot.LoadPlaylists()
	a.bot.SwitchPlaylist(a.bot.Config.CurrentPlaylist)
}

func (a *Actions) ShiftPlaylistSong() {
	if a.bot.Room.Song.DjId == a.bot.Config.UserId {
		a.bot.PushSongBottomPlaylist()
	}
}

func (a *Actions) ShowSongStats() {
	if a.bot.UserIsModerator(a.bot.Config.UserId) && a.bot.Config.AutoShowSongStats && a.bot.Room.Song.Title != "" {
		header, data := a.bot.ShowSongStats()
		a.bot.RoomMessage(header)
		utils.ExecuteDelayed(time.Duration(5)*time.Millisecond, func() {
			a.bot.RoomMessage(data)
		})
	}
}

func (a *Actions) UpdateRoom(e ttapi.RoomInfoRes) {
	a.bot.Room.Update(e)
}

func (a *Actions) UpdateRoomFromApi() {
	if roomInfo, err := a.bot.GetRoomInfo(); err == nil {
		a.UpdateRoom(roomInfo)
	} else {
		panic(err)
	}
}

func (a *Actions) EnforceSongDuration() {
	maxDurationSeconds := int(a.bot.Config.ModSongsMaxDuration) * 60
	durationDiff := maxDurationSeconds - a.bot.Room.Song.Length
	if a.bot.UserIsModerator(a.bot.Config.UserId) && maxDurationSeconds > 0 && durationDiff < 0 {
		utils.ExecuteDelayed(time.Duration(maxDurationSeconds)*time.Second, func() {
			a.bot.SkipSong()
		})
		msg := fmt.Sprintf("@%s song will be skipped at %s because it exceeds the limit.", a.bot.Room.Song.DjName, utils.FormatSecondsToMinutes(maxDurationSeconds))
		a.bot.RoomMessage(msg)
	}
}

func (a *Actions) UpdateSongStats(ups, downs, snags int) {
	a.bot.Room.Song.UpdateStats(ups, downs, snags)
}

func (a *Actions) SetBot() {
	if a.bot.Config.SetBot {
		a.bot.api.SetBot()
	}
}

func (a *Actions) RegisterUser(userId string) {
	if userId == a.bot.Config.UserId {
		return
	}

	user, _ := a.bot.UserFromId(userId)
	a.bot.Room.AddUser(user.Id, user.Name)
	if a.bot.Config.ModAutoWelcome && a.bot.UserIsModerator(a.bot.Config.UserId) {
		msg := fmt.Sprintf("Hey @%s, welcome to %s! :) Type `!help` to know how to interact with me.", user.Name, a.bot.Room.Name)
		a.bot.RoomMessage(msg)
	}
}

func (a *Actions) UnregisterUser(userId string) {
	a.bot.Room.RemoveDj(userId)
	a.bot.Room.RemoveUser(userId)

	if a.bot.UserIsModerator(a.bot.Config.UserId) && a.bot.Room.escorting.Size() > 0 {
		if a.bot.Room.Djs.HasElement(userId) {
			a.bot.RemoveDjEscorting(userId)
		}
	}

	if a.bot.UserIsModerator(a.bot.Config.UserId) && a.bot.Config.ModQueue {
		a.bot.Queue.Remove(userId)
	}
}

func (a *Actions) AddDj(userId string) {
	a.bot.Room.AddDj(userId)
}

func (a *Actions) RemoveDj(userId, modId string) {
	a.bot.Room.RemoveDj(userId)
	if a.bot.UserIsModerator(a.bot.Config.UserId) && a.bot.Room.escorting.Size() > 0 {
		if a.bot.Room.Djs.HasElement(userId) {
			a.bot.RemoveDjEscorting(userId)
		}
	}

	// disable auto-dj if bot has been removed by a moderator
	if a.bot.Config.AutoDj && userId == a.bot.Config.UserId && modId != "" {
		a.bot.AutoDj(false)
		return
	}
}
