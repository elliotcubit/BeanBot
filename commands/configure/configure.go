package configure

import (
	"fmt"
	"log"

	"beanbot/handlers"
	"beanbot/listener"
	"beanbot/state"

	"github.com/bwmarrin/discordgo"
)

func init() {
	handlers.RegisterActiveModule(
		Configure{},
	)
}

type Configure struct{}

func (h Configure) Do(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Verify the user has the Manage Server permission.
	perms, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		return
	}
	if perms&discordgo.PermissionManageChannels < 1 {
		s.ChannelMessageSend(m.ChannelID, "You must have the Manage Channels permission to make this change.")
		return
	}
	// This channel also handles updating the channel if it is already set.
	err = state.AddServerChannelToList(m.GuildID, m.ChannelID)
	if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "An internal error occured while registering this channel. Please try again later.")
		return
	}
	// Do not require a restart for this
	channelData := listener.UpdateLTCChannel(m.GuildID, m.ChannelID, m.ID)
	var start int
	if channelData.MostRecentNumber == -1 {
		start = 1
	} else {
		start = channelData.MostRecentNumber + 1
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("This channel is now registered as the learn to count channel. Counting starts at `%d` here.", start))
}

func (h Configure) Prefixes() []string {
	return []string{"configure"}
}
