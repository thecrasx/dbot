package promt

import (
	"dbot/internal/errors"
	"dbot/internal/games"
	"dbot/internal/share"
	"dbot/internal/state"
	"fmt"
	"log"
)

func IsPrompt(text string) bool {
	if text[0] == '!' {
		return true
	} else {
		return false
	}
}

func Promt(text string, d *share.Discord) {
	promt, content := getPromtAndContent(text)

	switch promt {
	case PROMT_TTC:
		um, err := getUserMention(content)
		if err != nil {
			d.ChannelMessageReturn("You need to challenge someone! [use !ttc @user]")
			break
		}

		userID_1 := d.MessageCreate.Author.ID
		userID_2 := getUserIDFromMention(um)
		if userID_1 == userID_2 {
			d.ChannelMessageReturn("You cannot challenge yourself " + um)
			break
		}
		member, err := d.Session.GuildMember(d.MessageCreate.GuildID, userID_2)
		if err != nil {
			log.Printf("Failed to get Guild member: %v", err)
			d.ChannelMessageReturn("User " + um + " cannot be found")
			break
		}

		gs1 := state.Global.GetGameStateByUserID(userID_1)
		gs2 := state.Global.GetGameStateByUserID(userID_2)

		if gs1 != nil {
			d.ChannelMessageReturn("You have ongoing game!")
			break
		} else if gs2 != nil {
			d.ChannelMessageReturn(member.DisplayName() + " has already ongoing game")
			break
		} else {
			ttcData := games.TTCData{
				UserID_1:  userID_1,
				UserID_2:  userID_2,
				ChannelID: d.MessageCreate.ChannelID,
				MessageID: "", // This is set by game when sends embed message
			}
			ttcGame := games.NewGameTCC(ttcData, d)
			gs1 = state.NewGameState(ttcGame)
			state.Global.AddGameState(gs1)

			ttcGame.InitReactions(d.Session)
		}
	case PROMT_TTC_QUIT:
		gs := state.Global.GetGameStateByUserID(d.MessageCreate.Author.ID)
		if gs == nil {
			d.ChannelMessageReturn("You are not in game")
			break
		}
		d.Session.ChannelMessageDelete(d.MessageCreate.ChannelID, gs.Game.Data.MessageID)
		state.Global.RemoveGameState(gs)

	default:
		d.ChannelMessageReturn(fmt.Sprintf("%s does not exist", promt))
	}

}

func getPromtAndContent(promt_content string) (string, string) {
	if len(promt_content) < 2 {
		return "", ""
	}

	spaceIndex := len(promt_content)
	promt := ""
	content := ""

	for i, c := range promt_content {
		if c == ' ' {
			spaceIndex = i
			break
		}
	}

	promt = promt_content[1:spaceIndex]

	if spaceIndex+1 < len(promt_content) {
		content = promt_content[spaceIndex+1:]
	}

	return promt, content
}

func getUserMention(user string) (string, error) {
	s := -1
	e := 0

	for i, c := range user {
		if c == '<' {
			if i+1 < len(user) {
				if user[i+1] == '@' {
					s = i
				}
			}

		}

		if s != -1 {
			e += 1
			if c == '>' {
				break
			}
		}
	}

	if s == -1 {
		return "", &errors.UserMentionNotFound{}
	}

	return user[s : s+e], nil
}

func getUserIDFromMention(mention string) string {
	if len(mention) < 4 {
		return ""
	}
	return mention[2 : len(mention)-1]

}
