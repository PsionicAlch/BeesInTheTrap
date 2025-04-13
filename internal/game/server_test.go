package game

import "testing"

func TestStartupServer(t *testing.T) {
	communication := StartupServer()

	if communication == nil {
		t.Error("Expected StartupServer to return non-nil CommunicationProtocol")
	}
}

func TestRun(t *testing.T) {
	mockProtocol := &MockProtocol{}

	server := &GameServer{
		finished:      false,
		communication: mockProtocol,
		state: GameState{
			Round:  0,
			Hits:   0,
			Stings: 0,
			Player: Player{
				Health:     100,
				MissChance: 0, // Always hit
			},
			Hive: []Bee{
				{
					Type:       QueenBee,
					Health:     1,
					MissChance: 100, // Always miss
				},
			},
		},
	}

	server.run()

	if !mockProtocol.finishedCalled {
		t.Errorf("Game did not finish as expected.")
	}

	if server.state.Round == 0 {
		t.Errorf("Expected at least one round to be completed.")
	}

	server = &GameServer{
		finished:      false,
		communication: mockProtocol,
		state: GameState{
			Player: Player{
				Health:     1,
				MissChance: 100, // Always miss
			},
			Hive: []Bee{
				{
					Type:       QueenBee,
					Health:     100,
					MissChance: 0, // Always hit
				},
			},
		},
	}

	server.run()

	if !mockProtocol.finishedCalled {
		t.Errorf("Game did not finish as expected.")
	}

	if server.state.Round == 0 {
		t.Errorf("Expected at least one round to be completed.")
	}
}

func TestPlayersTurn(t *testing.T) {
	mockProtocol := &MockProtocol{
		currentState: GameState{
			Hits: 100,
		},
		currentMessage: "A message that should never come from the server.",
	}

	server := &GameServer{
		communication: mockProtocol,
		state: GameState{
			Player: Player{
				Health:     100,
				MissChance: 0,
			},
			Hive: []Bee{
				{
					Type:       WorkerBee,
					Health:     1,
					MissChance: 100,
				},
			},
		},
	}

	server.playersTurn()

	if server.state.Hits != 1 {
		t.Errorf("Expected server.state.Hits to be 1.")
	}

	if len(server.state.Hive) != 0 {
		t.Errorf("Expected bee to be removed from hive after death.")
	}

	if mockProtocol.currentMessage == "" {
		t.Errorf("Expected playersTurn to return proper message.")
	}

	if mockProtocol.currentState.Hits != server.state.Hits {
		t.Errorf("Expected playersTurn to return proper game state.")
	}

	server = &GameServer{
		communication: &MockProtocol{},
		state: GameState{
			Player: Player{
				Health:     100,
				MissChance: 100,
			},
			Hive: []Bee{
				{
					Type:       WorkerBee,
					Health:     1,
					MissChance: 100,
				},
			},
		},
	}

	server.playersTurn()

	if server.state.Hits == 1 {
		t.Errorf("Expected server.state.Hits to be 0.")
	}

	if len(server.state.Hive) == 0 {
		t.Errorf("Expected bee to still be alive.")
	}
}

func TestHivesTurn(t *testing.T) {
	server := &GameServer{
		communication: &MockProtocol{},
		state: GameState{
			Player: Player{
				Health:     1,
				MissChance: 100,
			},
			Hive: []Bee{
				{
					Type:       WorkerBee,
					Health:     75,
					MissChance: 100,
				},
			},
		},
	}

	server.hivesTurn()

	if server.state.Stings == 1 {
		t.Errorf("Expected server.state.Stings to be 0.")
	}

	if server.state.Player.Health != 1 {
		t.Errorf("Expected player to still be alive.")
	}

	server = &GameServer{
		communication: &MockProtocol{},
		state: GameState{
			Player: Player{
				Health:     1,
				MissChance: 100,
			},
			Hive: []Bee{
				{
					Type:       WorkerBee,
					Health:     75,
					MissChance: 0,
				},
			},
		},
	}

	server.hivesTurn()

	if server.state.Stings != 1 {
		t.Errorf("Expected server.state.Stings to be 1.")
	}

	if server.state.Player.Health > 0 {
		t.Errorf("Expected player to be dead.")
	}

	mockProtocol := &MockProtocol{
		currentState: GameState{
			Stings: 100,
		},
		currentMessage: "A message that should never come from the server.",
	}

	server = &GameServer{
		communication: mockProtocol,
		state: GameState{
			Player: Player{
				Health:     100,
				MissChance: 100,
			},
			Hive: []Bee{
				{
					Type:       WorkerBee,
					Health:     75,
					MissChance: 0,
				},
			},
		},
	}

	server.hivesTurn()

	if mockProtocol.currentMessage == "" {
		t.Errorf("Expected hivesTurn to return proper message.")
	}

	if mockProtocol.currentState.Stings != server.state.Stings {
		t.Errorf("Expected playersTurn to return proper game state.")
	}
}

type MockProtocol struct {
	finishedCalled bool
	currentMessage string
	currentState   GameState
}

func (m *MockProtocol) Hit() Event {
	return Event{}
}

func (m *MockProtocol) WaitForCPU() Event {
	return Event{}
}

func (m *MockProtocol) WaitForPlayer() {}

func (m *MockProtocol) HitResponse(msg string, state GameState) {
	m.currentMessage = msg
	m.currentState = state
}

func (m *MockProtocol) StingResponse(msg string, state GameState) {
	m.currentMessage = msg
	m.currentState = state
}

func (m *MockProtocol) GameFinishedResponse(msg string, state GameState) {
	m.finishedCalled = true
}
