package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	PlayerPrompt             = "Please enter the number of players: "
	InvalidPlayerErrorPrompt = "Invalid input for the number of players... Try again."
	InvalidWinnerErrorPrompt = "Invalid input for the winner of the game... Try again."
)

type CLI struct {
	input  *bufio.Scanner
	output io.Writer
	game   Game
}

func NewCLI(input io.Reader, output io.Writer, game Game) *CLI {
	return &CLI{
		input:  bufio.NewScanner(input),
		output: output,
		game:   game,
	}
}

func (c *CLI) PlayPoker() error {
	fmt.Fprint(c.output, PlayerPrompt)
	numPlayersInput, err := strconv.Atoi(c.readLine())
	if err != nil {
		fmt.Fprint(c.output, InvalidPlayerErrorPrompt)
		return err
	}

	c.game.Start(numPlayersInput, c.output)

	input := c.readLine()
	processedInput := strings.Split(input, " ")
	if len(processedInput) != 2 || processedInput[1] != "wins" {
		fmt.Fprint(c.output, InvalidWinnerErrorPrompt)
		return nil
	}
	c.game.Finish(extractWinner(input))

	return nil
}

func extractWinner(input string) string {
	return strings.Split(input, " ")[0]
}

func (c *CLI) readLine() string {
	c.input.Scan()
	return c.input.Text()
}
