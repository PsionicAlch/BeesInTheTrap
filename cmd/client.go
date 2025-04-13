package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/PsionicAlch/BeesInTheTrap/internal/game"
)

type IOReader interface {
	ReadString(delim byte) (string, error)
}

type Client struct {
	communication game.Protocol
	reader        IOReader
	writer        io.Writer
	fatalErr      func(error)
}

func createClient(communication game.Protocol, input io.Reader, output io.Writer, errFunc func(error)) *Client {
	return &Client{
		communication: communication,
		reader:        bufio.NewReader(input),
		writer:        output,
		fatalErr:      errFunc,
	}
}

func (c *Client) run() {
	autoPlay := false

	c.printIntro()

	for {
		// Get user's input.
		if !autoPlay {
			command := c.readCommand()

			for !slices.Contains([]string{"hit", "auto"}, command) {
				c.printCommandError()

				command = c.readCommand()
			}

			if command == "auto" {
				autoPlay = true
			}
		}

		event := c.communication.Hit()

		fmt.Fprintln(c.writer, event.Message)

		if event.Type == game.GameFinished {
			c.printGameSummary(event.State)
			return
		}

		event = c.communication.WaitForCPU()

		fmt.Fprintln(c.writer, event.Message)

		if event.Type == game.GameFinished {
			c.printGameSummary(event.State)
			return
		}
	}
}

func (c *Client) printIntro() {
	fmt.Fprintln(c.writer, `🐝 Welcome to Bees In The Trap 🐝

The hive is restless, and you're standing right in the buzz zone.
Armed with nothing but courage and a sharp eye, you must take down the swarm before they sting you into oblivion.

Each turn, you can strike the hive... but beware:
- You might miss entirely.
- They might miss too.
- Every bee type fights differently. Watch out for the Queen.

⚔️ OBJECTIVE:
Destroy the hive before it destroys you.

Commands:
> hit       — Attempt a strike on the hive
> auto      — Let fate decide and simulate the entire game

Let the stinger-slinging begin...`)
}

func (c *Client) readCommand() string {
	fmt.Fprint(c.writer, "> ")

	command, err := c.reader.ReadString('\n')
	if err != nil {
		c.fatalErr(err)
	}

	command = strings.Replace(command, "\n", "", -1)

	return command
}

func (c *Client) printCommandError() {
	fmt.Fprintln(c.writer, `Invalid Command!

Commands:
> hit       — Attempt a strike on the hive
> auto      — Let fate decide and simulate the entire game`)
}

func (c *Client) printGameSummary(state game.GameState) {
	var playersFate string
	var queensFate string
	var workerBeesAlive uint
	var droneBeesAlive uint
	var finalCommentary string

	if state.Player.Health <= 0 {
		playersFate = "You perished in the swarm."
		finalCommentary = "The hive overwhelmed you. Your story ends in silence..."
	} else {
		playersFate = "You survived the hive!"
		finalCommentary = "Victory! The hive has fallen. Peace returns to the meadow."
	}

	queensIndex := slices.IndexFunc(state.Hive, func(bee game.Bee) bool {
		return bee.Type == game.QueenBee
	})

	if queensIndex == -1 {
		queensFate = "Unsure..."
	} else {
		if state.Hive[queensIndex].Health > 0 {
			queensFate = "Alive"
		} else {
			queensFate = "Dead"
		}
	}

	for _, bee := range state.Hive {
		if bee.Type == game.WorkerBee {
			workerBeesAlive++
		} else if bee.Type == game.DroneBee {
			droneBeesAlive++
		}
	}

	fmt.Fprintf(c.writer, `
📜 Game Summary
============================
Rounds played : %d
Total hits    : %d
Total stings  : %d

👤 Player Status
----------------------------
Final Health  : %d
Fate          : %s

🐝 Hive Status
----------------------------
Queen Bee     : %s
Worker Bees   : %d remaining
Drone Bees    : %d remaining

%s
`,
		state.Round,
		state.Hits,
		state.Stings,
		state.Player.Health,
		playersFate,
		queensFate,
		workerBeesAlive,
		droneBeesAlive,
		finalCommentary,
	)
}
