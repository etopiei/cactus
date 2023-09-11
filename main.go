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
	fmt.Println("id name cactus v0.3.0")
	fmt.Println("id author etopiei (Lachlan Jacob)")
	fmt.Println("uciok")
}

func orderMoves(moves []*chess.Move) []*chess.Move {
	// POC with Random move ordering, but put this here for program structure
	return moves
}

func findMoveOnDepth(position *chess.Position, depth int, alpha int, beta int) int {
	if depth == 0 {
		return evaluatePosition(position)
	}

	for _, move := range orderMoves(position.ValidMoves()) {
		newPosition := position.Update(move)
		evaluation := -findMoveOnDepth(newPosition, depth-1, -beta, -alpha)
		if evaluation >= beta {
			return beta
		}
		alpha = max(evaluation, alpha)
	}
	return alpha
}

func findMove(position *chess.Position) chess.Move {
	moves := position.ValidMoves()

	// TODO: Handle the case where there are no valid moves!
	// At the moment this will just return an uninitialized 'bestMove'

	best := -100000
	var bestMove chess.Move

	for _, move := range moves {
		newPosition := position.Update(move)
		score := -findMoveOnDepth(newPosition, 4, -10000, 100000)
		if score >= best {
			best = score
			bestMove = *move
		}
	}

	return bestMove
}

func pieceToValue(piece chess.Piece) int {
	switch piece.Type() {
	case chess.King:
		return 20000
	case chess.Queen:
		return 900
	case chess.Rook:
		return 500
	case chess.Bishop:
		return 330
	case chess.Knight:
		return 320
	case chess.Pawn:
		return 100
	default:
		return 0
	}
}

func squareTableValue(piece chess.Piece, square chess.Square) int {
	var boardIndex int
	if piece.Color() == chess.White {
		boardIndex = squareToIndex(square)
	} else {
		boardIndex = indexOfMirrorSquare(square)
	}

	switch piece.Type() {
	case chess.Bishop:
		return BISHOP[boardIndex]
	case chess.King:
		// TODO: Add midgame/endgame distinction with King
		return KING_MIDGAME[boardIndex]
	case chess.Pawn:
		return PAWN[boardIndex]
	case chess.Knight:
		return KNIGHT[boardIndex]
	case chess.Queen:
		return QUEEN[boardIndex]
	default:
		return 0
	}
}

func squareToXDirYDir(square chess.Square) (int, int) {
	var xDir int
	var yDir int
	switch square.File() {
	case chess.FileA:
		xDir = 0
	case chess.FileB:
		xDir = 1
	case chess.FileC:
		xDir = 2
	case chess.FileD:
		xDir = 3
	case chess.FileE:
		xDir = 4
	case chess.FileF:
		xDir = 5
	case chess.FileG:
		xDir = 6
	case chess.FileH:
		xDir = 7
	}

	switch square.Rank() {
	case chess.Rank1:
		yDir = 0
	case chess.Rank2:
		yDir = 1
	case chess.Rank3:
		yDir = 2
	case chess.Rank4:
		yDir = 3
	case chess.Rank5:
		yDir = 4
	case chess.Rank6:
		yDir = 5
	case chess.Rank7:
		yDir = 6
	case chess.Rank8:
		yDir = 7
	}

	return xDir, yDir
}

func squareToIndex(square chess.Square) int {
	xDir, yDir := squareToXDirYDir(square)
	return (yDir * 8) + xDir
}

func indexOfMirrorSquare(square chess.Square) int {
	xDir, yDir := squareToXDirYDir(square)
	return (7 - xDir) + ((7 - yDir) * 8)
}

func evaluatePosition(position *chess.Position) int {
	score := 0
	for square, piece := range position.Board().SquareMap() {
		if piece.Color() == position.Turn() {
			score += pieceToValue(piece) + squareTableValue(piece, square)
		} else {
			score -= pieceToValue(piece) + squareTableValue(piece, square)
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

				// TODO: Handle move string better (currently promotion doesn't work)
				moveStr := move.S1().String() + move.S2().String()

				fmt.Println("bestmove", moveStr)
			}
		}
	}
}
