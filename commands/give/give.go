package give

import (
	"strconv"
	"strings"

	"beanbot/handlers"
	"beanbot/state"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterActiveModule(
		Give{},
	)
}

type Give struct{}

// !bean give 50 @Someone
func (h Give) Do(s *discordgo.Session, m *discordgo.MessageCreate) {
	data := strings.SplitN(m.Content, " ", 4)
	if len(data) < 4 {
		return
	}
	if len(m.Mentions) < 1 {
		return
	}
	amount := 0
	amount, err := strconv.Atoi(data[2])
	if err != nil {
		return
	}
	if amount <= 0 {
		s.ChannelMessageSend(m.ChannelID, "You must send at least one bean.")
		return
	}

	recipientID := m.Mentions[0].ID
	donatorID := m.Author.ID

	if m.Mentions[0].Bot && m.Mentions[0].ID != s.State.User.ID {
		s.ChannelMessageSend(m.ChannelID, "You cannot give beans to other bots.")
		return
	}

	donatorBalance, err := state.GetUserBalance(m.GuildID, donatorID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Bean transfer failed.")
		return
	}
	if donatorBalance < amount {
		s.ChannelMessageSend(m.ChannelID, "You do not have enough beans.")
		return
	}
	err = state.TransferBeans(
		m.GuildID,
		donatorID,
		recipientID,
		amount,
	)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Bean transfer failed.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Bean transfer complete.")
	}
}

func (h Give) Prefixes() []string {
	return []string{"give"}
}
