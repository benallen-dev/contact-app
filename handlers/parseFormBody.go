package handlers

import (
	"strings"
)

func parseFormBody(body string) (out map[string][]string) {
	out = make(map[string][]string)

	parts := strings.Split(string(body), "&")

	for _, part := range parts {
		kv := strings.Split(part, "=")
		elem, ok := out[kv[0]]
		if !ok {
			out[kv[0]] = []string{kv[1]}
		} else {
			out[kv[0]] = append(elem, kv[1])
		}
	}

	return out
}
