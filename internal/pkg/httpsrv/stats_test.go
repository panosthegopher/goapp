package httpsrv

import (
	"testing"
)

func TestIncStats(t *testing.T) {
	server := &Server{}

	// Test my changes in 'incStats' for a new session
	server.incStats(sid1)
	if len(server.sessionStats) != 1 {
		t.Errorf("Expected 1 session and got %d", len(server.sessionStats))
	}
	if server.sessionStats[0].id != sid1 {
		t.Errorf("There is a mismatch in provided session's ID: %v", server.sessionStats[0].id)
	}
	if server.sessionStats[0].sent != 1 {
		t.Errorf("The message of session '%v' didn't sent", server.sessionStats[0].id)
	}

	// Test my changes in 'incStats' for 2 sessions
	server.incStats(sid2)
	if len(server.sessionStats) != 2 {
		t.Errorf("Expected 2 sessions and got %d", len(server.sessionStats))
	}
	if server.sessionStats[1].id != sid2 {
		t.Errorf("There is a mismatch in provided session's ID: %v", server.sessionStats[0].id)
	}
	if server.sessionStats[1].sent != 1 {
		t.Errorf("The message of session '%v' didn't sent", server.sessionStats[0].id)
	}
}
