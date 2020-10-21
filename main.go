package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"beanbot/handlers"
	"beanbot/listener"

	"github.com/bwmarrin/discordgo"

	// Importing a command module here takes care of everything

	// Active Modules
	_ "beanbot/commands/bet"
	_ "beanbot/commands/configure"
	_ "beanbot/commands/give"
	_ "beanbot/commands/help"
	_ "beanbot/commands/leaderboard"
	_ "beanbot/commands/max"
	_ "beanbot/commands/query"
	_ "beanbot/commands/risk"
)

func main() {
	log.Println("Loading golordsbot")

	DISCORD_TOKEN := os.Getenv("DISCORD_TOKEN")

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

	// Open connection
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error while opening connection to discord: %v", err)
	}

	rand.Seed(time.Now().Unix())

	// Load messages sent since we were offline in each server if applicable
	listener.LoadUnseenMessages(dg)

	log.Println("Golords bot is alive. SIGINT exits.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Ensure we don't quit in the middle of a message being parsed due to bad timing
	listener.CleanupHeap()

	log.Println("SIGINT Registered. Shutting down.")
	log.Println("Goodbye <3")
}
