package main

import (
	"bufio"
	"errors"
	"os"
)

type EnvVars []string

func LoadEnvVars(fn string) (ev EnvVars, err error) {
	ev = make(EnvVars, 0)

	file, err := os.Open(fn) // #nosec G304
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// missing file should not be an error, but rather
			// return an empty environment
			err = nil
		}

		return
	}

	defer file.Close() // #nosec G307

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ev = append(ev, scanner.Text())
	}

	err = scanner.Err()

	return
}
