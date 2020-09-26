package bet

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

var challenges []*Challenge

func init() {
	handlers.RegisterActiveModule(
		Bet{},
	)
}

type Bet struct{}

type Challenge struct {
	ServerID   string
	Challenger string
	Challengee string
	Amount     int
}

// !bean bet 50 @Someone
func (h Bet) Do(s *discordgo.Session, m *discordgo.MessageCreate) {
	data := strings.SplitN(m.Content, " ", 4)
	if len(m.Mentions) < 1 {
		s.ChannelMessageSend(m.ChannelID, "You must mention someone to make a challenge.")
		return
	}

	serverID := m.GuildID
	amount, _ := strconv.Atoi(data[2])
	challenger := m.Author.String()
	challengee := m.Mentions[0].String()

	// Check if this matches an existing challenge
	for index, challenge := range challenges {
		if challenge.ServerID == serverID &&
			challenge.Challenger == challengee &&
			challenge.Challengee == challenger {
			if amount == 0 || challenge.Amount == amount {
				s.ChannelMessageSend(m.ChannelID, executeBeanGame(index))
				return
			}
		}
	}

	if amount < 1 {
		s.ChannelMessageSend(m.ChannelID, "You cannot challenge someone for less than one bean.")
		return
	}

	challengerBalance, err := state.GetUserBalance(serverID, challenger)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "There was a problem creating your challenge.")
		return
	}
	if challengerBalance < amount {
		s.ChannelMessageSend(m.ChannelID, "You do not have enough beans to make that bet.")
		return
	}

	challengeesBalance, err := state.GetUserBalance(serverID, challengee)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "The person you're challenging does not have enough beans to make that bet.")
		return
	}

	// Verify the game doesn't already exist
	for _, challenge := range challenges {
		if challenger.ServerID == serverID &&
			challenge.Challenger == challenger &&
			challenge.Challengee == challengee {
			if challenge.Amount == amount {
				s.ChannelMessageSend(m.ChannelID, "You have already made that challenge, and it hasn't been accepted yet.")
				return
			}
		}
	}

	challenge := &Challenge{
		ServerID:   serverID,
		Challenger: challenger,
		Challengee: challengee,
		Amount:     amount,
	}
	challenges = append(challenges, challenge)
	s.ChannelMessageSend(m.ChannelID, "Challenge created for %d beans. Accept by challenging back.", amount)
}

func (h Bet) Prefixes() []string {
	return []string{"bet"}
}

func executeBeanGame(index int) string {
	challenge := challenges[index]

	// Select winner
	choice := rand.Intn(2)
	var winner string
	var loser string
	if choice == 0 {
		winner = challenge.Challenger
		loser = challenge.Challengee
	} else {
		winner = challenge.Challengee
		loser = challenge.Challenger
	}

	// Transfer funds
	_, err := state.AddToUserBalance(game.ServerID, winner, game.Amount)
	if err != nil {
		return "There was a problem finishing the challenge."
	}

	_, err = state.AddToUserBalance(game.ServerID, loser, -game.Amount)
	for err != nil {
		log.Println("Problem resolving bean challenge, trying again in 30s.")
		time.Sleep(30 * time.Second)
		_, err = state.AddToUserBalance(game.ServerID, loser, -game.Amount)
	}

	// Remove from the challenges list
	challenges[index] = challenges[len(challenges)-1]
	challenges[len(challenges)-1] = nil
	challenges = challenges[:len(challenges)-1]

	return fmt.Sprintf("%s won the bet between %s and %s for %d beans", winner, winner, loser, game.Amount)
}
