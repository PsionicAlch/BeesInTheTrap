package main

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/PsionicAlch/BeesInTheTrap/internal/game"
)

func TestCreateClient(t *testing.T) {
	client := createClient(&MockProtocol{}, os.Stdin, os.Stdout, func(err error) {})

	if client == nil {
		t.Error("createClient failed to return non-nil client instance.")
	}
}

func TestRun(t *testing.T) {
	input := strings.NewReader("wrong\nauto\n")
	output := &bytes.Buffer{}
	protocol := &MockProtocol{
		events: []game.Event{
			{Type: game.PlayerAttack, Message: "You killed a worker bee.", State: game.GameState{}},
			{Type: game.GameFinished, Message: "You lost!", State: game.GameState{}},
		},
	}

	client := createClient(protocol, input, output, func(err error) {})

	client.run()

	commandErrorStr := `Invalid Command!

Commands:
> hit       — Attempt a strike on the hive
> auto      — Let fate decide and simulate the entire game`

	result := output.String()

	if !strings.Contains(result, commandErrorStr) {
		t.Error("Expected output to contain command error string.")
	}

	if !strings.Contains(result, "You killed a worker bee.") {
		t.Error("Expected output to contain \"You killed a worker bee.\"")
	}

	if !strings.Contains(result, "You lost!") {
		t.Error("Expected output to contain \"You lost!\".")
	}

	input = strings.NewReader("auto\n")
	output = &bytes.Buffer{}
	protocol = &MockProtocol{
		events: []game.Event{
			{Type: game.GameFinished, Message: "You won!", State: game.GameState{}},
		},
	}

	client = createClient(protocol, input, output, func(err error) {})

	client.run()

	result = output.String()

	if !strings.Contains(result, "You won!") {
		t.Error("Expected output to contain \"You won!\".")
	}
}

func TestReadCommand(t *testing.T) {
	var testErr string

	client := &Client{
		communication: &game.CommunicationProtocol{},
		reader:        &MockReader{},
		writer:        os.Stdout,
		fatalErr: func(err error) {
			testErr = err.Error()
		},
	}

	client.readCommand()

	if testErr != "simulated read error" {
		t.Error("Expected readCommand to fail.")
	}
}

func TestPrintGameSummary(t *testing.T) {
	scenarios := []struct {
		state       game.GameState
		expectedStr string
	}{
		{game.GameState{Player: game.Player{Health: -1}}, "You perished in the swarm."},
		{game.GameState{Player: game.Player{Health: -1}}, "The hive overwhelmed you. Your story ends in silence..."},
		{game.GameState{Player: game.Player{Health: 1}}, "You survived the hive!"},
		{game.GameState{Player: game.Player{Health: 1}}, "Victory! The hive has fallen. Peace returns to the meadow."},
		{game.GameState{}, "Unsure..."},
		{game.GameState{Hive: []game.Bee{{Type: game.QueenBee, Health: 1, MissChance: 0}}}, "Alive"},
		{game.GameState{Hive: []game.Bee{{Type: game.QueenBee, Health: -1, MissChance: 0}}}, "Dead"},
		{game.GameState{Hive: []game.Bee{{Type: game.WorkerBee, Health: 1, MissChance: 0}}}, "Worker Bees   : 1 remaining"},
		{game.GameState{Hive: []game.Bee{{Type: game.DroneBee, Health: 1, MissChance: 0}}}, "Drone Bees    : 1 remaining"},
	}

	for _, scenario := range scenarios {
		output := &bytes.Buffer{}
		client := &Client{
			communication: &game.CommunicationProtocol{},
			reader:        &MockReader{},
			writer:        output,
			fatalErr:      func(err error) {},
		}

		client.printGameSummary(scenario.state)

		result := output.String()

		if !strings.Contains(result, scenario.expectedStr) {
			t.Errorf("Expected output to contain \"%s\". Result:\n\n%s", scenario.expectedStr, result)
		}
	}
}

type MockProtocol struct {
	events []game.Event
	index  int
}

func (protocol *MockProtocol) Hit() game.Event {
	event := protocol.events[protocol.index]
	protocol.index++

	return event
}

func (protocol *MockProtocol) WaitForCPU() game.Event {
	event := protocol.events[protocol.index]
	protocol.index++

	return event
}

func (protocol *MockProtocol) WaitForPlayer() {}

func (protocol *MockProtocol) HitResponse(string, game.GameState) {}

func (protocol *MockProtocol) StingResponse(string, game.GameState) {}

func (protocol *MockProtocol) GameFinishedResponse(string, game.GameState) {}

type MockReader struct{}

func (reader *MockReader) ReadString(delim byte) (string, error) {
	return "", errors.New("simulated read error")
}
