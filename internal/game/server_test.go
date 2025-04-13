package game

import "testing"

// TestStartupServer ensures that the StartupServer function returns a valid (non-nil) CommunicationProtocol.
func TestStartupServer(t *testing.T) {
	communication := StartupServer()

	if communication == nil {
		t.Error("Expected StartupServer to return non-nil CommunicationProtocol")
	}
}

// TestRun verifies the full run loop for both player-win and player-loss scenarios.
func TestRun(t *testing.T) {
	mockProtocol := &MockProtocol{}

	// Scenario: Player wins by killing a bee.
	server := &GameServer{
		finished:      false,
		communication: mockProtocol,
		state: GameState{
			Round:  0,
			Hits:   0,
			Stings: 0,
			Player: Player{
				Health:     100,
				MissChance: 0, // Always hit.
			},
			Hive: []Bee{
				{
					Type:       QueenBee,
					Health:     1,
					MissChance: 100, // Always miss.
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

	// Scenario: Player loses due to hive attacks.
	server = &GameServer{
		finished:      false,
		communication: mockProtocol,
		state: GameState{
			Player: Player{
				Health:     1,
				MissChance: 100, // Always miss.
			},
			Hive: []Bee{
				{
					Type:       QueenBee,
					Health:     100,
					MissChance: 0, // Always hit.
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

// TestPlayersTurn simulates the player's attack phase with various outcomes:
// - A successful attack.
// - A missed attack.
func TestPlayersTurn(t *testing.T) {
	// Successful hit and bee removed.
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

	// Missed attack, bee remains.
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

// TestHivesTurn simulates the hive’s retaliation phase against the player:
// - One test ensures the bee misses.
// - One ensures the player is killed.
// - One verifies correct state/message reporting.
func TestHivesTurn(t *testing.T) {
	// Hive misses the player.
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

	// Hive kills the player.
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

	// Confirm that message/state was sent back correctly.
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

// --- Mock Protocol ---

// MockProtocol implements the Protocol interface with test-friendly behavior.
type MockProtocol struct {
	finishedCalled bool
	currentMessage string
	currentState   GameState
}

// Hit is unused in these tests but required to implement the interface
func (m *MockProtocol) Hit() Event {
	return Event{}
}

// WaitForCPU is unused in server tests.
func (m *MockProtocol) WaitForCPU() Event {
	return Event{}
}

// WaitForPlayer is unused in server tests.
func (m *MockProtocol) WaitForPlayer() {}

// HitResponse simulates sending a result from a player hit.
func (m *MockProtocol) HitResponse(msg string, state GameState) {
	m.currentMessage = msg
	m.currentState = state
}

// StingResponse simulates sending a result from a bee attack.
func (m *MockProtocol) StingResponse(msg string, state GameState) {
	m.currentMessage = msg
	m.currentState = state
}

// GameFinishedResponse records that the game has concluded.
func (m *MockProtocol) GameFinishedResponse(msg string, state GameState) {
	m.finishedCalled = true
}
