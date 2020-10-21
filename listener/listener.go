package listener

import (
	"log"
	"strconv"
	"strings"
	"time"

	"beanbot/state"

	"github.com/bwmarrin/discordgo"
)

const (
	DISCORD_EPOCH = 1420070400000
)

var ltcChans map[string]*state.ChannelData

func init() {
	ltcChans = state.GetAllServers()
	// Ensure we remain at zero each reset (not including the beans that haven't popped)
	for _, data := range ltcChans {
		totalBeans := state.GetServerSum(data.ServerID)
		n := data.MostRecentNumber
		expectedBeans := (n * (n + 1)) / 2
		if totalBeans != expectedBeans {
			balance := expectedBeans - totalBeans
			// If we are giving + beans, give them to the poorest
			if balance > 0 {
				people, err := state.UglyBeanLeaderboard(data.ServerID, true, 1)
				if err != nil {
					continue
				}
				if len(people) < 1 {
					continue
				}
				userID := people[0].User
				state.AddBeans(data.ServerID, userID, balance)
				// Otherwise, take from the richest
			} else {
				people, err := state.UglyBeanLeaderboard(data.ServerID, false, 1)
				if err != nil {
					continue
				}
				if len(people) < 1 {
					continue
				}
				userID := people[0].User
				state.AddBeans(data.ServerID, userID, balance)
			}
		}
	}
}

// Loads all the messages sent after `mostrecentid`, called on startup
// to catch all of the messages we missed while offline.
// =======================FIXME=============================
// This assumes there won't be 100 messages sent in learn to count before the bot turns back on --
// This should be changed, but given that the bot will usually only be offline for a few seconds,
// It's good enough for now
func LoadUnseenMessages(s *discordgo.Session) {
	for _, data := range ltcChans {
		messages, err := s.ChannelMessages(data.ChannelID, 100, "", data.MostRecentID, "")
		if err != nil {
			continue
		}
		// Put all messages into the queue and let them be taken care of automatically
		for _, m := range messages {
			// Don't parse bot messages during this refresh
			if m.Author.Bot {
				continue
			}
			m.GuildID = data.ServerID
			EvaluateMessage(s, &discordgo.MessageCreate{m})
		}
	}
}

func GetServerData(guildID string) *state.ChannelData {
	for _, serverData := range ltcChans {
		if serverData.ServerID == guildID {
			return serverData
		}
	}
	return nil
}

func UpdateLTCChannel(serverID, channelID, messageID string) *state.ChannelData {
	// Unregister the old channel if needed
	var oldServerData *state.ChannelData
	var oldChannel string
	for channel, data := range ltcChans {
		if data.ServerID == serverID {
			oldServerData = data
			oldChannel = channel
			break
		}
	}
	if oldServerData != nil {
		delete(ltcChans, oldChannel)
		// Previous author will be the same, as it should.
		// Ensure we dont read old messages when a channel becomes LTC
		oldServerData.MostRecentID = messageID
		ltcChans[channelID] = oldServerData
	} else {
		ltcChans[channelID] = &state.ChannelData{
			ServerID:              serverID,
			ChannelID:             channelID,
			MostRecentID:          messageID,
			MostRecentAuthorID:    "",
			MostRecentNumber:      -1,
			HighestNumberAchieved: -1,
		}
	}
	return ltcChans[channelID]
}

func timeFromID(id string) (t time.Time, err error) {
	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return
	}
	n = (n >> 22) + DISCORD_EPOCH
	return time.Unix(0, int64(time.Millisecond)*n), nil
}

func EvaluateMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if _, ok := ltcChans[m.ChannelID]; !ok {
		return
	}

	thisNumber, err := strconv.Atoi(strings.TrimSpace(m.Content))

	// -1 will never a correct answer
	if err != nil {
		thisNumber = -1
	}

	ts, err := timeFromID(m.ID)
	// Something is very wrong if this happens. This could probably be safely ignored.
	if err != nil {
		log.Panic("HCF BECAUSE DISCORD DOESN'T SEND ME A NUMBER AS AN ID")
		return
	}

	AddMessageToQueue(&state.MessageData{
		Session:   s,
		GuildID:   m.GuildID,
		ChannelID: m.ChannelID,
		AuthorID:  m.Author.ID,
		ID:        m.ID,
		Number:    thisNumber,
		Timestamp: ts,
	})

}
