package ttfm

import (
	"fmt"
	"time"

	"github.com/alaingilbert/ttapi"
	"github.com/andreapavoni/ttfm_bot/utils"
	"github.com/sirupsen/logrus"
)

func onReady(b *Bot) {
	logrus.Info("BOT:READY")

	if b.Config.SetBot {
		utils.ExecuteDelayedRandom(30, func() { b.api.SetBot() })
	}

	b.LoadPlaylists()
	b.SwitchPlaylist(b.Config.CurrentPlaylist)
}

func onRoomChanged(b *Bot, e ttapi.RoomInfoRes) {
	if err := b.Room.Update(e); err != nil {
		panic(err)
	}

	if b.Config.AutoBop {
		utils.ExecuteDelayedRandom(15, b.Bop)
	}

	if b.Config.AutoDj && e.Room.Metadata.Djcount <= int(b.Config.AutoDjCountTrigger) {
		b.Dj()
	}

	logrus.WithFields(logrus.Fields{
		"moderators": e.Room.Metadata.ModeratorID,
		"room":       b.Room.Name,
		"roomId":     b.Room.Id,
		"shortcut":   b.Room.Shortcut,
		"djsCount":   b.Room.Djs.Size(),
		"maxDjs":     b.Room.MaxDjs,
		"listeners":  e.Room.Metadata.Listeners,
	}).Info("BOT:ROOM:CHANGED")
}

func onNewSong(b *Bot, e ttapi.NewSongEvt) {
	// show song stats
	if b.UserIsModerator(b.Config.UserId) && b.Config.AutoShowSongStats {
		header, data := b.ShowSongStats()
		b.RoomMessage(header)
		utils.ExecuteDelayed(time.Duration(5)*time.Millisecond, func() {
			b.RoomMessage(data)
		})
	}

	// escort people off the stage
	if b.UserIsModerator(b.Config.UserId) && b.Room.escorting.Size() > 0 {
		utils.ExecuteDelayed(time.Duration(10)*time.Millisecond, func() {
			for _, u := range b.Room.escorting.List() {
				if b.Room.Djs.HasElement(u) {
					b.EscortDj(u)
					b.RemoveDjEscorting(u)
				}
			}
		})
	}

	// forward queue
	if b.UserIsModerator(b.Config.UserId) && b.Config.ModQueue {
		stageIsFull := (b.Room.MaxDjs - b.Room.Djs.Size()) == 0
		queueNotEmpty := b.Queue.Size() > 0
		if err := b.EscortDj(b.Room.Song.djId); err == nil && stageIsFull && queueNotEmpty {
			b.Queue.Push(b.Room.Song.djId)
			b.PrivateMessage(b.Room.Song.djId, "Thank you for your awesome set. You've been temporarily removed from the stage and automatically added to the queue. I'll let you know when it'll be your turn again. If you want to opt-out, just type !q- and you'll be removed.")
		}

		if newDjId, err := b.Queue.Shift(int(b.Config.ModQueueInviteDuration)); err == nil {
			newDj, _ := b.UserFromId(newDjId)

			msg := fmt.Sprintf("Hey @%s! you can now jump on stage :) Your reserved slot will last %d minute(s) from now, grab it till you can!", newDj.Name, b.Config.ModQueueInviteDuration)
			b.RoomMessage(msg)
		}
	}

	// when bot is djing, push the last song to bottom of its playlist
	if b.Room.Song.djId == b.Config.UserId {
		b.PushSongBottomPlaylist()
	}

	logrus.WithFields(logrus.Fields{
		"djName": b.Room.Song.DjName,
		"djId":   b.Room.Song.djId,
		"title":  b.Room.Song.Title,
		"artist": b.Room.Song.Artist,
		"length": b.Room.Song.Length,
		"up":     b.Room.Song.up,
		"down":   b.Room.Song.down,
		"snag":   b.Room.Song.snag,
	}).Info("ROOM:LAST_SONG_STATS")

	// update room with new data
	b.Room.UpdateModerators(e.Room.Metadata.ModeratorID)
	b.Room.UpdateDjs(e.Room.Metadata.Djs)

	// update song
	song := e.Room.Metadata.CurrentSong
	b.Room.Song.Reset(song.ID, song.Metadata.Song, song.Metadata.Artist, song.Metadata.Length, song.Djname, song.Djid)

	// enforce song duration to avoid trolls with 2hours tracks
	maxDurationSeconds := int(b.Config.ModSongsMaxDuration) * 60
	durationDiff := maxDurationSeconds - b.Room.Song.Length
	if b.UserIsModerator(b.Config.UserId) && maxDurationSeconds > 0 && durationDiff < 0 {
		utils.ExecuteDelayed(time.Duration(maxDurationSeconds)*time.Second, func() {
			b.SkipSong()
		})
		msg := fmt.Sprintf("@%s song will be skipped at %s because it exceeds the limit.", song.Djname, utils.FormatSecondsToMinutes(maxDurationSeconds))
		b.RoomMessage(msg)
	}

	logrus.WithFields(logrus.Fields{
		"djName": song.Djname,
		"djId":   song.Djid,
		"title":  song.Metadata.Song,
		"artist": song.Metadata.Artist,
		"length": song.Metadata.Length,
	}).Info("ROOM:NEW_SONG")

	// auto bop
	if b.Config.AutoBop {
		utils.ExecuteDelayedRandom(30, b.Bop)
	}

	// auto snag
	if b.Config.AutoSnag {
		utils.ExecuteDelayedRandom(30, func() { b.Snag() })
	}
}

func onUpdateVotes(b *Bot, e ttapi.UpdateVotesEvt) {
	b.Room.Song.UpdateStats(e.Room.Metadata.Upvotes, e.Room.Metadata.Downvotes, b.Room.Song.snag)
	userId, vote := b.Room.Song.UnpackVotelog(e.Room.Metadata.Votelog)
	user, _ := b.UserFromId(userId)

	logrus.WithFields(logrus.Fields{
		"up":        e.Room.Metadata.Upvotes,
		"down":      e.Room.Metadata.Downvotes,
		"listeners": e.Room.Metadata.Listeners,
		"userId":    userId,
		"vote":      vote,
		"userName":  user.Name,
	}).Info("SONG:VOTE")
}

func onSnagged(b *Bot, e ttapi.SnaggedEvt) {
	b.Room.Song.UpdateStats(b.Room.Song.up, b.Room.Song.down, b.Room.Song.snag+1)
	user, _ := b.UserFromId(e.UserID)

	logrus.WithFields(logrus.Fields{
		"userId":   e.UserID,
		"userName": user.Name,
		"roomId":   e.RoomID,
	}).Info("SONG:SNAG")
}

func onRegistered(b *Bot, e ttapi.RegisteredEvt) {
	u := e.User[0]
	if u.ID == b.Config.UserId {
		return
	}

	b.Room.AddUser(u.ID, u.Name)
	user, _ := b.UserFromId(u.ID)

	if b.Config.ModAutoWelcome && b.UserIsModerator(b.Config.UserId) {
		msg := fmt.Sprintf("Hey @%s, welcome to %s! :) Type `!help` to know how to interact with me.", user.Name, b.Room.Name)
		b.RoomMessage(msg)
	}

	logrus.WithFields(logrus.Fields{
		"userId":   u.ID,
		"userName": u.Name,
		"fans":     u.Fans,
		"points":   u.Points,
	}).Info("ROOM:USER_JOINED")
}

func onDeregistered(b *Bot, e ttapi.DeregisteredEvt) {
	u := e.User[0]
	if u.ID == b.Config.UserId {
		return
	}

	b.Room.RemoveDj(u.ID)
	b.Room.RemoveUser(u.ID)

	if b.UserIsModerator(b.Config.UserId) && b.Room.escorting.Size() > 0 {
		if b.Room.Djs.HasElement(u.Userid) {
			b.RemoveDjEscorting(u.Userid)
		}
	}

	if b.UserIsModerator(b.Config.UserId) && b.Config.ModQueue {
		b.Queue.Remove(u.ID)
	}

	if b.Config.AutoDj && b.Room.Djs.Size() <= int(b.Config.AutoDjCountTrigger) {
		b.Dj()
	}

	logrus.WithFields(logrus.Fields{
		"userId":   u.ID,
		"userName": u.Name,
		"fans":     u.Fans,
		"points":   u.Points,
	}).Info("ROOM:USER_LEFT")
}

func onAddDj(b *Bot, e ttapi.AddDJEvt) {
	u := e.User[0]

	logrus.WithFields(logrus.Fields{
		"userId":   u.Userid,
		"userName": u.Name,
		"fans":     u.Fans,
		"points":   u.Points,
	}).Info("STAGE:DJ_JOINED")

	b.Room.AddDj(u.Userid)

	if b.UserIsModerator(b.Config.UserId) && b.Config.ModQueue && !b.Queue.CheckReservation(u.ID) {
		b.EscortDj(u.ID)
		msg := fmt.Sprintf("Sorry to remove you, @%s. But queue is active and the available slot is reserved. You can join the queue by typing !q+", u.Name)
		b.RoomMessage(msg)
		return
	}

	stageIsFull := (b.Room.MaxDjs - b.Room.Djs.Size()) == 0
	if b.UserIsModerator(b.Config.UserId) && !b.Config.ModQueue && stageIsFull {
		b.ModQueue(true)
		b.RoomMessage("/me has enabled queue mode")
	}

	if stageIsFull && b.UserIsDj(b.Config.UserId) {
		if b.UserIsCurrentDj(b.Config.UserId) {
			b.AddDjEscorting(b.Config.UserId)
			b.RoomMessage("/me will leave the stage at the end of this song to free a slot for humans")
		} else {
			b.EscortDj(b.Config.UserId)
			b.RoomMessage("/me leaves the stage to free a slot for humans")
		}
	}
}

func onRemDj(b *Bot, e ttapi.RemDJEvt) {
	u := e.User[0]
	b.Room.RemoveDj(u.Userid)

	if b.UserIsModerator(b.Config.UserId) && b.Room.escorting.Size() > 0 {
		if b.Room.Djs.HasElement(u.Userid) {
			b.RemoveDjEscorting(u.Userid)
		}
	}

	if b.Config.AutoDj && u.Userid == b.Config.UserId && e.Modid != "" {
		b.AutoDj(false)
		return
	}

	if b.Config.AutoDj && b.Room.Djs.Size() <= int(b.Config.AutoDjCountTrigger) {
		b.Dj()
	}

	logrus.WithFields(logrus.Fields{
		"userId":    u.Userid,
		"userName":  u.Name,
		"moderator": e.Modid,
	}).Info("STAGE:DJ_LEFT")

	stageIsAvailable := (b.Room.MaxDjs - b.Room.Djs.Size()) > 0
	if b.UserIsModerator(b.Config.UserId) && b.Config.ModQueue && stageIsAvailable {
		if newDjId, err := b.Queue.Shift(int(b.Config.ModQueueInviteDuration)); err == nil {
			var msg string
			newDj, _ := b.UserFromId(newDjId)

			if b.Queue.Size() > 0 {
				msg = fmt.Sprintf("Hey @%s! you can now jump on stage :) Your reserved slot will last %d minute(s) from now, grab it till you can!", newDj.Name, b.Config.ModQueueInviteDuration)
			} else {
				msg = fmt.Sprintf("Hey @%s! you can now jump on stage :)", newDj.Name)
				b.ModQueue(false)
				b.RoomMessage("/me has disabled queue mode")
			}

			b.RoomMessage(msg)
		}
	}
}
