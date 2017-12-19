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

type qunzhu struct {
	version  string
	freeze   bool
	boring   int
	sleeping int
	silence  int
}

var tyrael qunzhu = qunzhu{
	freeze:   false,
	boring:   0,
	sleeping: 0,
	silence:  0,
}

var gSession *discordgo.Session

func main() {
	rand.Seed(time.Now().UnixNano())
	viper.SetDefault("token", 0)
	viper.SetDefault("debugChannel", 0)
	viper.SetDefault("talkingChannel", 0)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	token := viper.Get("token").(string)
	fmt.Println("token=" + token)

	debug_channel = viper.Get("debugChannel").(string)
	fmt.Println("debugChannel=" + debug_channel)

	talking_channel = viper.Get("talkingChannel").(string)
	fmt.Println("talkingChannel=" + talking_channel)

	viper.SetConfigName("version")
	viper.SetConfigType("json")
	tyrael.version = viper.Get("version").(string)

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
	gSession = dg

	fmt.Println("群主上线.")

	dg.ChannelMessageSend(debug_channel, "当前版本:"+tyrael.version)
	msg, _ := dg.ChannelMessageSend(talking_channel, "前方高能反应，非战斗人员请迅速撤离")
	go func() {
		<-time.After(time.Second * 5)
		dg.ChannelMessageDelete(msg.ChannelID, msg.ID)
	}()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("群主下线.")
	dg.ChannelMessageSend(debug_channel, "群主下线.")
	dg.Close()
}

var clockInput = make(chan interface{})

// 启动信息
var startMessage = []string{
	"今天的女装已经准备好了，请各位赶快领取吧",
	"今天的女装，啧啧，是女仆装哦，请各位快来领取吧",
	"今天。。哇！有兔耳欸！请各位快来领取吧",
	"今天是普通的水手服呢，请各位赶紧换好",
}

// 随机状态
var gameName = []string{
	"女装山脉IV",
	"女装传说",
	"女装破环神III",
	"女装英雄",
	"女装先锋",
	"女装争霸",
	"女装的远征",
	"女装王座",
	"荒岛女装",
	"坎巴拉女装计划",
	"女装骑士",
	"微软模拟女装",
	"微软女装飞行",
	"女装谷物语",
}

// 被吵醒的回复
var forceWakeup = []string{
	"哇，你们不用睡觉的么",
}

func (*qunzhu) newStatus() {
	gSession.UpdateStatus(0, gameName[rand.Intn(len(gameName))])
}

func (*qunzhu) talk(channel string, str string, speed int) {
	go func(channel string, str string, speed int) {
		tyrael.boring /= 4
		<-time.After(time.Millisecond * time.Duration(speed) * time.Duration(len(str)))
		gSession.ChannelMessageSend(channel, str)
	}(channel, str, speed)
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.ChannelMessageSend(debug_channel, startMessage[rand.Intn(len(startMessage))])
	s.UpdateStatus(0, gameName[rand.Intn(len(gameName))])
	go clock(clockInput)
}

func clock(input chan interface{}) {
	min := time.NewTicker(1 * time.Minute)
	halfhour := time.NewTicker(29 * time.Minute)
	for {
		select {
		case <-min.C:
			if tyrael.sleeping == 0 {
				tyrael.silence += 1
				// 休眠模式
				if tyrael.silence > 54 {
					tyrael.sleeping += 1
					gSession.UpdateStatus(0, "打瞌睡Z.z.z..")
				}
				tyrael.boring += 1
				if rand.Intn(207) < tyrael.boring {
					tyrael.boring /= 4
					if !tyrael.freeze {
						gSession.ChannelMessageSend(talking_channel, IdleTalk())
					}
				}
			} else {
				tyrael.sleeping += 1
				if tyrael.sleeping > 330 {
					if rand.Intn(1000) < tyrael.sleeping-330 {
						gSession.ChannelMessageSend(talking_channel, fmt.Sprintf("哈欠，才睡了%d个多小时，好困", tyrael.sleeping/60))
						tyrael.newStatus()
						tyrael.sleeping = 0
						tyrael.silence = 0
						tyrael.boring = 0
					}
				}
			}
		case <-halfhour.C:
			if tyrael.sleeping == 0 {
				tyrael.newStatus()
			}
		}
	}
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

	if !m.Author.Bot && m.ChannelID == talking_channel {
		tyrael.silence = 0
		if tyrael.sleeping > 0 {
			if rand.Intn(100) < 10 {
				tyrael.sleeping = 0
				tyrael.talk(m.ChannelID, forceWakeup[rand.Intn(len(forceWakeup))], 300)
				tyrael.newStatus()
			}
		}
	}

	if !tyrael.freeze && m.ChannelID == talking_channel {
		// 图片
		if len(m.Attachments) > 0 && m.Attachments[0].Width > 0 {
			if rand.Intn(100) < 15 {
				tyrael.talk(m.ChannelID, PicTalk(), 500)
			}
		}
		// 特定人识别
		if len(m.Content) > 0 && IsVip(m.Author.ID) {
			rands := rand.Intn(100)
			words := Talk(m.Author.ID, m.Content, rands)
			if len(words) > 0 {
				tyrael.talk(m.ChannelID, words, 300)
			}
		}
	}

	// 临时禁言用
	if m.Author.ID == "377366407089881088" {
		if !tyrael.freeze && m.Content == "一二三木头人" {
			tyrael.freeze = true
			s.ChannelMessageSend(debug_channel, "唔，呜呜唔，唔~~~")
		}
		if tyrael.freeze && m.Content == "让他说话" {
			tyrael.freeze = false
			s.ChannelMessageSend(debug_channel, "呜~~ ~啊~~ ~憋死我了")
		}
	}

	// 愿此bot寿与天齐
	if m.Content == "苟利国家生死以" {
		s.ChannelMessageSend(m.ChannelID, "岂因祸福避趋之")
	}
}
