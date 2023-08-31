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

func findMoveOnDepth(game *chess.Game, depth int, color chess.Color) int {
	if depth == 0 {
		return evaluatePosition(game, color)
	}

	worst := 10000

	for _, move := range game.ValidMoves() {
		newGame := game.Clone()
		newGame.Move(move)
		score := findMoveOnDepth(newGame, depth-1, color)
		if score < worst {
			worst = score
		}
	}
	return worst
}

func findMove(game *chess.Game) string {
	moves := game.ValidMoves()

	// Here take the worst evaluation for the subtree of a move
	// and take the best worst evaluation. (minimax)
	best := -10000
	var bestMove chess.Move

	for _, move := range moves {
		newGame := game.Clone()
		newGame.Move(move)
		score := findMoveOnDepth(newGame, 3, game.Position().Turn())
		if score > best {
			best = score
			bestMove = *move
		}
	}

	// Apply move and return it
	game.Move(&bestMove)
	return bestMove.S1().String() + bestMove.S2().String()
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

func evaluatePosition(game *chess.Game, evaluateFor chess.Color) int {
	// TODO: Here we need more information than just pieces on the board
	// We also need to consider if a position is checkmate
	score := 0
	for _, piece := range game.Position().Board().SquareMap() {
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
				game.MoveStr(commandParts[len(commandParts)-1])
			}

			if commandParts[0] == "go" {
				move := findMove(game)
				fmt.Println("bestmove", move)
			}
		}
	}
}
