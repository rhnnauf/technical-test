package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type player struct {
	dices []int
	idx   int
	score int
	isPlaying bool
}

// function to randomize each player dices, ranging from (1-6)
func (p *player) randomizeDices() {
	for i := 0; i < len(p.dices); i++{
		randNumber := rand.Intn(6) + 1
		p.dices[i] = randNumber
	}
}

// function to remove value "1" and "6" from player dices, then evaluating the score, also return how many "1" occurrences
func (p *player) evaluate() int {
	var result []int
	oneCounter := 0

	for _, element := range p.dices{
		// if dice == 6, then increment the player score
		if element == 6{
			p.score++
		// if dice == 1, then increment the oneCounter to pass on the next player
		}else if element == 1{
			oneCounter++
		// else append to the new slices for the value except "1" and "6"
		}else{
			result = append(result, element)
		}
	}

	// set the newly created slice to the current player dices
	p.dices = result
	
	return oneCounter
}

// appending "1" to the player
func (p *player) addOneToTheCurrentPlayer(oneOccurrences int) {
	for i := 0; i < oneOccurrences; i++{
		p.dices = append(p.dices, 1)
	}
}

type diceGame struct {
	players   []player
	numPlayer int
	numDicePerPlayer int
	round     int
}

// game flow:
// 1. randomize dices for each player
// 2. display each player dices in step(1) + display game round
// 3. evaluate dices for each player
	// (a). removing "6" from each player dices, then incrementing the player score based on "6" occurrences
	// (b). removing "1" from each player dices, then move the "1" to the next player if there is an occurrence
// 4. display after evaluation in step(3)
// 5. check has winner, if not repeat step(1) || if only 1 player that has dices left in their slice then break

func (d *diceGame) start() {
	// keep playing until there is only one player remaining
	for{
		rand.Seed(time.Now().UnixNano())

		// step(1): randomize dices for each player
		for i := 0; i < len(d.players); i++{
			d.players[i].randomizeDices()
		}
	
		// step(2): display each player dices in step(1) + display game round
		d.round++
		fmt.Printf("\nGiliran %d lempar dadu\n", d.round)
		d.printPlayerScoreAndDices()

		// step(3): evaluate dices for each player
		oneCounterIdxArr := make([]int, len(d.players))

		for i := 0; i < len(d.players); i++{
			// evaluate the score
			oneCounter := d.players[i].evaluate()
			
			// push how many "1" occurrences for each player, for reference to next player
			oneCounterIdxArr[i] = oneCounter
		}

		// move "1" to the next available player
		for i := 0; i < len(d.players); i++{
			var nextPlayer int
			
			// if it is the rightmost index, then start again at index 0
			if i == len(d.players) - 1{
				nextPlayer = 0
			}else{
				nextPlayer = i + 1
			}

			// loop until we find the next available player
			for nextPlayer != i {

				// add "1" to the next player
				if d.players[nextPlayer].isPlaying {
					d.players[nextPlayer].addOneToTheCurrentPlayer(oneCounterIdxArr[i])
					break
				}
				
				if nextPlayer == len(d.players) - 1{
					nextPlayer = 0
				}else{
					nextPlayer++
				}
			}
		}

		// step(4): display each player dices after evaluation in step(3)
		fmt.Println("Setelah evaluasi: ")
		d.printPlayerScoreAndDices()

		// step(5): check has winner
		for idx, p := range d.players{
			if len(p.dices) == 0 {
				// if the player is still playing, then change the isPlaying value to false
				if d.players[idx].isPlaying {
					d.players[idx].isPlaying = false

					// decrement the current player
					d.numPlayer--
				}
			}
		}

		// if the current player is <= 1, then stop the game
		if d.numPlayer <= 1{
			d.printWinner()
			break
		}
	}
}

// print winner
func (d diceGame) printWinner() {
	maxScore := 0

	// multiple winner possibilities
	var winnerIdx []int

	for i, p := range d.players{
		// if current player score > maxScore then reset the slice to the player index
		if p.score > maxScore {
			winnerIdx = []int{p.idx + 1}
			maxScore = p.score
		// if we have multiple highscore, then append to the slice
		}else if p.score == maxScore{
			winnerIdx = append(winnerIdx, i + 1)
		}
	}

	fmt.Printf("Game dimenangkan oleh pemain #%d karena memiliki poin lebih banyak dari pemain lainnya.", winnerIdx)
}

// print the board
func (d diceGame) printPlayerScoreAndDices() {
	for _, p := range d.players{
		fmt.Printf("Pemain %d (%d): ", p.idx+1, p.score)

		if len(p.dices) == 0{
			fmt.Printf("_ (Berhenti bermain karena tidak memiliki dadu)")
		}else{
			for _, d := range p.dices{
				fmt.Printf("%d, ", d)
			}
		}

		fmt.Printf("\n")
	}
}

// new dice game constructor
func newDiceGame(inputPlayer, inputDice int) *diceGame {
	var p []player
	
	// initialize player for the game struct
	for i := 0; i < inputPlayer; i++{
		p = append(p, player{
			dices: make([]int, inputDice),
			idx: i,
			score: 0,
			isPlaying: true,
		})
	}

	return &diceGame{
		players: p,
		numPlayer: inputPlayer,
		numDicePerPlayer: inputDice,
		round: 0,
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Input jumlah pemain: ")
	scanner.Scan()
	inputPlayer, err := strconv.ParseInt(scanner.Text(), 10, 64)

	if err != nil {
		fmt.Println("input is not an integer")
		return
	}

	fmt.Printf("Input jumlah dadu: ")
	scanner.Scan()
	inputDice, err := strconv.ParseInt(scanner.Text(), 10, 64)
	if err != nil {
		fmt.Println("input is not an integer")
		return
	}

	// initialize new game, based on the input
	game := newDiceGame(int(inputPlayer), int(inputDice))

	fmt.Printf("Pemain = %d, Dadu = %d\n", inputPlayer, inputDice)

	// start game
	game.start()
}