package leaderboard

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"beanbot/handlers"
	"beanbot/state"
	"github.com/bwmarrin/discordgo"
)

var challenges []*Challenge

func init() {
	handlers.RegisterActiveModule(
		Leaderboard{},
	)
}

type Leaderboard struct{}

// !bean leaderboard top 5
func (h Leaderboard) Do(s *discordgo.Session, m *discordgo.MessageCreate) {
	data := strings.SplitN(m.Content, " ", 4)
	var err error
	direction := "top"
	amount := 5
	// Check which side of the leaderboard we're reading
	if len(data) > 2 {
		if data[2] == "top" || data[2] == "bottom" {
			direction = data[2]
		}
	}
	// Check if number was specified
	if len(data) > 3 {
		amount, err = strconv.Atoi(data[1])
		if err != nil {
			amount = 5
		}
	}
	results, err := state.GetBeanLeaderboard(m.GuildID, direction, amount)
	if err != nil {
		log.Println(err)
		return
	}
	out := "```"
	for ranking, data := range results {
		out += fmt.Sprintf("%-2d | %-32s %8d beans\n", ranking+1, data.User, data.Amount)
	}
	out += "```"
	s.ChannelMessageSend(m.ChannelID, out)
}

func (h Leaderboard) Prefixes() []string {
	return []string{"leaderboard"}
}
