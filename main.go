package main

import (
	"dbot/internal/handler"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func createOSInterruptSingal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func main() {
	godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Failed to load DISCORD_TOKEN from environment variable")
	}

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Failed to create a discord session: %v", err)
	}

	sess.AddHandler(handler.ReactionAddHandle)
	sess.AddHandler(handler.MessageHandle)
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatalf("Failed to open a discord session: %v", err)
	}
	defer sess.Close()

	fmt.Println("Bot is online!")
	createOSInterruptSingal()
}
