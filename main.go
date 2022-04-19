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
	b.Commands.Add("bop", commands.BopCommand())
	b.Commands.Add("dj", commands.DjCommand())
	b.Commands.Add("snag", commands.SnagCommand())
	b.Commands.Add("fan", commands.FanCommand())
	b.Commands.Add("unfan", commands.UnfanCommand())
	b.Commands.Add("props", commands.PropsCommand())
	b.Commands.Add("skip", commands.SkipCommand())
	b.Commands.Add("boot", commands.BootCommand())
	b.Commands.Add("q", commands.QueueCommand())
	b.Commands.Add("cfg", commands.SetConfigCommand())
	b.Commands.Add("say", commands.SayCommand())
	b.Commands.Add("help", commands.HelpCommand())
	b.Commands.Add("r", commands.ReactionCommand())
	b.Commands.Add("pl", commands.PlaylistCommand())
	b.Commands.Add("room", commands.FavRoomCommand())
	b.Commands.Add("die", commands.KillSwitchCommand())
	b.Commands.Add("admin", commands.AdminCommand())
	b.Start()
}
