package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DevDisc/thotbot/src/lib/app"

	"github.com/bwmarrin/discordgo"
)

func main() {
	application := app.Cli(&app.CliMethods{
		RunApp: runApp,
	})

	err := application.Run(os.Args)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func runApp(cfg *app.Config) error {

	// Create Session
	session, err := discordgo.New("Bot " + cfg.DiscordAuthToken)
	if err != nil {
		return err
	}

	// Add handlers
	session.AddHandler(handlePing)
	session.AddHandler(handleFutures)

	// Only care about inputs
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = session.Open()
	if err != nil {
		return err
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()
	return nil
}

func handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is ping, replt with pong
	if m.Content == "ping" {
		timeStamp, _ := m.Timestamp.Parse()
		log.Println(m.Author.ID + " pinged at " + timeStamp.String())
		s.ChannelMessageSend(m.ChannelID, "pong")
	}

	// Reverse
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "ping")
	}
}

func handleFutures(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is ping, replt with pong
	if m.Content == "!futures" || m.Content == "!market" {
		timeStamp, _ := m.Timestamp.Parse()
		if IsMarketClosed(timeStamp) {
			msg := "Futures are closed "
			msg += m.Author.ID
			s.ChannelMessageSend(m.ChannelID, msg)
		}
	}
}

func IsMarketClosed(t time.Time) bool {
	hour, _, _ := t.Clock()
	// Is it between 4 and 5 CT?
	fmt.Println(hour)
	if hour >= 22 && hour < 23 {
		return true
	}
	// Is it between Friday after 4 and Sunday before 5 CT
	if t.Weekday().String() == "Saturday" {
		return true
	}
	if t.Weekday().String() == "Friday" {
		if hour >= 22 {
			return true
		}
	}
	if t.Weekday().String() == "Sunday" {
		if hour > 0 && hour < 23 {
			return true
		}
	}
	return false
}
