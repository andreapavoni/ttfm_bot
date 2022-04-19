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
		utils.ExecuteDelayedRandom(30, a.bot.Room.CurrentSong.Bop)
	}
}

func (a *Actions) Bop() {
	a.bot.Room.CurrentSong.Bop()
}

func (a *Actions) AutoSnag() {
	if a.bot.Config.AutoSnagEnabled {
		utils.ExecuteDelayedRandom(30, func() { utils.MaybeLogError("BOT:SNAG", a.bot.CurrentPlaylist.Snag) })
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
	// adding a small delay to avoid overlapping with Stop Dj checks
	time.Sleep(2 * time.Second)
	if a.bot.Config.AutoDjEnabled && !a.bot.Users.UserIsDj(a.bot.Identity.Id) && a.bot.Room.Djs.Size() <= int(a.bot.Config.AutoDjMinDjs) {
		a.AutoDj()
	}
}

func (a *Actions) ConsiderStopAutoDj() {
	// adding a small delay to avoid overlapping with Start Dj checks
	time.Sleep(2 * time.Second)

	if a.bot.Users.UserIsDj(a.bot.Identity.Id) && (a.bot.Room.Djs.Size()-1) > int(a.bot.Config.AutoDjMinDjs) {
		if a.bot.Users.UserIsCurrentDj(a.bot.Identity.Id) {
			utils.MaybeLogError("BOT:ADD_DJ_ESCORTING", func() error { return a.bot.Room.AddDjEscorting(a.bot.Identity.Id) })
			a.bot.RoomMessage("/me will leave the stage at the end of this song to free a slot for humans")
		} else {
			utils.MaybeLogError("BOT:ESCORT_DJ", func() error { return a.bot.Room.EscortDj(a.bot.Identity.Id) })
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
	if availableSlots == 0 && !a.bot.Users.UserIsCurrentDj(a.bot.Room.CurrentSong.DjId) {
		if err := a.bot.Room.EscortDj(a.bot.Room.CurrentSong.DjId); err == nil {
			// put the escorted dj into queue
			a.bot.Queue.Push(a.bot.Room.CurrentSong.DjId)
			a.bot.PrivateMessage(a.bot.Room.CurrentSong.DjId, "Thank you for your awesome set. You've been temporarily removed from the stage and automatically added to the queue. I'll let you know when it'll be your turn again. If you want to opt-out, just type `!q rm` and you'll be removed.")
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
		utils.MaybeLogError("BOT:ESCORT_DJ", func() error { return a.bot.Room.EscortDj(userId) })
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

func (a *Actions) EscortEscortingDj(djId string) {
	if a.bot.Users.UserIsModerator(a.bot.Identity.Id) && a.bot.Room.escorting.HasElement(djId) {
		utils.ExecuteDelayed(time.Duration(10)*time.Millisecond, func() {
			utils.MaybeLogError("BOT:ESCORT_DJ", func() error { return a.bot.Room.EscortDj(djId) })
			utils.MaybeLogError("BOT:REMOVE_DJ_ESCORTING", func() error { return a.bot.Room.RemoveDjEscorting(djId) })
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
		utils.MaybeLogError("BOT:ADD_MAIN_ADMIN", func() error { return a.bot.Admins.Put(admin.Id, admin.Name) })
	}
}

func (a *Actions) InitPlaylists() {
	utils.MaybeLogError("BOT:LOAD_PLAYLISTS", a.bot.Playlists.LoadPlaylists)
	utils.MaybeLogError("BOT:SWITCH:PLAYLIST", func() error { return a.bot.Playlists.Switch(a.bot.Config.CurrentPlaylist) })
}

func (a *Actions) ShiftPlaylistSong() {
	if a.bot.Room.CurrentSong.DjId == a.bot.Identity.Id {
		utils.MaybeLogError("BOT:PLAYLIST_SONG_BOTTOM", a.bot.CurrentPlaylist.PushSongBottom)
	}
}

func (a *Actions) ShowSongStats() {
	if a.bot.Config.AutoShowSongStatsEnabled && a.bot.Room.CurrentSong.Title != "" {
		data := a.bot.Room.SongStats()
		a.bot.RoomMessage(data)
	}
}

func (a *Actions) ReloadFavRooms() error {
	if err := a.bot.FavRooms.LoadFavRoomsFromDb(); err != nil {
		return err
	}
	return nil
}

func (a *Actions) UpdateRoom(e ttapi.RoomInfoRes) error {
	if err := a.bot.Room.Update(e); err != nil {
		return err
	}
	return nil
}

func (a *Actions) UpdateRoomFromApi() error {
	if err := a.bot.Room.UpdateDataFromApi(); err != nil {
		return err
	}
	return nil
}

func (a *Actions) EnforceSongDuration() {
	maxDurationSeconds := int(a.bot.Config.MaxSongDuration) * 60
	durationDiff := maxDurationSeconds - a.bot.Room.CurrentSong.Length

	if a.bot.Users.UserIsModerator(a.bot.Identity.Id) && maxDurationSeconds > 0 && durationDiff < 0 {
		msg := fmt.Sprintf("@%s this song exceeds the duration limit of %s minutes: remaning %s will be skipped", a.bot.Room.CurrentSong.DjName, utils.FormatSecondsToMinutes(maxDurationSeconds), utils.FormatSecondsToMinutes(-durationDiff))
		a.bot.RoomMessage(msg)
	} else {
		a.bot.Room.CurrentSong.StopSkipTimer()
	}
}

func (a *Actions) UpdateSongStats(ups, downs, snags int) {
	a.bot.Room.CurrentSong.UpdateStats(ups, downs, snags)
}

func (a *Actions) UnpackVotelog(votelog [][]string) (string, string) {
	return a.bot.Room.CurrentSong.UnpackVotelog(votelog)
}

func (a *Actions) UpdateDjStatsVote(vote string) {
	switch vote {
	case "up":
		a.bot.Room.CurrentDj.UpdateStats(1, 0, 0, 0)
	case "down":
		a.bot.Room.CurrentDj.UpdateStats(0, 1, 0, 0)
	}
}

func (a *Actions) UpdateDjStatsSnag() {
	a.bot.Room.CurrentDj.UpdateStats(0, 0, 1, 0)
}

func (a *Actions) UpdateDjStatsPlays() {
	a.bot.Room.CurrentDj.UpdateStats(0, 0, 0, 1)
}

func (a *Actions) ShowDjStats(userId string) {
	if !a.bot.Config.AutoShowDjStatsEnabled {
		return
	}

	data, err := a.bot.Room.DjStats(userId)
	if err != nil {
		return
	}
	a.bot.RoomMessage(data)
}

func (a *Actions) SetBot() {
	if a.bot.Config.SetBot {
		utils.MaybeLogError("API:SET_BOT", a.bot.api.SetBot)
	}
}

func (a *Actions) RegisterUser(userId, userName string) {
	if userId == a.bot.Identity.Id {
		return
	}

	a.bot.Users.AddUser(userId, userName)
	// using a small delay to ensure the user sees the welcome message
	utils.ExecuteDelayed(time.Duration(300)*time.Millisecond, func() {
		if a.bot.Config.AutoWelcomeEnabled && a.bot.Users.UserIsModerator(a.bot.Identity.Id) {
			msg := fmt.Sprintf("Hey @%s, welcome to `%s`! Current theme is [%s]. Type `!help` to know how to interact with me 🤖", userName, a.bot.Room.Name, a.bot.Config.MusicTheme)
			a.bot.RoomMessage(msg)
		}
	})
}

func (a *Actions) UnregisterUser(userId string) {
	a.bot.Room.RemoveDj(userId)
	a.bot.Users.RemoveUser(userId)

	if a.bot.Users.UserIsModerator(a.bot.Identity.Id) && a.bot.Room.escorting.Size() > 0 {
		if a.bot.Room.Djs.HasKey(userId) {
			utils.MaybeLogError("BOT:REMOVE_DJ_ESCORTING", func() error { return a.bot.Room.RemoveDjEscorting(userId) })
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
		if a.bot.Room.Djs.HasKey(userId) {
			utils.MaybeLogError("BOT:REMOVE_DJ_ESCORTING", func() error { return a.bot.Room.RemoveDjEscorting(userId) })
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

func (a *Actions) KillSwitch() {
	a.bot.RoomMessage("I'll be back™ 🤖")
	a.bot.api.Stop()
}
