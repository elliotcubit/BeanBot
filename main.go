package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"beanbot/handlers"
	"github.com/bwmarrin/discordgo"

	// Importing a module here handles everything

	// Active Modules
	_ "beanbot/commands/bet"
	_ "beanbot/commands/leaderboard"
	// Passive Modules
)

func main() {
	log.Println("Loading golordsbot")

	// Get credentials
	DISCORD_ID := os.Getenv("DISCORD_ID")
	DISCORD_SECRET := os.Getenv("DISCORD_SECRET")
	DISCORD_TOKEN := os.Getenv("DISCORD_TOKEN")

	if DISCORD_ID == "" {
		log.Fatal("DISCORD_ID environment variable not set")
	}
	if DISCORD_SECRET == "" {
		log.Fatal("DISCORD_SECRET environment variable not set")
	}
	if DISCORD_TOKEN == "" {
		log.Fatal("DISCORD_TOKEN environment variable not set")
	}

	// Create Discord session
	dg, err := discordgo.New("Bot " + DISCORD_TOKEN)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}
	defer dg.Close()

	// # of message to respond to edited events
	dg.State.MaxMessageCount = 50

	// Register message handlers
	dg.AddHandler(handlers.OnMessageCreate)
	dg.AddHandler(handlers.OnMessageUpdate)

	// Open connection
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error while opening connection to discord: %v", err)
	}

	rand.Seed(time.Now().Unix())

	log.Println("Golords bot is alive. SIGINT exits.")
	dg.ChannelMessageSend(learntocount.LtcChan, "=====GOLORDS BOT IS ALIVE=====")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Give this message some time to send
	dg.ChannelMessageSend(learntocount.LtcChan, "===THE BOT IS GOING OFFLINE===")
	time.Sleep(5 * time.Second)

	log.Println("SIGINT Registered. Shutting down.")
	log.Println("Goodbye <3")
}
