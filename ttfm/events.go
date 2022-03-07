package ttfm

import (
	"fmt"

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

	b.room.Update(roomInfo)

	if b.config.SetBot {
		utils.ExecuteDelayedRandom(30, func() { b.api.SetBot() })
	}

	b.LoadPlaylists()
	b.SwitchPlaylist(b.config.CurrentPlaylist)
}

func onRoomChanged(b *Bot, e ttapi.RoomInfoRes) {
	logrus.WithFields(logrus.Fields{
		"room":       e.Room.Name,
		"moderators": e.Room.Metadata.ModeratorID,
		"downvotes":  e.Room.Metadata.Downvotes,
		"upvotes":    e.Room.Metadata.Upvotes,
		"maxdjs":     e.Room.Metadata.MaxDjs,
		"djcount":    e.Room.Metadata.Djcount,
		"djs":        e.Room.Metadata.Djs,
	}).Info("BOT:ROOM_CHANGED")

	b.room.Update(e)
	utils.ExecuteDelayedRandom(15, b.Bop)

	logrus.WithFields(logrus.Fields{
		"room":      b.room.name,
		"roomId":    b.room.id,
		"shortcut":  b.room.shortcut,
		"djs":       b.room.djs.Size(),
		"listeners": e.Room.Metadata.Listeners,
	}).Info("BOT:ROOM_CHANGED updated room data")

	if b.config.AutoDj && e.Room.Metadata.Djcount == 0 {
		b.AutoDj()
	}
}

func onNewSong(b *Bot, e ttapi.NewSongEvt) {
	if b.config.AutoShowSongStats {
		b.ShowSongStats()
	}

	if b.escorting.HasElement(b.room.song.djId) {
		b.EscortDj(b.room.song.djId)
		b.RemoveDjEscorting(b.room.song.djId)
	}

	if b.room.song.djId == b.config.UserId {
		b.PushSongBottomPlaylist()
	}

	logrus.WithFields(logrus.Fields{
		"dj":     b.room.song.djName,
		"djID":   b.room.song.djId,
		"song":   b.room.song.title,
		"artist": b.room.song.artist,
		"length": b.room.song.length,
		"up":     b.room.song.up,
		"down":   b.room.song.down,
		"snag":   b.room.song.snag,
	}).Info("ROOM:LAST_SONG_STATS")

	b.room.UpdateModerators(e.Room.Metadata.ModeratorID)
	song := e.Room.Metadata.CurrentSong
	b.room.song.Reset(song.ID, song.Metadata.Song, song.Metadata.Artist, song.Metadata.Length, song.Djname, song.Djid)

	logrus.WithFields(logrus.Fields{
		"dj":     song.Djname,
		"djId":   song.Djid,
		"song":   song.Metadata.Song,
		"artist": song.Metadata.Artist,
		"length": song.Metadata.Length,
	}).Info("ROOM:NEW_SONG")

	utils.ExecuteDelayedRandom(20, b.Bop)

	if b.config.AutoSnag {
		utils.ExecuteDelayedRandom(10, func() {
			b.Snag(b.room.song.id)
		})
	}
}

func onUpdateVotes(b *Bot, e ttapi.UpdateVotesEvt) {
	b.room.song.UpdateStats(e.Room.Metadata.Upvotes, e.Room.Metadata.Downvotes, b.room.song.snag)
	userId, vote := b.room.song.UnpackVotelog(e.Room.Metadata.Votelog)
	user, _ := b.UserFromId(userId)

	logrus.WithFields(logrus.Fields{
		"up":        e.Room.Metadata.Upvotes,
		"down":      e.Room.Metadata.Downvotes,
		"listeners": e.Room.Metadata.Listeners,
		"userID":    userId,
		"vote":      vote,
		"name":      user.Name,
	}).Info("SONG:VOTE")
}

func onSnagged(b *Bot, e ttapi.SnaggedEvt) {
	b.room.song.UpdateStats(b.room.song.up, b.room.song.down, b.room.song.snag+1)
	user, _ := b.UserFromId(e.UserID)

	logrus.WithFields(logrus.Fields{
		"userID": e.UserID,
		"name":   user.Name,
		"roomID": e.RoomID,
	}).Info("SONG:SNAG")
}

func onRegistered(b *Bot, e ttapi.RegisteredEvt) {
	u := e.User[0]
	if u.ID == b.config.UserId {
		return
	}

	b.room.AddUser(u.ID, u.Name)

	user, _ := b.UserFromId(u.ID)
	botUser, _ := b.UserFromId(b.config.UserId)

	if b.config.ModAutoWelcome && b.UserIsModerator(botUser) {
		msg := fmt.Sprintf("Hey @%s, welcome :)", user.Name)
		b.RoomMessage(msg)
	}

	logrus.WithFields(logrus.Fields{
		"userID": u.ID,
		"name":   u.Name,
		"laptop": u.Laptop,
		"fans":   u.Fans,
		"points": u.Points,
		"avatar": u.Avatarid,
	}).Info("ROOM:USER_JOINED")
}

func onDeregistered(b *Bot, e ttapi.DeregisteredEvt) {
	u := e.User[0]
	if u.ID == b.config.UserId {
		return
	}

	b.room.RemoveDj(u.ID)
	b.room.RemoveUser(u.ID)
	b.RemoveDjEscorting(b.room.song.djId)

	logrus.WithFields(logrus.Fields{
		"userID": u.ID,
		"name":   u.Name,
		"fans":   u.Fans,
		"points": u.Points,
	}).Info("ROOM:USER_LEFT")
}

func onAddDj(b *Bot, e ttapi.AddDJEvt) {
	u := e.User[0]
	b.room.AddDj(u.Userid)

	logrus.WithFields(logrus.Fields{
		"userID": u.Userid,
		"name":   u.Name,
	}).Info("STAGE:DJ_JOINED")
}

func onRemDj(b *Bot, e ttapi.RemDJEvt) {
	u := e.User[0]
	b.room.RemoveDj(u.Userid)
	b.RemoveDjEscorting(u.Userid)

	if b.config.AutoDj && u.Userid == b.config.UserId && e.Modid != "" {
		b.ToggleAutoDj()
		return
	}

	if b.config.AutoDj && b.room.djs.Size() == 0 {
		b.AutoDj()
	}

	logrus.WithFields(logrus.Fields{
		"userID":    u.Userid,
		"name":      u.Name,
		"moderator": e.Modid,
	}).Info("STAGE:DJ_LEFT")
}
