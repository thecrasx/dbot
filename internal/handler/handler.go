package handler

import (
	"dbot/internal/globals"
	"dbot/internal/promt"
	"dbot/internal/share"
	"dbot/internal/state"
	"log"

	"github.com/bwmarrin/discordgo"
)

func MessageHandle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.ChannelID != globals.DBOT_CHANNEL_ID {
		return
	}

	d := &share.Discord{
		Session:       s,
		MessageCreate: m,
	}

	if promt.IsPrompt(m.Content) {
		promt.Promt(m.Content, d)
	}
}

func ReactionAddHandle(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.ChannelID != globals.DBOT_CHANNEL_ID {
		return
	}
	if s.State.User.ID == r.Member.User.ID && r.Emoji.Name != "\u2716\uFE0F" {
		return
	}

	// fmt.Println(utils.EmojiToUnicodeString(r.Emoji.Name))
	gs := state.Global.GetGameStateByMessageID(r.MessageID)
	if gs == nil {
		log.Println("[handler.ReactionAddHandler] GameState not found")
		return
	}

	if r.Emoji.Name == "\u2716\uFE0F" {
		s.MessageReactionRemove(gs.Game.Data.ChannelID, gs.Game.Data.MessageID, "\u2716\uFE0F", s.State.User.ID)
		state.Global.RemoveGameState(gs)
	} else {
		gs.Game.HandleReactionAdd(s, r)
	}
}

