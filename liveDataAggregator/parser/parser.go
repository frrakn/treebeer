/*
======
Parser
======
Provides translation from raw TCP bytes to system-friendly format

*/

package parser

import (
	"encoding/json"
)

type Parser struct{}

func (p *Parser) Parse(b []byte) map[string]int {
	received := make(map[string]string)
	err := json.Unmarshal(b, &received)
	if err != nil {
		fmt.Sprintf("Parser: Unable to parse bytes into map[string]string - %s", err)
		return nil
	}
	result := make(map[string]int)
	for key, value := range received {
		result[key], err = strconv.Atoi(value)
		if err != nil {
			fmt.Sprintf("Parser: Unable to parse int from string %s", value)
			return nil
		}
	}
	return result
}

func (p *Parser) Unparse(m map[string]int) []byte {
	result, err := json.Marshal(m)
	if err != nil {
		fmt.Sprintf("Parser: Unable to marshal into bytes %+v", m)
		return nil
	}
	return result
}
