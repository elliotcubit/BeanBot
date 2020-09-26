package lottery

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"beanbot/handlers"
	"beanbot/state"

	"github.com/bwmarrin/discordgo"
)

const ticketPrice = 500

var runningLotteries map[string]*Lottery

func init() {
	handlers.RegisterActiveModule(
		Lottery{},
	)
	runningLotteries = make(map[string]*Lottery, 0)
}

type Lottery struct {
	ChannelID string
	Tickets   []string
}

func (h Lottery) Do(s *discordgo.Session, m *discordgo.MessageCreate) {
	data := strings.SplitN(m.Content, " ", 3)
	timer := 15
	if len(data) > 2 {
		timer, err := strconv.Atoi(data[2])
		if err != nil {
			timer = 15
		}
	}
	if _, exists := runningLotteries[m.GuildID]; exists {
		s.ChannelMessageSend(m.ChannelID, "There is already a running lottery in this server.")
		return
	}
	newLottery := &Lottery{
		ChannelID: m.ChannelID,
		Tickets:   []string{},
	}
	runningLotteries[m.GuildID] = newLottery
	go newLottery.AwaitBeanLottery(s, m, timer)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Lottery has been started and will end in %d minutes. Enter with !bean ticket", timer))
}

func (h Lottery) Prefixes() []string {
	return []string{"lottery"}
}
