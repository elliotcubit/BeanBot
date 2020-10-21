package listener

import (
	"container/heap"
	"fmt"
	"log"
	"time"

	"beanbot/state"
	"github.com/bwmarrin/discordgo"
)

var mq MessageHeap

// Initialize the message heap and start the periodic parser
func init() {
	heap.Init(&mq)
	go PeriodicFlush()
}

func AddMessageToQueue(m *state.MessageData) {
	mq.s.Lock()
	defer mq.s.Unlock()
	heap.Push(&mq, m)
}

func SafeProcessMessages() []*state.MessageData {
	mq.s.Lock()
	defer mq.s.Unlock()
	return UnsafeProcessMessages()
}

func UnsafeProcessMessages() []*state.MessageData {
	messages := (&mq).Flush()
	for _, v := range messages {
		ProcessMessage(v)
	}
	return messages
}

// Called on bot shutdown to flush the rest of the messages from the queue
func CleanupHeap() []*state.MessageData {
	// Never unlock. This function is called at program end.
	mq.s.Lock()
	time.Sleep(MESSAGE_QUEUE_FLUSH_INTERVAL * time.Millisecond)
	return UnsafeProcessMessages()
}

func ProcessMessage(m *state.MessageData) {
	context := ltcChans[m.ChannelID]

	var isCorrect bool

	// Only 1 or the next number is correct
	if context.MostRecentNumber == -1 {
		isCorrect = (m.Number == 1)
		// Do not "punish" mistakes before we've restarted counting
		if !isCorrect {
			return
		}
	} else {
		isCorrect = (m.Number == context.MostRecentNumber+1)
	}

	// No double counting
	if context.MostRecentAuthorID == m.AuthorID {
		isCorrect = false
	}

	if isCorrect {
		_, err := state.AddBeans(m.GuildID, m.AuthorID, m.Number)
		if err != nil {
			return
		}
		// Mark ourselves as parsed
		ltcChans[m.ChannelID].MostRecentNumber = m.Number
		ltcChans[m.ChannelID].MostRecentAuthorID = m.AuthorID
		ltcChans[m.ChannelID].MostRecentID = m.ID
		// Set highest number locally if needed
		if m.Number > ltcChans[m.ChannelID].HighestNumberAchieved {
			ltcChans[m.ChannelID].HighestNumberAchieved = m.Number
		}
	} else {
		num := ltcChans[m.ChannelID].MostRecentNumber
		// Subtract the sum of all awarded beans to keep the net at zero
		amount := (num * (num + 1)) / 2
		_, err := state.AddBeans(m.GuildID, m.AuthorID, -amount)
		if err != nil {
			return
		}

		// Back to original state
		ltcChans[m.ChannelID].MostRecentNumber = -1
		ltcChans[m.ChannelID].MostRecentAuthorID = ""
		ltcChans[m.ChannelID].MostRecentID = m.ID

		embed := &discordgo.MessageEmbed{Color: 0x3498DB}
		embed.Title = "Uh Oh"
		embed.Description = fmt.Sprintf("<@%s> spilled and lost %d beans!", m.AuthorID, amount)
		m.Session.ChannelMessageSendEmbed(m.ChannelID, embed)

		// Update the message for when we send it to the db
		m.Number = -1   // Recover a failed state on reset
		m.AuthorID = "" // Allow doublecounting after a mistake
	}

	err := state.UpdateMessageInServer(m)
	if err != nil {
		log.Printf("Error when updating message: %+v\n", err)
	}
}

func PeriodicFlush() {
	for {
		time.Sleep(MESSAGE_QUEUE_FLUSH_INTERVAL)
		SafeProcessMessages()
	}
}
