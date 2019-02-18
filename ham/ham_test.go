package ham

import "testing"

func TestHam(t *testing.T) {
	var machineState, incomingState, currentState, incomingValue, currentValue float64
	h, err := NewHam(machineState, incomingState, currentState, incomingValue, currentValue)

	if err != nil {
		t.Error(err)
	}

	_ = h
}
