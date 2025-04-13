package game

import (
	"testing"
	"time"
)

// TestCreateCommunicationProtocol verifies that the protocol and its channels are properly initialized.
func TestCreateCommunicationProtocol(t *testing.T) {
	communication := createCommunicationProtocol()

	if communication == nil {
		t.Fatal("Expected non-nil CommunicationProtocol.")
	}

	if communication.hitSignal == nil {
		t.Error("Expected hitSignal to be initialized.")
	}

	if communication.eventChannel == nil {
		t.Error("Expected eventChannel to be initialized.")
	}
}

// TestHit ensures the Hit method waits for a hit signal and correctly receives an event.
func TestHit(t *testing.T) {
	communication := createCommunicationProtocol()

	expectedEvent := Event{
		Type:    PlayerAttack,
		Message: "Test Message",
		State: GameState{
			Round: 25,
		},
	}

	go func() {
		select {
		case <-communication.hitSignal:
			communication.eventChannel <- expectedEvent
		case <-time.After(time.Second * 3):
			t.Errorf("Timed out waiting for hitSignal")
		}
	}()

	actualEvent := communication.Hit()

	if actualEvent.Type != expectedEvent.Type && actualEvent.Message != expectedEvent.Message && actualEvent.State.Round != expectedEvent.State.Round {
		t.Errorf("Expected event does not match received event.")
	}
}

// TestWaitForCPU ensures CPU event messages are correctly received from the event channel.
func TestWaitForCPU(t *testing.T) {
	communication := createCommunicationProtocol()

	expectedEvent := Event{
		Type:    PlayerAttack,
		Message: "Test Message",
		State: GameState{
			Round: 25,
		},
	}

	go func() {
		communication.eventChannel <- expectedEvent
	}()

	actualEvent := communication.WaitForCPU()

	if actualEvent.Type != expectedEvent.Type && actualEvent.Message != expectedEvent.Message && actualEvent.State.Round != expectedEvent.State.Round {
		t.Errorf("Expected event does not match received event.")
	}
}

// TestWaitForPlayer ensures that WaitForPlayer properly blocks until a signal is received.
func TestWaitForPlayer(t *testing.T) {
	communication := createCommunicationProtocol()

	go func() {
		communication.hitSignal <- struct{}{}
	}()

	communication.WaitForPlayer()
}

// TestHitResponse ensures that HitResponse sends the correct event over the channel.
func TestHitResponse(t *testing.T) {
	communication := createCommunicationProtocol()

	msg := "Test Message"
	state := GameState{
		Round: 25,
	}

	go func() {
		communication.HitResponse(msg, state)
	}()

	select {
	case event := <-communication.eventChannel:
		if event.Type != PlayerAttack && event.Message != msg && event.State.Round != state.Round {
			t.Errorf("Expected event does not match received event.")
		}
	case <-time.After(time.Second * 3):
		t.Errorf("Timed out waiting for hitSignal")
	}
}

// TestStingResponse ensures the StingResponse method delivers the correct event.
func TestStingResponse(t *testing.T) {
	communication := createCommunicationProtocol()

	msg := "Test Message"
	state := GameState{
		Round: 25,
	}

	go func() {
		communication.StingResponse(msg, state)
	}()

	select {
	case event := <-communication.eventChannel:
		if event.Type != HiveAttack && event.Message != msg && event.State.Round != state.Round {
			t.Errorf("Expected event does not match received event.")
		}
	case <-time.After(time.Second * 3):
		t.Errorf("Timed out waiting for hitSignal")
	}
}

// TestGameFinishedResponse ensures the GameFinishedResponse sends the correct event at game end.
func TestGameFinishedResponse(t *testing.T) {
	communication := createCommunicationProtocol()

	msg := "Test Message"
	state := GameState{
		Round: 25,
	}

	go func() {
		communication.GameFinishedResponse(msg, state)
	}()

	select {
	case event := <-communication.eventChannel:
		if event.Type != GameFinished && event.Message != msg && event.State.Round != state.Round {
			t.Errorf("Expected event does not match received event.")
		}
	case <-time.After(time.Second * 3):
		t.Errorf("Timed out waiting for hitSignal")
	}
}
