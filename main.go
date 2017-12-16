package main

import (
	"fmt"
	//"io/ioutil"
	"os"
	"os/signal"
	//"strings"
	"math/rand"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

var debug_channel string
var talking_channel string

var freeze bool = false

func main() {
	rand.Seed(time.Now().UnixNano())
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
	s.ChannelMessageSend(debug_channel, "今天的女装已经准备好了，请各位赶快领取吧")
	s.UpdateStatus(0, "女装山脉IV")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Printf(m.Author.Username + ":" + m.Content + "\n")
	for _, v := range m.Embeds {
		fmt.Printf(v.Type)
	}
	for _, v := range m.Attachments {
		fmt.Printf("图片尺寸:%dx%d\n", v.Width, v.Height)
	}

	if !freeze && m.ChannelID == talking_channel {
		// 图片
		if len(m.Attachments) > 0 && m.Attachments[0].Width > 0 {
			if rand.Intn(100) < 15 {
				go func() {
					str := PicTalk()
					<-time.After(time.Millisecond * 500 * time.Duration(len(str)))
					s.ChannelMessageSend(m.ChannelID, str)
				}()
			}
		}
		// m.Type
		// 特定人识别
		if len(m.Content) > 0 && IsVip(m.Author.ID) {
			rands := rand.Intn(100)
			fmt.Printf("rands=%d\n", rands)
			if rands < 10 {
				go func() {
					str := Talk(m.Author.ID, m.Content)
					<-time.After(time.Millisecond * 300 * time.Duration(len(str)))
					s.ChannelMessageSend(m.ChannelID, str)
				}()
			}
		} else {
			if rand.Intn(100) < 5 {
				s.ChannelMessageSend(m.ChannelID, IdleTalk())
			}
		}
	}

	// 临时禁言用
	if m.Author.ID == "377366407089881088" {
		if !freeze && m.Content == "一二三木头人" {
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
