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

func NewActions(b *Bot) *Actions {
	return &Actions{bot: b}
}

func (a *Actions) AutoBop() {
	if a.bot.Config.AutoBopEnabled {
		utils.ExecuteDelayedRandom(30, a.bot.Room.Song.Bop)
	}
}

func (a *Actions) Bop() {
	a.bot.Room.Song.Bop()
}

func (a *Actions) AutoSnag() {
	if a.bot.Config.AutoSnagEnabled {
		utils.ExecuteDelayedRandom(30, func() { a.bot.CurrentPlaylist.Snag() })
	}
}

func (a *Actions) AutoDj() {
	if !a.bot.Users.UserIsDj(a.bot.Identity.Id) {
		if err := a.bot.api.AddDj(); err == nil {
			a.bot.RoomMessage("/me is going on stage")
		}
	}
}

func (a *Actions) ConsiderStartAutoDj() {
	if a.bot.Config.AutoDjEnabled && a.bot.Room.Djs.Size() <= int(a.bot.Config.AutoDjMinDjs) {
		a.AutoDj()
	}
}

func (a *Actions) ConsiderStopAutoDj() {
	if a.bot.Users.UserIsDj(a.bot.Identity.Id) && a.bot.Room.Djs.Size() > int(a.bot.Config.AutoDjMinDjs) {
		if a.bot.Users.UserIsCurrentDj(a.bot.Identity.Id) {
			a.bot.Room.AddDjEscorting(a.bot.Identity.Id)
			a.bot.RoomMessage("/me will leave the stage at the end of this song to free a slot for humans")
		} else {
			a.bot.Room.EscortDj(a.bot.Identity.Id)
			a.bot.RoomMessage("/me leaves the stage to free a slot for humans")
		}
	}
}

func (a *Actions) ForwardQueue() {
	if !(a.bot.Users.UserIsModerator(a.bot.Identity.Id) && a.bot.Config.QueueEnabled && a.bot.Queue.Size() > 0) {
		return
	}

	availableSlots := a.bot.Room.MaxDjs - a.bot.Room.Djs.Size()

	// stage is full, escort last dj if it's not djing (ex. true on newSong, but false on remDj or deregisterUser)
	if availableSlots == 0 && !a.bot.Users.UserIsCurrentDj(a.bot.Room.Song.DjId) {
		if err := a.bot.Room.EscortDj(a.bot.Room.Song.DjId); err == nil {
			// put the escorted dj into queue
			a.bot.Queue.Push(a.bot.Room.Song.DjId)
			a.bot.PrivateMessage(a.bot.Room.Song.DjId, "Thank you for your awesome set. You've been temporarily removed from the stage and automatically added to the queue. I'll let you know when it'll be your turn again. If you want to opt-out, just type `!q rm` and you'll be removed.")
		}
	}

	// stage is available and queue is > 0, so shift queue and invite next in line
	if availableSlots > 0 {
		if newDjId, err := a.bot.Queue.Shift(int(a.bot.Config.QueueInviteDuration)); err == nil {
			var msg string
			newDj, _ := a.bot.Users.UserFromId(newDjId)
			msg = fmt.Sprintf("Hey @%s! you can now jump on stage! Your reserved slot will last %d minute(s) from now, grab it till you can!", newDj.Name, a.bot.Config.QueueInviteDuration)
			a.bot.RoomMessage(msg)
		}
	}
}

func (a *Actions) EnforceQueueStageReservation(userId string) {
	user, _ := a.bot.Users.UserFromId(userId)
	if a.bot.Users.UserIsModerator(a.bot.Identity.Id) && a.bot.Config.QueueEnabled && !a.bot.Queue.CheckReservation(userId) {
		a.bot.Room.EscortDj(userId)
		msg := fmt.Sprintf("Sorry to remove you, @%s. But queue is active and the available slot is reserved. You can join the queue by typing !q+", user.Name)
		a.bot.RoomMessage(msg)
		return
	}
}

func (a *Actions) ConsiderQueueStart() {
	stageIsFull := a.bot.Room.MaxDjs == a.bot.Room.Djs.Size()

	if a.bot.Users.UserIsModerator(a.bot.Identity.Id) && !a.bot.Config.QueueEnabled && stageIsFull {
		a.bot.Config.EnableQueue(true)
		a.bot.Queue.Empty()
		a.bot.RoomMessage("/me has enabled queue mode")
	}
}

func (a *Actions) ConsiderQueueStop() {
	stageIsAvailable := a.bot.Room.MaxDjs-a.bot.Room.Djs.Size() > 0

	if a.bot.Users.UserIsModerator(a.bot.Identity.Id) && a.bot.Config.QueueEnabled && stageIsAvailable && a.bot.Queue.Size() == 0 {
		a.bot.Config.EnableQueue(false)
		a.bot.Queue.Empty()
		a.bot.RoomMessage("/me has disabled queue mode")
	}
}

func (a *Actions) EscortDjs() {
	if a.bot.Users.UserIsModerator(a.bot.Identity.Id) && a.bot.Room.escorting.Size() > 0 {
		utils.ExecuteDelayed(time.Duration(10)*time.Millisecond, func() {
			for _, u := range a.bot.Room.escorting.List() {
				if a.bot.Room.Djs.HasElement(u) {
					a.bot.Room.EscortDj(u)
					a.bot.Room.RemoveDjEscorting(u)
				}
			}
		})
	}
}

func (a *Actions) LoadBotIdentity() error {
	user, err := a.bot.Users.UserFromApi(a.bot.Identity.Id)
	if err != nil {
		return err
	}

	a.bot.Identity.Name = user.Name
	return nil
}

func (a *Actions) LoadMainAdmin() {
	if len(a.bot.Admins.Keys()) == 0 {
		admin, err := a.bot.Users.UserFromApi(a.bot.Config.MainAdminId)
		if err != nil {
			panic("can't find main admin user on server")
		}
		a.bot.Admins.Put(admin.Id, admin.Name)
		a.bot.Admins.Save()
	}
}

func (a *Actions) InitPlaylists() {
	a.bot.Playlists.LoadPlaylists()
	a.bot.Playlists.Switch(a.bot.Config.CurrentPlaylist)
}

func (a *Actions) ShiftPlaylistSong() {
	if a.bot.Room.Song.DjId == a.bot.Identity.Id {
		a.bot.CurrentPlaylist.PushSongBottom()
	}
}

func (a *Actions) ShowSongStats() {
	if a.bot.Config.AutoShowSongStatsEnabled && a.bot.Room.Song.Title != "" {
		header, data := a.bot.Room.SongStats()
		a.bot.RoomMessage(header)
		utils.ExecuteDelayed(time.Duration(5)*time.Millisecond, func() {
			a.bot.RoomMessage(data)
		})
	}
}

func (a *Actions) ReloadFavRooms() {
	a.bot.FavRooms.LoadFavRoomsFromDb()
}

func (a *Actions) UpdateRoom(e ttapi.RoomInfoRes) {
	a.bot.Room.Update(e)
}

func (a *Actions) UpdateRoomFromApi() {
	a.bot.Room.UpdateDataFromApi()
}

func (a *Actions) EnforceSongDuration() {
	maxDurationSeconds := int(a.bot.Config.MaxSongDuration) * 60
	durationDiff := maxDurationSeconds - a.bot.Room.Song.Length

	if a.bot.Users.UserIsModerator(a.bot.Identity.Id) && maxDurationSeconds > 0 && durationDiff < 0 {
		msg := fmt.Sprintf("@%s this song exceeds the duration limit of %s minutes: remaning %s will be skipped", a.bot.Room.Song.DjName, utils.FormatSecondsToMinutes(maxDurationSeconds), utils.FormatSecondsToMinutes(-durationDiff))
		a.bot.RoomMessage(msg)
	} else {
		a.bot.Room.Song.StopSkipTimer()
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

func (a *Actions) RegisterUser(userId, userName string) {
	if userId == a.bot.Identity.Id {
		return
	}

	a.bot.Users.AddUser(userId, userName)
	if a.bot.Config.AutoWelcomeEnabled && a.bot.Users.UserIsModerator(a.bot.Identity.Id) {
		msg := fmt.Sprintf("Hey @%s, welcome to %s! Type `!help` to know how to interact with me.", userName, a.bot.Room.Name)
		a.bot.RoomMessage(msg)
	}
}

func (a *Actions) UnregisterUser(userId string) {
	a.bot.Room.RemoveDj(userId)
	a.bot.Users.RemoveUser(userId)

	if a.bot.Users.UserIsModerator(a.bot.Identity.Id) && a.bot.Room.escorting.Size() > 0 {
		if a.bot.Room.Djs.HasElement(userId) {
			a.bot.Room.RemoveDjEscorting(userId)
		}
	}

	if a.bot.Users.UserIsModerator(a.bot.Identity.Id) && a.bot.Config.QueueEnabled {
		a.bot.Queue.Remove(userId)
	}
}

func (a *Actions) AddDj(userId string) {
	a.bot.Room.AddDj(userId)
}

func (a *Actions) RemoveDj(userId, modId string) {
	a.bot.Room.RemoveDj(userId)
	if a.bot.Users.UserIsModerator(a.bot.Identity.Id) && a.bot.Room.escorting.Size() > 0 {
		if a.bot.Room.Djs.HasElement(userId) {
			a.bot.Room.RemoveDjEscorting(userId)
		}
	}

	// disable auto-dj if bot has been removed by a moderator
	if a.bot.Config.AutoDjEnabled && userId == a.bot.Identity.Id && modId != "" {
		status := false
		a.bot.Config.EnableAutoDj(status)

		if a.bot.Users.UserIsDj(a.bot.Identity.Id) {
			if a.bot.Users.UserIsCurrentDj(a.bot.Identity.Id) {
				a.bot.Room.AddDjEscorting(a.bot.Identity.Id)
				return
			}
			a.bot.api.RemDj("")
		}
		return
	}
}

func (a *Actions) BootUser(userId, reason string) error {
	return a.bot.Room.BootUser(userId, reason)
}

func (a *Actions) EscortDj(userId string) error {
	return a.bot.Room.EscortDj(userId)
}

func (a *Actions) AddDjEscorting(userId string) error {
	return a.bot.Room.AddDjEscorting(userId)
}
