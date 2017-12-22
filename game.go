// game.go
// 吃鸡相关逻辑和接口
package main

import (
	//"math/rand"

	"github.com/HWSkynet/cpgame"
	"github.com/bwmarrin/discordgo"
)

var GameState string = "idle"
var GameChannel *discordgo.Channel
var Players cpgame.Player

func GameNewRoom() {
	var err error
	GameChannel, err = GSession.GuildChannelCreate("377366788322623491", "GAMEROOM-eat-chicken", "text")
	if err != nil {
		panic(err)
	}
	GameState = "ready"
}

func GameClear() {
	GSession.ChannelDelete(GameChannel.ID)
	GameState = "idle"
}

func GameRoomMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

}
