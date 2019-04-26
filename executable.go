package sheller

import (
	"log"
	"os"
	"os/exec"
	"time"
)

// EnvironmentVariableMap is a map[string]string which is meant to hold environment variables.
type EnvironmentVariableMap map[string]string

// ArgumentList is a []string that is meant to hold the list of arguments.
type ArgumentList []string

// Executable defines the type for a generic executable.
type Executable struct {
	exe                   string
	environmentVariables  EnvironmentVariableMap
	arguments             ArgumentList
	standardOutputHandler func(string)
	standardErrorHandler  func(string)
	standardOutputChannel chan string
	standardErrorChannel  chan string
}

// CreateExecutable a new Executable instance.
func CreateExecutable(exe string) *Executable {
	return &Executable{
		exe:                  exe,
		environmentVariables: environmentToMap(os.Environ()),
		//arguments:            make(ArgumentList, 0),
	}
}

// Clone the execution context represented by this Executable.
func (e *Executable) Clone() *Executable {
	return &Executable{
		exe:                  e.exe,
		environmentVariables: e.environmentVariables.clone(),
		arguments:            e.arguments.clone(),
	}
}

// Execute synchronously executes the command defined by this execution context.
func (e *Executable) Execute() (commandResult *CommandResult, err error) {
	command := exec.Command(e.exe, e.arguments...)
	command.Env = e.environmentVariables.toCommandEnv()

	es := ExecutionScanner{
		Command:       command,
		StdoutHandler: e.standardOutputHandler,
		StderrHandler: e.standardErrorHandler,
		StdoutChannel: e.standardOutputChannel,
		StderrChannel: e.standardErrorChannel,
	}
	es.Start()

	startTime := time.Now()

	if err := command.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); !ok {
			log.Fatal(exiterr)
		}
	}

	endTime := time.Now()

	es.Wait()

	if e.standardOutputChannel != nil {
		close(e.standardOutputChannel)
	}
	if e.standardErrorChannel != nil {
		close(e.standardErrorChannel)
	}

	result := &CommandResult{
		Succeeded:      command.ProcessState.Success(),
		ExitCode:       command.ProcessState.ExitCode(),
		StandardOutput: es.Stdout,
		StandardError:  es.Stderr,
		StartTime:      startTime,
		EndTime:        endTime,
		RunTime:        endTime.Sub(startTime),
	}

	return result, nil
}

// ExecuteAsync executes the command defined by this execution context and allows a done channel.
func (e *Executable) ExecuteAsync(ch chan<- *CommandResult) {
	result, err := e.Execute()
	if err != nil {
		log.Fatal(err)
	}

	ch <- result
}

// WithEnvironmentVariable adds an environment variable to the execution context.  This method may be called multiple times.
func (e *Executable) WithEnvironmentVariable(key string, value string) *Executable {
	newExe := e.Clone()
	newExe.environmentVariables[key] = value

	return newExe
}

// WithArguments adds an argument to the execution context.  This method may be called multiple times.
func (e *Executable) WithArguments(args ...string) *Executable { return e.WithArgument(args...) }

// WithArgument adds an argument to the execution context.  This method may be called multiple times.
func (e *Executable) WithArgument(args ...string) *Executable {
	newExe := e.Clone()
	newExe.arguments = newExe.arguments.appendMany(args...)

	return newExe
}

// UseStandardOutputHandler sets the standard output handler on the execution context.  This method
// may be called once, and subsequent invocations will overwrite the previous value in the new context.
func (e *Executable) UseStandardOutputHandler(handler func(string)) *Executable {
	newExe := e.Clone()
	newExe.standardOutputHandler = handler

	return newExe
}

// UseStandardErrorHandler sets the standard output handler on the execution context.  This method
// may be called once, and subsequent invocations will overwrite the previous value in the new context.
func (e *Executable) UseStandardErrorHandler(handler func(string)) *Executable {
	newExe := e.Clone()
	newExe.standardErrorHandler = handler

	return newExe
}

// UseStandardOutputChannel sets the standard output channel on the execution context.  This method
// may be called once, and subsequent invocations will overwrite the previous value in the new context.
func (e *Executable) UseStandardOutputChannel(ch chan string) *Executable {
	newExe := e.Clone()
	newExe.standardOutputChannel = ch

	return newExe
}

// UseStandardErrorChannel sets the standard output channel on the execution context.  This method
// may be called once, and subsequent invocations will overwrite the previous value in the new context.
func (e *Executable) UseStandardErrorChannel(ch chan string) *Executable {
	newExe := e.Clone()
	newExe.standardErrorChannel = ch

	return newExe
}
