package app

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Server struct {
	*Config
	router *discordgo.Session
}

func NewServer(cfg *Config) (*Server, error) {

	// Ensure port dir exists
	_, err := os.Stat(cfg.PortPath)
	if os.IsNotExist(err) {
		return nil, err
	}

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

	errorMsg := "Invalid syntax, look at `!port help` for help"

	// Need a command
	if len(route) == 1 {
		return
	}

	// Handle help and return
	if route[1] == "help" {
		s.RunHelp(m.ChannelID)
		return
	}

	// Handle show
	if route[1] == "show" {
		s.RunShow(m)
		return
	}

	// Handle chart
	if route[1] == "chart" {
		s.RunChart(m)
		return
	}

	// Handle quote
	if route[1] == "quote" {
		s.RunQuote(m)
		return
	}

	// Handle add
	if route[1] == "add" {
		if len(route) == 4 {
			quantity, err := strconv.Atoi(route[2])
			if err == nil {
				s.RunAdd(m, quantity, route[3])
				return
			} else {
				errorMsg = "Unable to read quantity of: " + route[2]
			}
		}
	}

	// Handle remove
	if route[1] == "remove" {
		if len(route) == 4 {
			quantity, err := strconv.Atoi(route[2])
			if err == nil {
				s.RunRemove(m, quantity, route[3])
				return
			} else {
				errorMsg = "Unable to read quantity of: " + route[2]
			}
		}
	}

	// Otherwise invalid
	s.router.ChannelMessageSend(m.ChannelID, errorMsg)
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
		Value:  "Display your current port and daily change (UNIMPLEMENTED)",
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
		Description: "Quick intro on how to use the port app. To create a port all you need to do is run `!port add QUANITIY SYMBOL` and ThotBot will create a port for you.",
		Color:       10177720,
		Fields:      fields,
	}
	s.router.ChannelMessageSendEmbed(channelID, &helpMessage)
}

func (s *Server) RunShow(m *discordgo.MessageCreate) {

	if !s.UserExists(m) {
		return
	}

	s.router.ChannelMessageSend(m.ChannelID, "Unimplemented feature, will implement soon....")
}

func (s *Server) RunChart(m *discordgo.MessageCreate) {
	s.RunStella("?c2 ", m)
}

func (s *Server) RunQuote(m *discordgo.MessageCreate) {
	s.RunStella("?quote ", m)
}

func (s *Server) RunStella(command string, m *discordgo.MessageCreate) {
	if !s.UserExists(m) {
		return
	}

	// Load existing port
	userPortPath := s.GetUserPortPath(m)
	port, err := LoadPort(userPortPath)
	if err != nil {
		s.router.ChannelMessageSend(m.ChannelID, "Unable to load your port (@Dev, wtf?)")
		return
	}

	// Add symbols
	for symbol, _ := range port.Holdings {
		command = command + symbol + " "
	}
	s.router.ChannelMessageSend(m.ChannelID, command)
}

func (s *Server) RunAdd(m *discordgo.MessageCreate, quantity int, symbol string) {
	// If not exists then a port will be made
	port := new(Port)
	userPortPath := s.GetUserPortPath(m)
	if _, err := os.Stat(userPortPath); err == nil {
		// Load existing port
		port, err = LoadPort(userPortPath)
		if err != nil {
			s.router.ChannelMessageSend(m.ChannelID, "Unable to load your port (@Dev, wtf?)")
			return
		}

	}

	if port.Holdings == nil {
		port.Holdings = make(map[string]int)
		port.Holdings[symbol] = quantity
		s.router.ChannelMessageSend(m.ChannelID, "Started new holding in "+symbol+" of "+strconv.Itoa(quantity)+" shares")
	} else {
		// Check holds to see if symbol exists and add
		if _, ok := port.Holdings[symbol]; ok {
			oldQuantity := port.Holdings[symbol]
			newQuantity := oldQuantity + quantity
			port.Holdings[symbol] = newQuantity

			s.router.ChannelMessageSend(m.ChannelID, "Increased holdings in "+symbol+" from "+strconv.Itoa(oldQuantity)+" to "+strconv.Itoa(newQuantity))
		} else {
			port.Holdings[symbol] = quantity
			s.router.ChannelMessageSend(m.ChannelID, "Started new holding in "+symbol+" of "+strconv.Itoa(quantity)+" shares")
		}
	}
	// Save changes
	SavePort(port, userPortPath)
}

func (s *Server) RunRemove(m *discordgo.MessageCreate, quantity int, symbol string) {

	if !s.UserExists(m) {
		return
	}

	// Load existing port
	userPortPath := s.GetUserPortPath(m)
	port, err := LoadPort(userPortPath)

	if err != nil {
		s.router.ChannelMessageSend(m.ChannelID, "Unable to load your port (@Dev, wtf?)")
		return
	}

	// Check holds to see if symbol exists and subtract
	if _, ok := port.Holdings[symbol]; ok {
		oldQuantity := port.Holdings[symbol]
		newQuantity := oldQuantity - quantity
		if newQuantity < 0 {
			newQuantity = 0
		}
		port.Holdings[symbol] = newQuantity
		s.router.ChannelMessageSend(m.ChannelID, "Decreased holdings in "+symbol+" from "+strconv.Itoa(oldQuantity)+" to "+strconv.Itoa(newQuantity))

		// If 0 remove from map
		if newQuantity == 0 {
			delete(port.Holdings, symbol)
		}
	} else {
		s.router.ChannelMessageSend(m.ChannelID, "Unable to find "+symbol+" in port")
	}

	// Save changes or delete file if map is empty
	if len(port.Holdings) == 0 {
		err := os.Remove(userPortPath)
		if err != nil {
			s.router.ChannelMessageSend(m.ChannelID, "Someone fucked up ThotBot @Dev")
		}
	} else {
		SavePort(port, userPortPath)
	}
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

func (s *Server) GetUserPortPath(m *discordgo.MessageCreate) string {
	return s.PortPath + m.Author.ID + ".json"
}

func (s *Server) UserExists(m *discordgo.MessageCreate) bool {
	userPortPath := s.GetUserPortPath(m)
	_, err := os.Stat(userPortPath)
	if os.IsNotExist(err) {
		s.router.ChannelMessageSend(m.ChannelID, "You don't have anything in your port, take a look at `!port help` to see how to add symbols")
		return false
	}
	return true
}
