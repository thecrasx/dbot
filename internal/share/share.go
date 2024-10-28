package share

import (
	"dbot/internal/globals"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Session       *discordgo.Session
	MessageCreate *discordgo.MessageCreate
}

func (d *Discord) ChannelMessageReturn(text string) *discordgo.Message {
	if globals.ENABLE_MESSAGE {
		msg, err := d.Session.ChannelMessageSend(d.MessageCreate.ChannelID, text)
		if err != nil {
			log.Printf("Error sending embed: %v", err)
		}
		return msg
	}
	return &discordgo.Message{}
}

func (d *Discord) ChannelMessageEmbedReturn(embed *discordgo.MessageEmbed) *discordgo.Message {
	if globals.ENABLE_MESSAGE {
		msg, err := d.Session.ChannelMessageSendEmbed(d.MessageCreate.ChannelID, embed)
		if err != nil {
			log.Printf("Error sending embed: %v", err)
		}
		return msg
	}
	return &discordgo.Message{}
}
