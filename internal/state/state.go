package state

import (
	"dbot/internal/games"
	"fmt"

	"github.com/google/uuid"
)

var Global *GlobalState = NewGlobalState()

type GlobalState struct {
	Games []*GameState
}

type GameState struct {
	ID        uuid.UUID
	Game      *games.TTCGame
	ChannelID string
	MessageID string
}

func NewGlobalState() *GlobalState {
	return &GlobalState{
		Games: []*GameState{},
	}
}

func NewGameState(game *games.TTCGame) *GameState {
	return &GameState{
		ID:        uuid.New(),
		Game:      game,
		ChannelID: "",
		MessageID: "",
	}
}

func (global *GlobalState) GetGameStateByGameID(id uuid.UUID) *GameState {
	for _, gs := range global.Games {
		if gs.ID == id {
			return gs
		}
	}
	return nil
}

func (global *GlobalState) GetGameStateByMessageID(messageID string) *GameState {
	for _, gs := range global.Games {
		if gs.Game.Data.MessageID == messageID {
			return gs
		}
	}
	return nil
}

func (global *GlobalState) AddGameState(gs *GameState) {
	global.Games = append(global.Games, gs)
}

func (global *GlobalState) GetGameStateByUserID(userID string) *GameState {
	for _, gs := range global.Games {
		if gs.Game.Data.UserID_1 == userID || gs.Game.Data.UserID_2 == userID {
			return gs
		}
	}
	return nil
}

func (global *GlobalState) RemoveGameState(gs *GameState) {
	for i, gstate := range global.Games {
		if gs == gstate {
			fmt.Println("GameState Removed")
			global.Games = append(global.Games[:i], global.Games[i+1:]...)
			fmt.Printf("%+v\n", global.Games)
			break
		}
	}
}
