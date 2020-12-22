/*
 * go-openvpn -- Go gettable library for wrapping Openvpn functionality in go way.
 *
 * Copyright (C) 2020 BlockDev AG.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License Version 3
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program in the COPYING file.
 * If not, see <http://www.gnu.org/licenses/>.
 */

// Command test util adapted from https://gist.github.com/kglee79/db8f0bf3eafe962e0feddac8451387da

package openvpn

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const execTestExitCodeKey = "EXEC_HELPER_EXIT_CODE"
const execTestStdOutputKey = "EXEC_HELPER_STDOUT"
const execTestStdErrorKey = "EXEC_HELPER_STDERR"
const execTestDelayKey = "EXEC_HELPER_DELAY_MILISECONDS"
const execTestArgsKey = "EXEC_HELPER_ARGS"

// ExecCmdTestResult represents the mocked cmd exec result
type ExecCmdTestResult struct {
	command  string
	exitCode int
	stdOut   string
	stdErr   string
	delay    int
}

// ExecCmdTestHelper provides a way to test code that uses exec.Command by providing a mockable function that replaces
// the real exec.Command function during the test.
// Usage:
//  - Create a test function starting with the prefix "Test" such as 'func TestHelperProcess(t *testing.T)'. This
//    function must contain a call to testutils.RunTestExecCmd().
//
//    func TestHelperProcess(t *testing.T) {
//      testutils.RunTestExecCmd()
//    }
//
//  - Create a ExecCmdTestHelper instance using NewExecCmdTestHelper and pass in the name of the test function created.
//  - For each command that you want to mock, call "AddExecResult" on the ExecCmdTestHelper instance.
//	- The code which calls exec.Command must use a variable which holds the "exec.Command" function as this variable
//    must be replaced in the test file with the ExecCmdTestHelper's ExecCommand function so that it can mock the result.
//    For example, in your code under test you should have a variable such as var myexec = exec.Command, then where you
//    would normally use exec.Command you would use myexec instead. In your test file you would set myexec to
//    the ExecCmdTestHelper's ExecuteCommand function.
type ExecCmdTestHelper struct {
	testResults        map[string][]ExecCmdTestResult
	testHelperFuncName string
}

// NewExecCmdTestHelper creates a new ExecCmdTestHelper instance which will run the test function with the name
// specified by testHelperFuncName when the command is executed in order to mock the command's response.
func NewExecCmdTestHelper(testHelperFuncName string) *ExecCmdTestHelper {
	return &ExecCmdTestHelper{
		testResults:        make(map[string][]ExecCmdTestResult),
		testHelperFuncName: testHelperFuncName,
	}
}

// AddExecResult adds a mock response for the command given where the command stdout will contain the output string and
// the process will exit with the exit code given.
func (e *ExecCmdTestHelper) AddExecResult(stdOut, stdErr string, exitCode int, delayMiliseconds int, command ...string) {
	fullCommand := strings.Join(command, " ")
	base64Command := base64.StdEncoding.EncodeToString([]byte(command[0]))

	result := ExecCmdTestResult{
		stdOut:   stdOut,
		stdErr:   stdErr,
		exitCode: exitCode,
		command:  fullCommand,
		delay:    delayMiliseconds,
	}

	if e.testResults[base64Command] == nil {
		e.testResults[base64Command] = make([]ExecCmdTestResult, 0)
	}

	e.testResults[base64Command] = append(e.testResults[base64Command], result)
}

// ExecCommand is the stub for the real exec.Command function. This is called in place of exec.Command.
// Ensure you set the var back to exec.Command once your test is complete.
func (e *ExecCmdTestHelper) ExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=" + e.testHelperFuncName, "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)

	fullCommand := command

	if len(args) > 0 {
		fullCommand = command + " " + strings.Join(args, " ")
	}

	base64Command := base64.StdEncoding.EncodeToString([]byte(strings.Fields(fullCommand)[0]))

	if len(e.testResults[base64Command]) == 0 {
		fmt.Println("No result was setup for command: ", fullCommand)
		return nil
	}

	// Retrieve next result
	mockResults := e.testResults[base64Command][0]

	// Remove current result so that next time it will use next result that was setup.  If no next result, re-use same result.
	if len(e.testResults[base64Command]) > 1 {
		e.testResults[base64Command] = e.testResults[base64Command][1:]
	}

	ar := execTestArgsKey + "=" + strings.Join(args, " ")
	stdout := execTestStdOutputKey + "=" + mockResults.stdOut
	stderr := execTestStdErrorKey + "=" + mockResults.stdErr
	exitCode := execTestExitCodeKey + "=" + strconv.FormatInt(int64(mockResults.exitCode), 10)
	delay := execTestDelayKey + "=" + strconv.FormatInt(int64(mockResults.delay), 10)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", stdout, stderr, exitCode, delay, ar}

	return cmd
}

// RunTestExecCmd will simulate the execution of a command by returning a mocked response which includes output to stdout, stderr
// and a specific exit code.
func RunTestExecCmd() {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	stdout := os.Getenv(execTestStdOutputKey)
	stderr := os.Getenv(execTestStdErrorKey)
	exitCode, err := strconv.ParseInt(os.Getenv(execTestExitCodeKey), 10, 64)

	if err != nil {
		os.Exit(1)
	}

	delay := os.Getenv(execTestDelayKey)
	delayMillis, err := strconv.ParseInt(delay, 10, 64)
	if err != nil {
		os.Exit(1)
	}

	if delayMillis != 0 {
		time.Sleep(time.Millisecond * time.Duration(delayMillis))
	}

	fmt.Fprintf(os.Stdout, stdout)
	fmt.Fprintf(os.Stderr, stderr)
	os.Exit(int(exitCode))
}

func startFakeOpenvpnManagement(address, port string, stop chan struct{}) {
	addr := fmt.Sprintf("%v:%v", address, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("start error", err)
	}
	conn.Write([]byte(">INFO:OpenVPN Management Interface Version 1 -- type 'help' for more info\n"))
	conn.Write([]byte(">STATE:1522855903,CONNECTING,,,,,,\n"))

	conn.Write([]byte(">STATE:1522855903,WAIT,,,,,,\n"))
	conn.Write([]byte(">STATE:1522855903,AUTH,,,,,,\n"))
	conn.Write([]byte(">STATE:1522855904,GET_CONFIG,,,,,,\n"))
	conn.Write([]byte(">STATE:1522855904,ASSIGN_IP,,10.8.0.133,,,,\n"))
	conn.Write([]byte(fmt.Sprintf(">STATE:1522855905,CONNECTED,SUCCESS,10.8.0.133,%v,%v,,\n", address, port)))

	for {
		select {
		case <-stop:
			conn.Write([]byte(">STATE:1522855911,EXITING,SIGTERM,,,,,\n"))
			conn.Close()
			return
		case <-time.After(100 * time.Millisecond):
			conn.Write([]byte(">BYTECOUNT:36987,32252\n"))
		}
	}
}

// RunTestExecOpenvpn will run a simulated openvpn management
func RunTestExecOpenvpn() {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	sigc := make(chan os.Signal, 1)
	stop := make(chan struct{})
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		stop <- struct{}{}
	}()

	args := strings.Fields(os.Getenv(execTestArgsKey))
	port := args[2]
	address := args[1]

	startFakeOpenvpnManagement(address, port, stop)

	os.Exit(int(0))
}
