package games

import (
	"dbot/internal/games/ttc"
	"dbot/internal/share"
	"log"

	"github.com/bwmarrin/discordgo"
)

type unicodeEmoji struct {
	One   string
	Two   string
	Three string
	Four  string
	Five  string
	Six   string
	Seven string
	Eight string
	Nine  string
}

func (u *unicodeEmoji) IsAny(uniEmoji string) (string, int) {
	switch uniEmoji {
	case UnicodeEmoji.One:
		return UnicodeEmoji.One, 1
	case UnicodeEmoji.Two:
		return UnicodeEmoji.Two, 2
	case UnicodeEmoji.Three:
		return UnicodeEmoji.Three, 3
	case UnicodeEmoji.Four:
		return UnicodeEmoji.Four, 4
	case UnicodeEmoji.Five:
		return UnicodeEmoji.Five, 5
	case UnicodeEmoji.Six:
		return UnicodeEmoji.Six, 6
	case UnicodeEmoji.Seven:
		return UnicodeEmoji.Seven, 7
	case UnicodeEmoji.Eight:
		return UnicodeEmoji.Eight, 8
	case UnicodeEmoji.Nine:
		return UnicodeEmoji.Nine, 9
	default:
		return "", -1
	}

}

func (u *unicodeEmoji) AtIndex(index int) string {
	switch index {
	case 0:
		return UnicodeEmoji.One
	case 1:
		return UnicodeEmoji.Two
	case 2:
		return UnicodeEmoji.Three
	case 3:
		return UnicodeEmoji.Four
	case 4:
		return UnicodeEmoji.Five
	case 5:
		return UnicodeEmoji.Six
	case 6:
		return UnicodeEmoji.Seven
	case 7:
		return UnicodeEmoji.Eight
	case 8:
		return UnicodeEmoji.Nine
	default:
		return ""
	}
}

var UnicodeEmoji unicodeEmoji = unicodeEmoji{
	One:   "\u0031\uFE0F\u20E3",
	Two:   "\u0032\uFE0F\u20E3",
	Three: "\u0033\uFE0F\u20E3",
	Four:  "\u0034\uFE0F\u20E3",
	Five:  "\u0035\uFE0F\u20E3",
	Six:   "\u0036\uFE0F\u20E3",
	Seven: "\u0037\uFE0F\u20E3",
	Eight: "\u0038\uFE0F\u20E3",
	Nine:  "\u0039\uFE0F\u20E3",
}

type TTCGame struct {
	Data                 TTCData
	currentPlayer        string // UserID
	reactionsInitialised bool
	Game                 ttc.TicTacToe
}

type TTCData struct {
	UserID_1  string
	UserID_2  string
	ChannelID string
	MessageID string
}

func NewGameTCC(data TTCData, d *share.Discord) *TTCGame {
	game := &TTCGame{
		Data:                 data,
		currentPlayer:        data.UserID_1,
		reactionsInitialised: false,
		Game:                 ttc.New(),
	}

	member, err := d.Session.GuildMember(d.MessageCreate.GuildID, game.currentPlayer)
	if err != nil {
		log.Printf("[game.NewGameTTC] Failed to get Member")
	}
	msg := d.ChannelMessageEmbedReturn(
		&discordgo.MessageEmbed{
			Title:       game.Game.TableToString(),
			Description: member.DisplayName() + " turn",
		},
	)
	game.Data.MessageID = msg.ID
	return game
}

func (g *TTCGame) SwitchPlayer() {
	if g.currentPlayer == g.Data.UserID_1 {
		g.currentPlayer = g.Data.UserID_2
	} else {
		g.currentPlayer = g.Data.UserID_1
	}
}

func (g *TTCGame) SetCurrentPlayer(userID string) {
	if userID != g.Data.UserID_1 || userID != g.Data.UserID_2 {
		panic("This UserID doesn't exist")
	}
	g.currentPlayer = userID
}

func (g *TTCGame) CurrentPlayer() string {
	return g.currentPlayer
}

func (g *TTCGame) InitReactions(s *discordgo.Session) {
	if g.reactionsInitialised || g.Data.MessageID == "" {
		return
	}

	s.MessageReactionAdd(g.Data.ChannelID, g.Data.MessageID, UnicodeEmoji.One)
	s.MessageReactionAdd(g.Data.ChannelID, g.Data.MessageID, UnicodeEmoji.Two)
	s.MessageReactionAdd(g.Data.ChannelID, g.Data.MessageID, UnicodeEmoji.Three)
	s.MessageReactionAdd(g.Data.ChannelID, g.Data.MessageID, UnicodeEmoji.Four)
	s.MessageReactionAdd(g.Data.ChannelID, g.Data.MessageID, UnicodeEmoji.Five)
	s.MessageReactionAdd(g.Data.ChannelID, g.Data.MessageID, UnicodeEmoji.Six)
	s.MessageReactionAdd(g.Data.ChannelID, g.Data.MessageID, UnicodeEmoji.Seven)
	s.MessageReactionAdd(g.Data.ChannelID, g.Data.MessageID, UnicodeEmoji.Eight)
	s.MessageReactionAdd(g.Data.ChannelID, g.Data.MessageID, UnicodeEmoji.Nine)

	g.reactionsInitialised = true
}

func (g *TTCGame) HandleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	_, pos := UnicodeEmoji.IsAny(r.Emoji.Name)

	if pos == -1 {
		s.MessageReactionRemove(g.Data.ChannelID, g.Data.MessageID, r.Emoji.Name, r.UserID)
		return
	}

	if r.UserID != g.currentPlayer {
		s.MessageReactionRemove(g.Data.ChannelID, g.Data.MessageID, r.Emoji.Name, r.UserID)
		return
	}

	go s.MessageReactionRemove(g.Data.ChannelID, g.Data.MessageID, r.Emoji.Name, r.UserID)
	go s.MessageReactionRemove(g.Data.ChannelID, g.Data.MessageID, r.Emoji.Name, s.State.User.ID)

	g.Game.SetSignAtPosition(pos - 1)
	if g.Game.CheckWinner() {
		member, _ := s.GuildMember(r.GuildID, g.currentPlayer)
		g.WinGame(s, member.DisplayName())
		return
	}

	if g.Game.AvailableTurns() < 1 {
		g.TieGame(s)
		return
	}

	g.SwitchPlayer()
	if g.currentPlayer == s.State.User.ID {
		s.ChannelMessageEditEmbed(g.Data.ChannelID, g.Data.MessageID,
			&discordgo.MessageEmbed{
				Title:       g.Game.TableToString(),
				Description: "dbot turn",
			},
		)
		pos, _ := g.Game.AutoSet()

		if g.Game.CheckWinner() {
			g.WinGame(s, "dbot")
			return
		} else if g.Game.AvailableTurns() < 1 {
			g.TieGame(s)
			return
		}
		g.SwitchPlayer()
		go s.MessageReactionRemove(g.Data.ChannelID, g.Data.MessageID, UnicodeEmoji.AtIndex(pos), s.State.User.ID)
	}

	member, err := s.GuildMember(r.GuildID, g.currentPlayer)
	if err != nil {
		log.Printf("[game.HandleReactionAdd] Failed to get Member")
	}

	go s.ChannelMessageEditEmbed(g.Data.ChannelID, g.Data.MessageID,
		&discordgo.MessageEmbed{
			Title:       g.Game.TableToString(),
			Description: member.DisplayName() + " turn",
		},
	)
}

func (g *TTCGame) RemoveAllReactions(s *discordgo.Session) {
	msg, err := s.ChannelMessage(g.Data.ChannelID, g.Data.MessageID)
	if err != nil {
		log.Printf("RemoveAllReactions failed: %v", err)
	}

	for _, react := range msg.Reactions {
		s.MessageReactionRemove(
			g.Data.ChannelID,
			g.Data.MessageID,
			react.Emoji.Name,
			s.State.User.ID)
	}
}

func (g *TTCGame) WinGame(s *discordgo.Session, name string) {
	s.ChannelMessageEditEmbed(g.Data.ChannelID, g.Data.MessageID,
		&discordgo.MessageEmbed{
			Title:       name + " Won",
			Description: g.Game.TableToString(),
		},
	)
	g.FinishGame(s)
}

func (g *TTCGame) TieGame(s *discordgo.Session) {
	s.ChannelMessageEditEmbed(g.Data.ChannelID, g.Data.MessageID,
		&discordgo.MessageEmbed{
			Title:       "Tie",
			Description: g.Game.TableToString(),
		},
	)
	g.FinishGame(s)
}

func (g *TTCGame) FinishGame(s *discordgo.Session) {
	s.MessageReactionAdd(g.Data.ChannelID, g.Data.MessageID, "\u2716\uFE0F")
	g.RemoveAllReactions(s)
}
