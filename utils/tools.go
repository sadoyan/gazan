package utils

import (
	"errors"
	"math/rand"
)

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func RetRandomMap(key string) (string, error) {

	if len(Dconf.Upstreams[key]) >= 1 {
		randomIndex := rand.Intn(len(Dconf.Upstreams[key]))
		pick := Dconf.Upstreams[key][randomIndex]
		return pick, nil

	} else {
		return "", errors.New("upstream not found, or upstreams list is empty")
	}

}

func RetRandom(input []string) string {
	randomIndex := rand.Intn(len(input))
	pick := input[randomIndex]
	return pick
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
