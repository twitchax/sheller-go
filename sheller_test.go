package sheller_test

import (
	"testing"

	"strings"

	sheller "../sheller-go"
	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	expected := "lol"

	echoResult, _ := sheller.
		UseExecutable("/bin/echo").
		WithArgument(expected).
		Execute()

	assert.True(t, echoResult.Succeeded, "Expected Succeeded.")
	assert.Equal(t, expected, trim(echoResult.StandardOutput), "Unexpected output.")
}

func TestHandlers(t *testing.T) {
	expected := "lol"
	result := ""

	echoResult, _ := sheller.
		UseExecutable("/bin/echo").
		WithArgument(expected).
		UseStandardOutputHandler(func(s string) { result += s }).
		Execute()

	assert.True(t, echoResult.Succeeded, "Expected Succeeded.")
	assert.Equal(t, expected, result, "Unexpected output.")
}

func TestChannelOutput(t *testing.T) {
	expected := "lol"
	result := ""

	resultCh := make(chan string, 1)
	doneCh := make(chan bool, 1)

	go func(result *string, resultCh <-chan string, doneCh chan<- bool) {
		for {
			line, ok := <-resultCh
			*result += line
			if !ok {
				doneCh <- true
			}
		}
	}(&result, resultCh, doneCh)

	echoResult, _ := sheller.
		UseExecutable("/bin/echo").
		WithArgument(expected).
		UseStandardOutputChannel(resultCh).
		Execute()

	<-doneCh

	assert.True(t, echoResult.Succeeded, "Expected Succeeded.")
	assert.Equal(t, expected, result, "Unexpected output.")
}

func TestAsync(t *testing.T) {
	expected := "lol"

	ch := make(chan *sheller.CommandResult, 1)

	go sheller.
		UseExecutable("/bin/echo").
		WithArgument(expected).
		ExecuteAsync(ch)

	echoResult := <-ch

	assert.True(t, echoResult.Succeeded, "Expected Succeeded.")
	assert.Equal(t, expected, trim(echoResult.StandardOutput), "Unexpected output.")
}

func trim(s string) string {
	return strings.Trim(s, " \t\n")
}
