//go:generate sh -c "wget https://github.com/dyne/slangroom-exec/releases/latest/download/slangroom-exec-$(uname)-$(uname -m) -O ./slangroom-exec && chmod +x ./slangroom-exec"

package slangroom

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	b64 "encoding/base64"

	"github.com/amenzhinsky/go-memexec"
)

type SlangResult struct {
	Output string
	Logs   string
}

var binaryPath = "./slangroom-exec"

func loadBinary(binaryPath string) ([]byte, error) {

	binary, err := os.ReadFile(binaryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read binary file: %v", err)
	}
	return binary, nil
}

func SlangroomExec(conf string, contract string, data string, keys string, extra string, context string) (SlangResult, error) {

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return SlangResult{}, fmt.Errorf("binary %s does not exist. Please build it using 'go generate'", binaryPath)
	}

	binary, err := loadBinary(binaryPath)
	if err != nil {
		return SlangResult{}, err
	}

	exec, err := memexec.New(binary)
	if err != nil {
		return SlangResult{}, fmt.Errorf("failed to load Slangroom executable from memory: %v", err)
	}
	defer exec.Close()

	execCmd := exec.Command("slangroom-exec")

	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return SlangResult{}, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := execCmd.StderrPipe()
	if err != nil {
		return SlangResult{}, fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	stdin, err := execCmd.StdinPipe()
	if err != nil {
		return SlangResult{}, fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	inputs := []string{conf, contract, keys, data, extra, context}
	for _, input := range inputs {
		b64Input := b64.StdEncoding.EncodeToString([]byte(input))
		fmt.Fprintln(stdin, b64Input)
	}

	stdin.Close()

	err = execCmd.Start()
	if err != nil {
		return SlangResult{}, fmt.Errorf("failed to start command: %v", err)
	}

	stdoutOutput := make(chan string)
	stderrOutput := make(chan string)
	go captureOutput(stdout, stdoutOutput)
	go captureOutput(stderr, stderrOutput)

	err = execCmd.Wait()

	stdoutStr := <-stdoutOutput
	stderrStr := <-stderrOutput

	return SlangResult{Output: stdoutStr, Logs: stderrStr}, err
}

func captureOutput(pipe io.ReadCloser, output chan<- string) {
	defer close(output)

	buf := new(strings.Builder)
	_, err := io.Copy(buf, pipe)
	if err != nil {
		log.Printf("Failed to capture output: %v", err)
		return
	}
	output <- buf.String()
}
