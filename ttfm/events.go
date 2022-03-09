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
	if b.Config.AutoShowSongStats {
		b.ShowSongStats()
	}

	if b.escorting.HasElement(b.Room.Song.djId) {
		b.EscortDj(b.Room.Song.djId)
		b.RemoveDjEscorting(b.Room.Song.djId)
	}

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

	utils.ExecuteDelayedRandom(30, b.Bop)

	if b.Config.AutoSnag {
		utils.ExecuteDelayedRandom(30, func() {
			b.Snag(b.Room.Song.Id)
		})
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
	botUser, _ := b.UserFromId(b.Config.UserId)

	if b.Config.ModAutoWelcome && b.UserIsModerator(botUser) {
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

	logrus.WithFields(logrus.Fields{
		"userId":   u.ID,
		"userName": u.Name,
		"fans":     u.Fans,
		"points":   u.Points,
	}).Info("ROOM:USER_LEFT")
}

func onAddDj(b *Bot, e ttapi.AddDJEvt) {
	u := e.User[0]
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
