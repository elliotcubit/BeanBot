package max

import (
	"fmt"

	"beanbot/handlers"
	"beanbot/listener"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterActiveModule(
		Risk{},
	)
}

type Risk struct{}

func (h Risk) Do(s *discordgo.Session, m *discordgo.MessageCreate) {
	data := listener.GetServerData(m.ChannelID)
	amt := 0
	if data != nil {
		n := data.MostRecentNumber
		amt = (n * (n + 1)) / 2
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You will lose %d beans if you make a mistake right now.", amt))
}

func (h Risk) Prefixes() []string {
	return []string{"risk"}
}
