package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/notnil/chess"
)

func printIdInfo() {
	fmt.Println("id name cactus 0.1")
	fmt.Println("id author etopiei (Lachlan Jacob)")
	fmt.Println("uciok")
}

func findMoveOnDepth(position *chess.Position, depth int, color chess.Color) int {
	if depth == 0 {
		return evaluatePosition(position, color)
	}

	worst := 10000

	for _, move := range position.ValidMoves() {
		newPosition := position.Update(move)
		score := findMoveOnDepth(newPosition, depth-1, color)
		if score < worst {
			worst = score
		}
	}
	return worst
}

func findMove(position *chess.Position) chess.Move {
	moves := position.ValidMoves()

	// Here take the worst evaluation for the subtree of a move
	// and take the best worst evaluation. (minimax)
	best := -10000
	var bestMove chess.Move

	// TODO: Handle the case where there are no valid moves!
	// At the moment this will just return an uninitialized 'bestMove'

	for _, move := range moves {
		newPosition := position.Update(move)
		score := findMoveOnDepth(newPosition, 4, newPosition.Turn())
		if score > best {
			best = score
			bestMove = *move
		}
	}

	return bestMove
}

func pieceToValue(piece chess.Piece) int {
	switch piece.Type() {
	case chess.King:
		return 30
	case chess.Queen:
		return 10
	case chess.Rook:
		return 5
	case chess.Bishop:
		return 3
	case chess.Knight:
		return 3
	case chess.Pawn:
		return 1
	default:
		return 0
	}
}

func evaluatePosition(position *chess.Position, evaluateFor chess.Color) int {
	// TODO: Make the evaluation function smarter
	score := 0
	for _, piece := range position.Board().SquareMap() {
		if piece.Color() == evaluateFor {
			score += pieceToValue(piece)
		} else {
			score -= pieceToValue(piece)
		}
	}
	return score
}

func main() {
	printIdInfo()
	game := chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	stdin := bufio.NewScanner(os.Stdin)

	f, err := os.Create("cactus-log.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for {
		for stdin.Scan() {
			// Got a line of text in stdin.Text()
			// Check for line matching position
			command := stdin.Text()
			commandParts := strings.Split(command, " ")

			fmt.Fprintln(f, command)

			if commandParts[0] == "ping" {
				fmt.Println(commandParts[1])
			}

			if commandParts[0] == "isready" {
				fmt.Println("readyok")
			}

			if commandParts[0] == "position" && len(commandParts) > 2 {
				// TODO: Make this handle when the position is being set not added to
				game.MoveStr(commandParts[len(commandParts)-1])
			}

			if commandParts[0] == "go" {
				move := findMove(game.Position())
				game.Move(&move)
				moveStr := move.S1().String() + move.S2().String()
				fmt.Println("bestmove", moveStr)
			}
		}
	}
}
