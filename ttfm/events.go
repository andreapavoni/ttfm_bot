package ttfm

import (
	"github.com/alaingilbert/ttapi"
	"github.com/sirupsen/logrus"
)

func onReady(b *Bot) {
	b.Actions.SetBot()
	b.Actions.InitPlaylists()
	logrus.Info("BOT:READY")
}

func onRoomChanged(b *Bot, e ttapi.RoomInfoRes) {
	b.Actions.UpdateRoom(e)
	b.Actions.AutoBop()
	b.Actions.AutoDj()

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
	b.Actions.ShowSongStats()
	logrus.WithFields(logrus.Fields{
		"djName": b.Room.Song.DjName,
		"djId":   b.Room.Song.DjId,
		"title":  b.Room.Song.Title,
		"artist": b.Room.Song.Artist,
		"length": b.Room.Song.Length,
		"up":     b.Room.Song.up,
		"down":   b.Room.Song.down,
		"snag":   b.Room.Song.snag,
	}).Info("ROOM:LAST_SONG_STATS")

	b.Actions.EscortDjs()
	b.Actions.ForwardQueue()
	// when bot is djing, push the last played song to bottom of its playlist
	b.Actions.ShiftPlaylistSong()
	b.Actions.UpdateRoomFromApi()
	// enforce song duration to avoid trolls with 2hours tracks
	b.Actions.EnforceSongDuration()
	b.Actions.AutoBop()
	b.Actions.AutoSnag()

	logrus.WithFields(logrus.Fields{
		"djName": b.Room.Song.DjName,
		"djId":   b.Room.Song.DjId,
		"title":  b.Room.Song.Title,
		"artist": b.Room.Song.Artist,
		"length": b.Room.Song.Length,
	}).Info("ROOM:NEW_SONG")
}

func onUpdateVotes(b *Bot, e ttapi.UpdateVotesEvt) {
	b.Actions.UpdateSongStats(e.Room.Metadata.Upvotes, e.Room.Metadata.Downvotes, b.Room.Song.snag)
	logrus.WithFields(logrus.Fields{
		"up":        e.Room.Metadata.Upvotes,
		"down":      e.Room.Metadata.Downvotes,
		"listeners": e.Room.Metadata.Listeners,
	}).Info("SONG:VOTE")
}

func onSnagged(b *Bot, e ttapi.SnaggedEvt) {
	b.Actions.UpdateSongStats(b.Room.Song.up, b.Room.Song.down, b.Room.Song.snag+1)

	logrus.WithFields(logrus.Fields{
		"userId": e.UserID,
		"roomId": e.RoomID,
	}).Info("SONG:SNAG")
}

func onRegistered(b *Bot, e ttapi.RegisteredEvt) {
	u := e.User[0]
	b.Actions.RegisterUser(u.ID)

	logrus.WithFields(logrus.Fields{
		"userId":   u.ID,
		"userName": u.Name,
		"fans":     u.Fans,
		"points":   u.Points,
	}).Info("ROOM:USER_JOINED")
}

func onDeregistered(b *Bot, e ttapi.DeregisteredEvt) {
	u := e.User[0]
	b.Actions.UnregisterUser(u.ID)
	b.Actions.AutoDj()

	logrus.WithFields(logrus.Fields{
		"userId":   u.ID,
		"userName": u.Name,
		"fans":     u.Fans,
		"points":   u.Points,
	}).Info("ROOM:USER_LEFT")
}

func onAddDj(b *Bot, e ttapi.AddDJEvt) {
	u := e.User[0]
	b.Actions.AddDj(u.Userid)
	b.Actions.EnforceQueueStageReservation(u.ID)
	b.Actions.ConsiderQueueActivation()

	logrus.WithFields(logrus.Fields{
		"userId":   u.Userid,
		"userName": u.Name,
		"fans":     u.Fans,
		"points":   u.Points,
	}).Info("STAGE:DJ_JOINED")
}

func onRemDj(b *Bot, e ttapi.RemDJEvt) {
	u := e.User[0]
	b.Actions.RemoveDj(u.Userid, e.Modid)
	b.Actions.AutoDj()
	b.Actions.ForwardQueue()

	logrus.WithFields(logrus.Fields{
		"userId":    u.Userid,
		"userName":  u.Name,
		"moderator": e.Modid,
	}).Info("STAGE:DJ_LEFT")
}
