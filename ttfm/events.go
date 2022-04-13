package ttfm

import (
	"github.com/alaingilbert/ttapi"
	"github.com/andreapavoni/ttfm_bot/utils"
	"github.com/sirupsen/logrus"
)

func onReady(b *Bot) {
	b.Actions.SetBot()
	b.Actions.LoadBotIdentity()
	b.Actions.LoadMainAdmin()
	b.Actions.InitPlaylists()
	logrus.Info("BOT:READY")
}

func onRoomChanged(b *Bot, e ttapi.RoomInfoRes) {
	utils.MaybeLogError("BOT:ROOM:UPDATE", func () error {return b.Actions.UpdateRoom(e)})
	b.Actions.AutoBop()
	b.Actions.ConsiderStartAutoDj()

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
	b.Actions.UpdateDjStatsPlays()
	b.Actions.ShowSongStats()
	logrus.WithFields(logrus.Fields{
		"djName": b.Room.CurrentSong.DjName,
		"djId":   b.Room.CurrentSong.DjId,
		"title":  b.Room.CurrentSong.Title,
		"artist": b.Room.CurrentSong.Artist,
		"length": b.Room.CurrentSong.Length,
		"up":     b.Room.CurrentSong.up,
		"down":   b.Room.CurrentSong.down,
		"snag":   b.Room.CurrentSong.snag,
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
	b.Actions.ConsiderStartAutoDj()

	logrus.WithFields(logrus.Fields{
		"songSourceId": e.Room.Metadata.CurrentSong.Sourceid,
		"songSource":   e.Room.Metadata.CurrentSong.Source,
		"songId":       e.Room.Metadata.CurrentSong.ID,
		"djName":       b.Room.CurrentSong.DjName,
		"djId":         b.Room.CurrentSong.DjId,
		"title":        b.Room.CurrentSong.Title,
		"artist":       b.Room.CurrentSong.Artist,
		"length":       b.Room.CurrentSong.Length,
	}).Info("ROOM:NEW_SONG")
}

func onUpdateVotes(b *Bot, e ttapi.UpdateVotesEvt) {
	b.Actions.UpdateSongStats(e.Room.Metadata.Upvotes, e.Room.Metadata.Downvotes, b.Room.CurrentSong.snag)
	userId, vote := b.Actions.UnpackVotelog(e.Room.Metadata.Votelog)
	b.Actions.UpdateDjStatsVote(vote)

	logrus.WithFields(logrus.Fields{
		"up":        e.Room.Metadata.Upvotes,
		"down":      e.Room.Metadata.Downvotes,
		"listeners": e.Room.Metadata.Listeners,
		"vote":      vote,
		"userId":    userId,
	}).Info("SONG:VOTE")
}

func onSnagged(b *Bot, e ttapi.SnaggedEvt) {
	b.Actions.UpdateSongStats(b.Room.CurrentSong.up, b.Room.CurrentSong.down, b.Room.CurrentSong.snag+1)
	b.Actions.UpdateDjStatsSnag()

	logrus.WithFields(logrus.Fields{
		"userId": e.UserID,
		"roomId": e.RoomID,
	}).Info("SONG:SNAG")
}

func onRegistered(b *Bot, e ttapi.RegisteredEvt) {
	u := e.User[0]
	b.Actions.RegisterUser(u.ID, u.Name)
	b.Actions.ConsiderStartAutoDj()

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
	b.Actions.ConsiderStartAutoDj()

	logrus.WithFields(logrus.Fields{
		"userId":   u.ID,
		"userName": u.Name,
		"fans":     u.Fans,
		"points":   u.Points,
	}).Info("ROOM:USER_LEFT")
}

func onBootedUser(b *Bot, e ttapi.BootedUserEvt) {
	user, _ := b.Users.UserFromId(e.Userid)
	b.Actions.UnregisterUser(user.Id)
	b.Actions.ConsiderStartAutoDj()

	logrus.WithFields(logrus.Fields{
		"userId":    e.Userid,
		"userName":  user.Name,
		"moderator": e.Modid,
		"reason":    e.Reason,
	}).Info("ROOM:USER_BOOTED")
}

func onAddDj(b *Bot, e ttapi.AddDJEvt) {
	u := e.User[0]
	b.Actions.AddDj(u.Userid)
	b.Actions.EnforceQueueStageReservation(u.ID)
	b.Actions.ConsiderQueueStart()

	logrus.WithFields(logrus.Fields{
		"userId":   u.Userid,
		"userName": u.Name,
		"fans":     u.Fans,
		"points":   u.Points,
	}).Info("STAGE:DJ_JOINED")
}

func onRemDj(b *Bot, e ttapi.RemDJEvt) {
	u := e.User[0]
	b.Actions.ShowDjStats(u.Userid)
	b.Actions.RemoveDj(u.Userid, e.Modid)
	b.Actions.ConsiderQueueStop()
	b.Actions.ForwardQueue()
	b.Actions.ConsiderStartAutoDj()

	logrus.WithFields(logrus.Fields{
		"userId":    u.Userid,
		"userName":  u.Name,
		"moderator": e.Modid,
	}).Info("STAGE:DJ_LEFT")
}
