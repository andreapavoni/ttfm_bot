package main

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
	"github.com/andreapavoni/ttfm_bot/ttfm/commands"
)

var Version = "v0.0.0-dev"

func main() {
	b := ttfm.New()
	b.Commands.Add("dj", commands.DjCommand())
	b.Commands.Add("escort", commands.EscortCommand())
	b.Commands.Add("escortme", commands.EscortMeCommand())
	b.Commands.Add("autobop", commands.AutoBopCommand())
	b.Commands.Add("autodj", commands.AutoDjCommand())
	b.Commands.Add("autosnag", commands.AutoSnagCommand())
	b.Commands.Add("bop", commands.BopCommand())
	b.Commands.Add("dj", commands.DjCommand())
	b.Commands.Add("snag", commands.SnagCommand())
	b.Commands.Add("fan", commands.FanCommand())
	b.Commands.Add("unfan", commands.UnfanCommand())
	b.Commands.Add("props", commands.PropsCommand())
	b.Commands.Add("skip", commands.SkipCommand())
	b.Commands.Add("boot", commands.BootCommand())
	b.Commands.Add("padd", commands.PlaylistCreateCommand())
	b.Commands.Add("pdel", commands.PlaylistDeleteCommand())
	b.Commands.Add("pls", commands.PlaylistListCommand())
	b.Commands.Add("prm", commands.PlaylistRemoveSongCommand())
	b.Commands.Add("pc", commands.PlaylistSwitchCommand())
	b.Commands.Add("q", commands.QueueCommand())
	b.Commands.Add("qadd", commands.QueueAddCommand())
	b.Commands.Add("qrm", commands.QueueRemoveCommand())
	b.Commands.Add("cfg", commands.SetConfigCommand())
	b.Commands.Add("say", commands.SayCommand())
	b.Commands.Add("help", commands.HelpCommand())
	b.Commands.Add("r", commands.ReactionCommand())
	b.Commands.Add("room", commands.FavRoomCommand())
	b.Start()
}
