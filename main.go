package main

import (
	"fmt"
	//"io/ioutil"
	"os"
	"os/signal"
	//"strings"
	"syscall"
	//"time"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

var debug_channel string
var talking_channel string

var freeze bool = false

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

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("群主上线.")
	dg.ChannelMessageSend(talking_channel, "<:xyx:389356458539614208>")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("群主下线.")
	dg.ChannelMessageSend(debug_channel, "群主下线.")
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.ChannelMessageSend(debug_channel, "女装已经换好，请各位来撩")
	s.UpdateStatus(0, "女装山脉IV")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf(m.Content + "\n")
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !freeze && m.ChannelID == talking_channel {
		// 图片
		// m.Type
		// 特定人识别
		if IsVip(m.Author.ID) {
			// 关键词识别

			// 普通随机回复
			if rand.Intn(1000) > 900 {
				s.ChannelMessageSend(m.ChannelID, GetRandom(m.Author.ID))
			}
		} else {
			if rand.Intn(1000) > 950 {
				s.ChannelMessageSend(m.ChannelID, "<:xyx:389356458539614208>")
			}
		}
	}

	if m.Author.ID == "377366407089881088" {
		if !freeze && m.Content == "一二三稻草人" {
			freeze = true
			s.ChannelMessageSend(debug_channel, "唔，呜呜唔，唔~~~")
		}
		if freeze && m.Content == "让他说话" {
			freeze = false
			s.ChannelMessageSend(debug_channel, "呜~~~啊~~~憋死我了")
		}
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "苟利国家生死以" {
		s.ChannelMessageSend(m.ChannelID, "岂因祸福避趋之")
	}
}
