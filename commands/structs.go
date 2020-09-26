package commands

import (
	"github.com/bwmarrin/discordgo"
)

type ActiveModule interface {
	// Do what the module Do
	Do(s *discordgo.Session, m *discordgo.MessageCreate)
	// An array of each !bean [command] that should execute this module
	Prefixes() []string
}
