package query

import (
	"fmt"

	"beanbot/handlers"
	"beanbot/state"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterActiveModule(
		Query{},
	)
}

type Query struct{}

func (h Query) Do(s *discordgo.Session, m *discordgo.MessageCreate) {
	var user *discordgo.User
	if len(m.Mentions) < 1 {
		user = m.Author
	} else {
		user = m.Mentions[0]
	}
	amount, err := state.GetUserBalance(m.GuildID, user.ID)
	if err != nil {
		return
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %d beans", user.String(), amount))
}

func (h Query) Prefixes() []string {
	return []string{"query"}
}
