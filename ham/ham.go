package ham

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"
)

// Ham ...
type Ham struct {
	Defer      bool
	Historical bool
	Converge   bool
	Incoming   bool
	Current    bool
	State      bool
}

// NewHam ...
func NewHam(
	machineState float64,
	incomingState float64,
	currentState float64,
	incomingValue interface{},
	currentValue interface{},
) (*Ham, error) {
	if machineState < incomingState {
		return &Ham{Defer: true}, nil
	}
	if incomingState < currentState {
		return &Ham{Historical: true}, nil
	}
	if currentState < incomingState {
		return &Ham{Converge: true, Incoming: true}, nil
	}
	var incomingVal string
	var currentVal string
	if incomingState == currentState {
		incomingVal = Lexical(incomingValue)
		currentVal = Lexical(currentValue)
		if incomingVal == currentVal {
			return &Ham{State: true}, nil
		}
		if incomingVal < currentVal {
			return &Ham{Converge: true, Current: true}, nil
		}
		if currentVal < incomingVal {
			return &Ham{Converge: true, Incoming: true}, nil
		}
	}

	return nil, fmt.Errorf("Invalid CRDT Data: %s to %s at %v to %v", incomingVal, currentVal, incomingState, currentState)
}

// Mix ...
func Mix(change map[string]interface{}, graph map[string]interface{}) map[string]interface{} {
	machine := time.Now()
	var diff map[string]interface{}

	for soul := range change {
		node := change[soul].(map[string]interface{})
		for key := range node {
			val := node[key]
			if key == "_" {
				continue
			}

			// equiv: state = node._['>'][key]
			state := node["_"].(map[string]interface{})[">"].(map[string]interface{})[key]

			// equiv: was = (graph[soul]||{_:{'>':{}}})._['>'][key] || -Infinity
			var was interface{}
			soulv, ok := graph[soul]

			if ok {
				was = soulv.(map[string]interface{})["_"].(map[string]interface{})[">"].(map[string]interface{})[key]
			}

			if was == nil {
				// 'infinity'
				was = float64(math.MinInt64)
			}

			// equiv: known = (graph[soul]||{})[key]
			var known interface{}
			graphsoul, ok := graph[soul]
			if ok {
				known = graphsoul.(map[string]interface{})[key]
			}

			hm, err := NewHam(
				float64(machine.Unix()),
				state.(float64),
				was.(float64),
				val,
				known,
			)

			if err != nil {
				log.Fatal(err)
			}
			if !hm.Incoming {
				if hm.Defer {
					fmt.Println("defer", key, val)
					// need to implement this
				}

				continue
			}

			if diff == nil {
				diff = make(map[string]interface{})
			}

			_, ok = diff[soul]
			if !ok {
				// equiv: graph[soul] = {_:{'#':soul, '>':{}}}
				diff[soul] = make(map[string]interface{})
				diff[soul].(map[string]interface{})["_"] = make(map[string]interface{})
				diff[soul].(map[string]interface{})["_"].(map[string]interface{})["#"] = soul
				diff[soul].(map[string]interface{})["_"].(map[string]interface{})[">"] = make(map[string]interface{})
			}

			_, ok = graph[soul]
			if !ok {
				// equiv: graph[soul] = {_:{'#':soul, '>':{}}}
				graph[soul] = make(map[string]interface{})
				graph[soul].(map[string]interface{})["_"] = make(map[string]interface{})
				graph[soul].(map[string]interface{})["_"].(map[string]interface{})["#"] = soul
				graph[soul].(map[string]interface{})["_"].(map[string]interface{})[">"] = make(map[string]interface{})
			}

			graph[soul].(map[string]interface{})[key] = val
			diff[soul].(map[string]interface{})[key] = val

			diff[soul].(map[string]interface{})["_"].(map[string]interface{})[">"].(map[string]interface{})[key] = state
			graph[soul].(map[string]interface{})["_"].(map[string]interface{})[">"].(map[string]interface{})[key] = state
		}
	}

	return diff
}

// Lexical ...
func Lexical(value interface{}) string {
	js, err := json.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}

	return string(js)
}
