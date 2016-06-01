package core

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type Result struct {
	StartDate float64
	EndDate   float64
}

func Run(command string, producers int, scripts int, workers int, bar func() bool) ([]Result, []error) {
	sem := make(chan bool, workers)
	results := make([]Result, producers*scripts)
	errors := make([]error, producers*scripts)
	processed := 1

	for i := 0; i < producers; i++ {
		for j := 0; j < scripts; j++ {
			sem <- true
			go func(producerIndex int, scriptIndex int) {
				defer func() { <-sem }()
				index := producerIndex*scripts + scriptIndex
				result, err := RunCommand(command)
				if err != nil {
					errors[index] = err
				} else {
					results[index] = *result
				}
				processed++
				if bar != nil {
					bar()
				}
			}(i, j)
		}
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return results, errors
}

func RunCommand(command string) (*Result, error) {
	parts := strings.Fields(command)
	head := parts[0]
	parts = parts[1:len(parts)]

	cmd := exec.Command(head, parts...)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(stdout, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
