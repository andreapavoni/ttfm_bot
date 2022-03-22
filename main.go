package main

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
	"github.com/andreapavoni/ttfm_bot/ttfm/commands"
)

var Version = "v0.0.0-dev"

func main() {
	b := ttfm.New()
	b.AddCommand("dj", commands.DjCommand())
	b.AddCommand("escort", commands.EscortCommand())
	b.AddCommand("escortme", commands.EscortMeCommand())
	b.AddCommand("autobop", commands.AutoBopCommand())
	b.AddCommand("autodj", commands.AutoDjCommand())
	b.AddCommand("autosnag", commands.AutoSnagCommand())
	b.AddCommand("bop", commands.BopCommand())
	b.AddCommand("dj", commands.DjCommand())
	b.AddCommand("snag", commands.SnagCommand())
	b.AddCommand("fan", commands.FanCommand())
	b.AddCommand("unfan", commands.UnfanCommand())
	b.AddCommand("props", commands.PropsCommand())
	b.AddCommand("skip", commands.SkipCommand())
	b.AddCommand("boot", commands.BootCommand())
	b.AddCommand("padd", commands.PlaylistCreateCommand())
	b.AddCommand("pdel", commands.PlaylistDeleteCommand())
	b.AddCommand("pls", commands.PlaylistListCommand())
	b.AddCommand("prm", commands.PlaylistRemoveSongCommand())
	b.AddCommand("pc", commands.PlaylistSwitchCommand())
	b.AddCommand("q", commands.QueueCommand())
	b.AddCommand("qadd", commands.QueueAddCommand())
	b.AddCommand("qrm", commands.QueueRemoveCommand())
	b.AddCommand("cfg", commands.SetConfigCommand())
	b.AddCommand("say", commands.SayCommand())
	b.AddCommand("help", commands.HelpCommand())
	b.AddCommand("r", commands.ReactionCommand())
	b.AddCommand("room", commands.RoomCommand())
	b.Start()
}
