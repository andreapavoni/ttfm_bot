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

	roomInfo, err := b.GetRoomInfo()
	if err != nil {
		logrus.WithError(err).Error("BOT:READY:ERR")
	}

	b.Room.Update(roomInfo)

	if b.Config.SetBot {
		utils.ExecuteDelayedRandom(30, func() { b.api.SetBot() })
	}

	b.LoadPlaylists()
	b.SwitchPlaylist(b.Config.CurrentPlaylist)
}

func onRoomChanged(b *Bot, e ttapi.RoomInfoRes) {
	logrus.WithFields(logrus.Fields{
		"room":       e.Room.Name,
		"moderators": e.Room.Metadata.ModeratorID,
		"maxdjs":     e.Room.Metadata.MaxDjs,
		"djsCount":   e.Room.Metadata.Djcount,
		"djs":        e.Room.Metadata.Djs,
	}).Info("BOT:ROOM:CHANGED")

	b.Room.Update(e)
	utils.ExecuteDelayedRandom(15, b.Bop)

	if b.Config.AutoDj && e.Room.Metadata.Djcount == 0 {
		b.AutoDj()
	}

	logrus.WithFields(logrus.Fields{
		"room":      b.Room.name,
		"roomId":    b.Room.id,
		"shortcut":  b.Room.shortcut,
		"djs":       b.Room.djs.Size(),
		"listeners": e.Room.Metadata.Listeners,
	}).Info("BOT:ROOM:UPDATED")
}

func onNewSong(b *Bot, e ttapi.NewSongEvt) {
	// show song stats
	if b.UserIsModerator(b.Config.UserId) && b.Config.AutoShowSongStats {
		header, data := b.ShowSongStats()

		b.RoomMessage(header)
		delay := time.Duration(10) * time.Millisecond
		utils.ExecuteDelayed(delay, func() {
			b.RoomMessage(data)
		})
	}

	// escort people off the stage
	if b.UserIsModerator(b.Config.UserId) && b.escorting.HasElement(b.Room.Song.djId) {
		b.EscortDj(b.Room.Song.djId)
		b.RemoveDjEscorting(b.Room.Song.djId)
	}

	// forward queue
	if b.UserIsModerator(b.Config.UserId) && b.Config.ModQueue {
		if err := b.EscortDj(b.Room.Song.djId); err == nil {
			b.Queue.Push(b.Room.Song.djId)
			b.PrivateMessage(b.Room.Song.djId, "Thank you for your awesome set. You've been temporarily removed from the stage and automatically added to the queue. I'll let you know when it'll be your turn again. If you want to opt-out, just type !q- and you'll be removed.")
		}

		if newDjId, err := b.Queue.Shift(int(b.Config.ModQueueInviteDuration)); err == nil {
			newDj, _ := b.UserFromId(newDjId)

			msg := fmt.Sprintf("Hey %s! you can now jump on stage :) Your reserved slot will last %d minute(s) from now, grab it till you can!", newDj.Name, b.Config.ModQueueInviteDuration)
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
	song := e.Room.Metadata.CurrentSong
	b.Room.Song.Reset(song.ID, song.Metadata.Song, song.Metadata.Artist, song.Metadata.Length, song.Djname, song.Djid)

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
		msg := fmt.Sprintf("Hey @%s, welcome! :)", user.Name)
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
	b.RemoveDjEscorting(b.Room.Song.djId)

	if b.UserIsModerator(b.Config.UserId) && b.Config.ModQueue {
		b.Queue.Remove(u.ID)
	}

	if b.Config.AutoDj && b.Room.djs.Size() == 0 {
		b.AutoDj()
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

	if b.UserIsModerator(b.Config.UserId) && b.Config.ModQueue && !b.Queue.CheckReservation(u.ID) {
		b.EscortDj(u.ID)
		msg := fmt.Sprintf("Sorry to remove you, @%s. But queue is active and the available slot is reserved. You can join the queue by typing !q+", u.Name)
		b.RoomMessage(msg)
		return
	}

	b.Room.AddDj(u.Userid)

	logrus.WithFields(logrus.Fields{
		"userId":   u.Userid,
		"userName": u.Name,
		"fans":     u.Fans,
		"points":   u.Points,
	}).Info("STAGE:DJ_JOINED")
}

func onRemDj(b *Bot, e ttapi.RemDJEvt) {
	u := e.User[0]
	b.Room.RemoveDj(u.Userid)
	b.RemoveDjEscorting(u.Userid)

	if b.Config.AutoDj && u.Userid == b.Config.UserId && e.Modid != "" {
		b.ToggleAutoDj()
		return
	}

	if b.Config.AutoDj && b.Room.djs.Size() == 0 {
		b.AutoDj()
	}

	logrus.WithFields(logrus.Fields{
		"userId":    u.Userid,
		"userName":  u.Name,
		"moderator": e.Modid,
	}).Info("STAGE:DJ_LEFT")
}
