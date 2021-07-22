package app

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Server struct {
	*Config
	router *discordgo.Session
}

func NewServer(cfg *Config) (*Server, error) {

	session, err := discordgo.New("Bot " + cfg.DiscordAuthToken)

	if err != nil {
		return nil, err
	}

	server := &Server{
		Config: cfg,
		router: session,
	}

	// Add Handlers
	server.router.AddHandler(server.HandlePing)
	server.router.AddHandler(server.HandleFutures)
	server.router.AddHandler(server.HandlePort)

	// Only care about inputs
	server.router.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	return server, nil
}

// Start the server
func (s *Server) Start() error {
	return s.router.Open()
}

func (s *Server) Stop() error {
	return s.router.Close()
}

// Handlers

func (s *Server) HandlePort(sess *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.ID == s.router.State.User.ID {
		return
	}

	// If the message starts with "!port" then handle
	if !strings.HasPrefix(m.Content, "!port") {
		return
	}

	// Handle desired function
	route := strings.Split(m.Content, " ")

	// Need a command
	if len(route) == 1 {
		return
	}

	// Handle help and return
	if route[1] == "help" {
		s.RunHelp(m.ChannelID)
	}
}

func (s *Server) RunHelp(channelID string) {
	// Create help message fields
	fields := make([]*discordgo.MessageEmbedField, 6)
	fields[0] = &discordgo.MessageEmbedField{
		Name:   "!port help",
		Value:  "Use to show help function",
		Inline: true,
	}
	fields[1] = &discordgo.MessageEmbedField{
		Name:   "!port show",
		Value:  "Display your current port and daily change",
		Inline: true,
	}
	fields[2] = &discordgo.MessageEmbedField{
		Name:   "!port chart",
		Value:  "Display 5 min chart for every symbol in your port",
		Inline: true,
	}
	fields[3] = &discordgo.MessageEmbedField{
		Name:   "!port quote",
		Value:  "Display quotes for every symbol in your port",
		Inline: true,
	}
	fields[4] = &discordgo.MessageEmbedField{
		Name:  "!port add QUANTITY SYMBOL",
		Value: "Add a quantity of stock to your port (!port add 10 AAPL)",
	}
	fields[5] = &discordgo.MessageEmbedField{
		Name:  "!port remove QUANTITY SYMBOL",
		Value: "Add a quantity of stock to your port (!port remove 10 AAPL). Will not go below 0",
	}

	// Create help message
	helpMessage := discordgo.MessageEmbed{
		Title:       "ThotBot Port App Help",
		Description: "Quick intro on how to use the port app",
		Color:       10177720,
		Fields:      fields,
	}
	s.router.ChannelMessageSendEmbed(channelID, &helpMessage)
}

func (s *Server) HandlePing(sess *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.ID == s.router.State.User.ID {
		return
	}

	// If the message is ping, replt with pong
	if m.Content == "ping" {
		timeStamp, _ := m.Timestamp.Parse()
		log.Println(m.Author.ID + " pinged at " + timeStamp.String())
		s.router.ChannelMessageSend(m.ChannelID, "pong")
	}

	// Reverse
	if m.Content == "pong" {
		s.router.ChannelMessageSend(m.ChannelID, "ping")
	}
}

func (s *Server) HandleFutures(sess *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.ID == s.router.State.User.ID {
		return
	}

	// If the message is ping, replt with pong
	if m.Content == "!futures" || m.Content == "!market" {
		timeStamp, _ := m.Timestamp.Parse()
		if IsMarketClosed(timeStamp) {
			msg := "Futures are closed "
			msg += m.Author.Username
			s.router.ChannelMessageSend(m.ChannelID, msg)
		}
	}
	if m.Content == "!f" {
		msg := "?quote /ES /NQ /RTY /YM /GC /SI /CL"
		s.router.ChannelMessageSend(m.ChannelID, msg)
	}

	if m.Content == "!devport" {
		msg := "?c2 nio pltr open vldr coin"
		s.router.ChannelMessageSend(m.ChannelID, msg)
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
