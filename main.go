package main

import (
	"fmt"
	//"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

const version string = "Alpha build 12204"

var debug_channel string
var talking_channel string

type qunzhu struct {
	version  string
	freeze   bool
	boring   int
	sleeping int
	silence  int
	shy      int
	energy   int // 元气
}

var tyrael qunzhu = qunzhu{
	freeze:   false,
	boring:   0,
	sleeping: 0,
	silence:  0,
	shy:      0,
	energy:   1000,
}

var GSession *discordgo.Session

func main() {
	rand.Seed(time.Now().UnixNano())
	viper.SetDefault("token", 0)
	viper.SetDefault("debugChannel", 0)
	viper.SetDefault("talkingChannel", 0)
	viper.SetDefault("oldversion", "unknown")
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

	oldVersion := viper.Get("oldversion").(string)
	tyrael.version = version

	var newVersion bool
	if oldVersion == tyrael.version {
		newVersion = false
	} else {
		newVersion = true
	}

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
	GSession = dg

	fmt.Println("群主上线.")
	if newVersion {
		dg.ChannelMessageSend(debug_channel, "升级成功！\n旧版本:"+oldVersion+"\n当前版本:"+tyrael.version)
		viper.Set("oldversion", tyrael.version)
		viper.WriteConfig()
	}
	viper.Reset()

	tyrael.initEnergy()
	dg.ChannelMessageSend(debug_channel, fmt.Sprintf("当前元气值:%d", tyrael.energy))

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
	GSession.UpdateStatus(0, gameName[rand.Intn(len(gameName))])
}

func (*qunzhu) talk(channel string, str string, speed int) {
	go func(channel string, str string, speed int) {
		tyrael.boring /= 4
		<-time.After(time.Millisecond * time.Duration(speed) * time.Duration(len(str)))
		GSession.ChannelMessageSend(channel, str)
	}(channel, str, speed)
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.ChannelMessageSend(debug_channel, startMessage[rand.Intn(len(startMessage))])
	s.UpdateStatus(0, gameName[rand.Intn(len(gameName))])
	go clock(clockInput)
}

func getMinute() int {
	h, m, _ := time.Now().Clock()
	return h*60 + m
}

func (self *qunzhu) initEnergy() {
	m := getMinute()
	if m < 420 || m > 1380 {
		self.energy = rand.Intn(60)
	} else {
		self.energy = (700 + rand.Intn(150)) * (m - 420) / 1000
	}
}

func clock(input chan interface{}) {
	min := time.NewTicker(1 * time.Minute)
	halfhour := time.NewTicker(23 * time.Minute)
	var lastBoring *discordgo.Message
	for {
		select {
		case <-min.C:
			if getMinute() == 0 {
				GSession.ChannelMessageSend(debug_channel, fmt.Sprintf("零时迷子即将启动\n当前元气值：%d", tyrael.energy))
				tyrael.energy /= 4 // 零时迷子
				tyrael.talk(debug_channel, "零时迷子启动", 300)
				tyrael.talk(debug_channel, fmt.Sprintf("当前元气值：%d", tyrael.energy), 600)
			}
			if tyrael.sleeping == 0 {
				tyrael.energy -= 1
				tyrael.silence += 1
				tyrael.boring += 1
				if rand.Intn(60) > tyrael.energy { // 累得睡着
					tyrael.sleeping += 1
					tyrael.boring = 0
					GSession.UpdateStatus(0, "睡着了Z.z.z.")
				} else if tyrael.boring > 40+rand.Intn(100) { // 无聊得打瞌睡
					tyrael.sleeping += 1
					GSession.UpdateStatus(0, "打瞌睡Z.z.z.")
				} else if rand.Intn(100) < 5 {
					if !tyrael.freeze {
						if tyrael.shy > 0 {
							GSession.ChannelMessageDelete(talking_channel, lastBoring.ID)
							tyrael.shy = 0
						}
						lastBoring, _ = GSession.ChannelMessageSend(talking_channel, IdleTalk())
						tyrael.shy += 1
					}
				}
			} else {
				tyrael.sleeping += 1
				tyrael.energy += 1
				if tyrael.boring < 30 {
					tyrael.energy += 1
					if tyrael.energy > 800+rand.Intn(160) {
						if tyrael.sleeping/60 < 7 {
							GSession.ChannelMessageSend(talking_channel, fmt.Sprintf("哈欠，虽然只睡了%d个多小时，但是感觉元气满满的呢", tyrael.sleeping/60))
						} else if tyrael.sleeping/60 > 8 {
							GSession.ChannelMessageSend(talking_channel, fmt.Sprintf("哇，不小心睡了%d个多小时，糟了糟了", tyrael.sleeping/60))
						} else {
							GSession.ChannelMessageSend(talking_channel, "<:xyx:389356458539614208><:xyx:389356458539614208><:xyx:389356458539614208>")
						}
						tyrael.newStatus()
						tyrael.sleeping = 0
						tyrael.silence = 0
						tyrael.boring = 0
					}
				} else {
					tyrael.energy -= 1
					if rand.Intn(1000) < tyrael.boring {
						tyrael.sleeping = 0
						tyrael.silence = 0
						tyrael.boring = 0
						lastBoring, _ = GSession.ChannelMessageSend(talking_channel, "啊，好无聊啊")
						tyrael.shy += 1
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
	channelName, _ := s.Channel(m.ChannelID).Name
	fmt.Printf("[" + channelName + "]" + m.Author.Username + ":" + m.Content + "\n")
	for _, v := range m.Embeds {
		fmt.Printf(v.Type)
	}
	for _, v := range m.Attachments {
		fmt.Printf("图片尺寸:%dx%d\n", v.Width, v.Height)
	}

	if !m.Author.Bot && m.ChannelID == talking_channel {
		tyrael.silence = 0
		tyrael.shy = 0
		if tyrael.sleeping > 0 {
			if tyrael.boring > 30 {
				if rand.Intn(100) < 20 {
					tyrael.sleeping = 0
					tyrael.boring /= 2
					tyrael.talk(m.ChannelID, "<:xyx:389356458539614208>", 300)
					tyrael.newStatus()
				}
			} else {
				if rand.Intn(100) < 9 {
					tyrael.sleeping = 0
					tyrael.talk(m.ChannelID, forceWakeup[rand.Intn(len(forceWakeup))], 300)
					tyrael.newStatus()
				}
			}
		} else {
			if tyrael.boring > 0 {
				tyrael.boring -= 1
			}
			if strings.Contains(m.Content, "来吃鸡吧") {
				if GameState == "idle" {
					GameNewRoom()
					tyrael.talk(m.ChannelID, "吃鸡房间建立完成，现在可以加入战局", 100)
				} else {
					tyrael.talk(m.ChannelID, "请等待上一只鸡吃完", 300)
				}
			}
			if strings.Contains(m.Content, "强制关闭战局") {
				if GameState == "idle" {
					tyrael.talk(m.ChannelID, "空即是色，施主怕是杂念太多", 100)
				} else {
					GameClear()
					tyrael.talk(m.ChannelID, "世界，又归于和平...", 100)
				}
			}
		}
	}
	if !tyrael.freeze && tyrael.sleeping == 0 && m.ChannelID == talking_channel {
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
			s.ChannelMessageSend(debug_channel, "呜~呜嗯~ 嗯~~ 啊~ 憋死我了")
		}
	}

	// 愿此bot寿与天齐
	if m.Content == "苟利国家生死以" {
		s.ChannelMessageSend(m.ChannelID, "岂因祸福避趋之")
	}

	if GameState != "idle" && m.ChannelID == GameChannel.ID {
		GameRoomMessageHandler(s, m)
	}
}
