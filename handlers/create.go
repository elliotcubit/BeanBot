package handlers

import (
	"beanbot/commands"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var activeModules = []commands.ActiveModule{}

func RegisterActiveModule(handler commands.ActiveModule) {
	activeModules = append(activeModules, handler)
}

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bots
	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, "!bean") {
		return
	}

	// Supply the help command if there is no argument
	data := strings.SplitN(m.Content, " ", 3)
	if len(data) < 2 {
		data = append(data, "help")
	}

	for _, handler := range activeModules {
		for _, prefix := range handler.Prefixes() {
			if data[1] == prefix {
				handler.Do(s, m)
				return
			}
		}
	}

	// TODO passive listener
}
