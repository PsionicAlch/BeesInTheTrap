package game

type Protocol interface {
	Hit() Event
	WaitForCPU() Event
	WaitForPlayer()
	HitResponse(string, GameState)
	StingResponse(string, GameState)
	GameFinishedResponse(string, GameState)
}

// CommunicationProtocol encapsulates the message flow between
// the server-side game engine and the client.
//
// The protocol uses channels to coordinate player input and emit
// game events in response to game logic execution.
type CommunicationProtocol struct {
	hitSignal    chan struct{}
	eventChannel chan Event
}

// createCommunicationProtocol initializes and returns a new CommunicationProtocol instance.
func createCommunicationProtocol() *CommunicationProtocol {
	return &CommunicationProtocol{
		hitSignal:    make(chan struct{}),
		eventChannel: make(chan Event),
	}
}

// Hit is called when the player takes an action. It blocks until the server
// processes the player's move and returns an Event describing the outcome.
func (protocol *CommunicationProtocol) Hit() Event {
	protocol.hitSignal <- struct{}{}

	return <-protocol.eventChannel
}

// WaitForCPU blocks until the hive has finished it's turn.
func (protocol *CommunicationProtocol) WaitForCPU() Event {
	return <-protocol.eventChannel
}

// waitForPlayer blocks until a player action is received.
func (protocol *CommunicationProtocol) WaitForPlayer() {
	<-protocol.hitSignal
}

// hitResponse sends an event indicating the player has performed an attack.
func (protocol *CommunicationProtocol) HitResponse(msg string, state GameState) {
	event := Event{
		Type:    PlayerAttack,
		Message: msg,
		State:   state,
	}

	protocol.eventChannel <- event
}

// stingResponse sends an event indicating the hive has performed an attack.
func (protocol *CommunicationProtocol) StingResponse(msg string, state GameState) {
	event := Event{
		Type:    HiveAttack,
		Message: msg,
		State:   state,
	}

	protocol.eventChannel <- event
}

// gameFinishedResponse sends an event indicating that the game has finished.
func (protocol *CommunicationProtocol) GameFinishedResponse(msg string, state GameState) {
	event := Event{
		Type:    GameFinished,
		Message: msg,
		State:   state,
	}

	protocol.eventChannel <- event
}
