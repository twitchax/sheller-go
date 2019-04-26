package sheller

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
	"sync"
)

// ExecutionScanner wraps a bufio.Scanner to stream lines of text.
type ExecutionScanner struct {
	Command       *exec.Cmd
	Stdout        string
	Stderr        string
	StdoutHandler func(string)
	StderrHandler func(string)
	StdoutChannel chan<- string
	StderrChannel chan<- string

	wg     *sync.WaitGroup
	stdout strings.Builder
	stderr strings.Builder
}

// Start reading lines from the underlying Buffer.
func (es *ExecutionScanner) Start() {
	es.wg = &sync.WaitGroup{}

	stdoutPipe, err := es.Command.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderrPipe, err := es.Command.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	stdoutScanner := bufio.NewScanner(stdoutPipe)
	stderrScanner := bufio.NewScanner(stderrPipe)

	go func() {
		es.wg.Add(1)
		defer es.wg.Done()
		for stdoutScanner.Scan() {
			text := stdoutScanner.Text()

			es.stdout.WriteString(text)
			text = trim(text)

			if es.StdoutHandler != nil {
				es.StdoutHandler(text)
			}
			if es.StdoutChannel != nil {
				es.StdoutChannel <- text
			}
		}
	}()

	go func() {
		es.wg.Add(1)
		defer es.wg.Done()
		for stderrScanner.Scan() {
			text := stderrScanner.Text()

			es.stderr.WriteString(text)
			text = trim(text)

			if es.StderrHandler != nil {
				es.StderrHandler(text)
			}
			if es.StderrChannel != nil {
				es.StderrChannel <- text
			}
		}
	}()
}

// Wait the completion of the stdout and stderr buffers.
func (es *ExecutionScanner) Wait() {
	es.wg.Wait()

	es.Stdout = es.stdout.String()
	es.Stderr = es.stderr.String()
}
