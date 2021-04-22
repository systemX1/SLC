package game

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	winPoint		= 100
	gamesPerSeries	= 10
)

type players map[string]player

type player struct {
	totalScore int
}

var p player

func MainWindows()  {
	var input string

	for {
		_, _ = fmt.Scanln(&input)
		if input == "p" {
			p.play()
		} else if input == "h" {
			printHelp()
		} else {
			break
		}
	}
}


func (p *player) play()  {
	for i := 1; i <= gamesPerSeries; i++ {
		thisTurnScore := 0
		turnIsOver := false
		for {
			fmt.Println("Choose your operation, 1 - roll, 2 - stay")
			var input string
			_, _ = fmt.Scanln(&input)

			if input == "1" {
				var rollResult int
				p.roll(&rollResult, &turnIsOver)
				if rollResult == 1 {
					thisTurnScore = 0
				} else {
					thisTurnScore += rollResult
				}

				fmt.Printf("rollResult is %d, thisTurnScore is %d\n", rollResult, thisTurnScore)
			} else if input == "2"{
				turnIsOver = true
			} else {
				fmt.Println("Please choose again")
			}

			if turnIsOver {
				p.totalScore += thisTurnScore
				fmt.Printf("Turn %d is over\nyour totalScore score: %d\n", i, p.totalScore)
				break
			}
		}
	}
}

func (p *player) roll(rollResult *int, turnIsOver *bool) () {
	rand.Seed(time.Now().Unix())
	*rollResult = rand.Intn(6) + 1 // A random int in [1, 6]
	if *rollResult == 1 {
		*turnIsOver = true
	}
	return
}

func printHelp()  {
	fmt.Println(`• If the player rolls a 1, the player scores nothing and it becomes the opponent’s turn.
• If the player rolls a number other than 1, the number is added to the player’s turn total and the player’s turn continues.
• If the player holds, the turn total, the sum of the rolls during the turn, is added to the player’s score, and it becomes the opponent’s turn.`)
}



















