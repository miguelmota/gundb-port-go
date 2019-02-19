package get

// Get ...
func Get(lex, graph map[string]interface{}) map[string]interface{} {
	if lex == nil {
		lex = make(map[string]interface{})
	}
	if graph == nil {
		graph = make(map[string]interface{})
	}
	soul := lex["#"].(string)
	node, ok := graph[soul]
	if !ok {
		return nil
	}
	key, ok := lex["."]
	if ok {
		tmp, ok := node.(map[string]interface{})[key.(string)]
		if !ok {
			return nil
		}
		// equiv: (node = {_: node._})[key] = tmp
		node = map[string]interface{}{
			"_": node.(map[string]interface{})["_"],
		}
		node.(map[string]interface{})[key.(string)] = tmp

		// equiv: tmp = node._['>']
		tmp = node.(map[string]interface{})["_"].(map[string]interface{})[">"]

		// equiv: (node._['>'] = {})[key] = tmp[key]
		node.(map[string]interface{})["_"].(map[string]interface{})[">"] = make(map[string]interface{})

		node.(map[string]interface{})["_"].(map[string]interface{})[">"].(map[string]interface{})[key.(string)] = tmp.(map[string]interface{})[key.(string)]
	}

	ack := make(map[string]interface{})
	ack[soul] = node
	return ack
}
