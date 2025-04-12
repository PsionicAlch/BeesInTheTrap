package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/PsionicAlch/BeesInTheTrap/internal/game"
)

func main() {
	autoPlay := false
	input := bufio.NewReader(os.Stdin)
	communication := game.StartupServer()

	printIntro()

	for {
		// Get user's input.
		if !autoPlay {
			command := readCommand(input)

			for !slices.Contains([]string{"hit", "auto"}, command) {
				printCommandError()

				command = readCommand(input)
			}

			if command == "auto" {
				autoPlay = true
			}
		}

		event := communication.Hit()

		fmt.Println(event.Message)

		if event.Type == game.GameFinished {
			printGameSummary(event.State)
			return
		}

		event = communication.WaitForCPU()

		fmt.Println(event.Message)

		if event.Type == game.GameFinished {
			printGameSummary(event.State)
			return
		}
	}
}

// printIntro displays the opening banner and game instructions.
func printIntro() {
	fmt.Println(`🐝 Welcome to Bees In The Trap 🐝

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

// readCommand reads a command from standard input.
func readCommand(input *bufio.Reader) string {
	fmt.Print("> ")

	command, err := input.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	command = strings.Replace(command, "\n", "", -1)

	return command
}

// printCommandError displays an error message when the input is invalid.
func printCommandError() {
	fmt.Println(`Invalid Command!

Commands:
> hit       — Attempt a strike on the hive
> auto      — Let fate decide and simulate the entire game`)
}

// printGameSummary summarizes the outcome of the game.
func printGameSummary(state game.GameState) {
	var playersFate string
	var queensFate string
	var workerBeesAlive uint
	var droneBeesAlive uint
	var finalCommentary string

	if state.Player.Health <= 0 {
		playersFate = "You perished in the swarm."
	} else {
		playersFate = "You survived the hive!"
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

	if state.Player.Health > 0 {
		finalCommentary = "Victory! The hive has fallen. Peace returns to the meadow."
	} else {
		finalCommentary = "The hive overwhelmed you. Your story ends in silence..."
	}

	fmt.Printf(`
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
