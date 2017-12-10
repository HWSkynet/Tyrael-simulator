package main

import (
	"fmt"
	//"io/ioutil"
	"os"
	"os/signal"
	//"strings"
	"syscall"
	//"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

var debug_channel string
var talking_channel string

func main() {
	viper.SetDefault("token", 0)
	viper.SetDefault("debugChannel", 0)
	viper.SetDefault("talkingChannel", 0)
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	token := viper.Get("token").(string)
	fmt.Print("token=" + token + "\r\n")

	debug_channel = viper.Get("debugChannel").(string)
	fmt.Print("debugChannel=" + debug_channel + "\r\n")

	talking_channel = viper.Get("talkingChannel").(string)
	fmt.Print("talkingChannel=" + talking_channel + "\r\n")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(ready)
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("群主上线.")
	dg.ChannelMessageSend(talking_channel, "萌七！")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("群主下线.")
	dg.ChannelMessageSend(debug_channel, "群主下线.")
	// Cleanly close down the Discord session.
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.ChannelMessageSend(debug_channel, "群主上线.")
	s.UpdateStatus(0, "Artifact Idiot")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "苟利国家生死以" {
		s.ChannelMessageSend(m.ChannelID, "岂因祸福避趋之")
	}
}
