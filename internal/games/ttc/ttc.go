package ttc

import (
	"dbot/internal/errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	SIGN_X  byte = 88
	SIGN_O  byte = 79
	NO_SIGN byte = 45 // -
	// NO_SIGN byte = 32 // space
)

var winpos [8][3]int = [8][3]int{
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
	{0, 4, 8},
	{2, 4, 6},
}

//

type Table [9]byte

//

type TicTacToe struct {
	Table          Table
	CurrentSign    byte
	availableTurns int
}

//
//

func New() TicTacToe {
	var table Table
	for i := range table {
		table[i] = NO_SIGN
	}

	return TicTacToe{
		Table:          table,
		CurrentSign:    SIGN_X,
		availableTurns: 9,
	}
}

//
//

func (ttc *TicTacToe) SwitchSign() {
	if ttc.CurrentSign == SIGN_X {
		ttc.CurrentSign = SIGN_O
	} else {
		ttc.CurrentSign = SIGN_X
	}
}

func (ttc *TicTacToe) SetSignAtPosition(pos int) error {
	if ttc.availableTurns < 1 {
		return &errors.TTCNoAvailablePosition{}
	}
	if pos < 0 || pos > 8 {
		return &errors.InvalidRange{Start: 0, End: 8}
	}
	if ttc.Table[pos] != NO_SIGN {
		return &errors.TTCPositionTaken{Position: pos}
	}

	ttc.Table[pos] = ttc.CurrentSign
	ttc.availableTurns -= 1
	ttc.SwitchSign()
	return nil
}

func (ttc *TicTacToe) SetSign(sign byte, pos int) error {
	if pos < 0 || pos > 8 {
		return &errors.InvalidRange{Start: 0, End: 8}
	}
	ttc.Table[pos] = sign
	ttc.availableTurns -= 1
	return nil
}

func (ttc *TicTacToe) CheckWinner() bool {
	for i := range 8 {
		counter := 0
		for j := range 3 {
			if ttc.Table[winpos[i][j]] == SIGN_X {
				counter += 1
			} else if ttc.Table[winpos[i][j]] == SIGN_O {
				counter -= 1
			}
		}
		if counter == 3 || counter == -3 {
			return true
		}
	}

	return false
}

func (ttc *TicTacToe) Reset(pos int) {
	for i := range ttc.Table {
		ttc.Table[i] = 0
	}
	ttc.CurrentSign = SIGN_O
	ttc.availableTurns = 9
}

func (ttc *TicTacToe) CurrentSignStr() string {
	return SignToStr(ttc.CurrentSign)
}

func (ttc *TicTacToe) AvailableTurns() int {
	return ttc.availableTurns
}

func (ttc *TicTacToe) TableToString() string {
	row1 := fmt.Sprintf("| %c | %c | %c |\n", ttc.Table[0], ttc.Table[1], ttc.Table[2])
	row2 := fmt.Sprintf("| %c | %c | %c |\n", ttc.Table[3], ttc.Table[4], ttc.Table[5])
	row3 := fmt.Sprintf("| %c | %c | %c |\n", ttc.Table[6], ttc.Table[7], ttc.Table[8])
	// spacer := "--------\n"
	spacer := ""
	return row1 + spacer + row2 + spacer + row3
}

func (ttc *TicTacToe) AutoSet() (int, error) {
	if ttc.availableTurns < 1 {
		return -1, &errors.TTCNoAvailablePosition{}
	}
	openPos := map[int][]int{}
	otherSignWinPos := []int{-1, -1}

	for i := range 8 {
		csCounter := 0
		osCounter := 0
		for j := range 3 {
			sign := ttc.Table[winpos[i][j]]
			if sign == NO_SIGN {
				openPos[i] = append(openPos[i], j)

			} else if sign == ttc.CurrentSign {
				csCounter += 1
			} else {
				osCounter += 1
			}
		}

		if csCounter == 2 && len(openPos[i]) == 1 {
			pos := winpos[i][openPos[i][0]]
			ttc.SetSignAtPosition(pos)
			return pos, nil

		} else if osCounter == 2 && len(openPos[i]) == 1 {
			otherSignWinPos[0] = i
			otherSignWinPos[1] = openPos[i][0]
		}
	}

	if otherSignWinPos[0] != -1 {
		pos := winpos[otherSignWinPos[0]][otherSignWinPos[1]]
		ttc.SetSignAtPosition(pos)
		return pos, nil
	}

	openPosArray := []int{}
	for k, v := range openPos {
		for x := range v {
			openPosArray = append(openPosArray, winpos[k][x])
		}
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	pos := openPosArray[r.Intn(len(openPosArray))]
	ttc.SetSignAtPosition(pos)
	return pos, nil
}

func SignToStr(sign byte) string {
	return fmt.Sprintf("%c", sign)
}

func FillTableRandom(table *Table) {
	positions := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	sign := SIGN_X

	for len(positions) > 0 {
		if sign == SIGN_X {
			sign = SIGN_O
		} else {
			sign = SIGN_X
		}
		n := r.Intn(len(positions))
		table[positions[n]] = sign

		positions = append(positions[:n], positions[n+1:]...)
	}
}
