package max

import (
	"fmt"

	"beanbot/handlers"
	"beanbot/listener"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterActiveModule(
		Max{},
	)
}

type Max struct{}

func (h Max) Do(s *discordgo.Session, m *discordgo.MessageCreate) {
	data := listener.GetServerData(m.GuildID)
	amt := 0
	if data != nil {
		amt = data.HighestNumberAchieved
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("This server has counted to %d.", amt))
}

func (h Max) Prefixes() []string {
	return []string{"max"}
}
