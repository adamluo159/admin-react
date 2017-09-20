package yada

import "testing"

func TestAgentServer(t *testing.T) {
	y := New("./config.json")
	y.Run()
}
