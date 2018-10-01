/*
 * Copyright (C) 2018 The "MysteriumNetwork/go-openvpn" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package openvpn

import (
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testProcessPrefix = "[process-test-log] "

// TestHelperProcess IS ESENTIAL FOR CMD MOCKING - DO NOT DELETE
func TestHelperProcess(t *testing.T) {
	RunTestExecCmd()
}

func TestWrapperStartReturnsErrorOnNoArgs(t *testing.T) {
	execTestHelper := NewExecCmdTestHelper("TestHelperProcess")
	execCommand := func(arg ...string) *exec.Cmd {
		cmd := execTestHelper.ExecCommand("openvpn", arg...)
		cmd.Args = nil
		return cmd
	}
	execTestHelper.AddExecResult("", "", 0, 10000, "openvpn")

	process := NewCmdWrapper(testProcessPrefix, execCommand)
	err := process.Start([]string{})
	assert.NotNil(t, err)
}

func TestWaitAndStopProcessDoesNotDeadLocks(t *testing.T) {
	execTestHelper := NewExecCmdTestHelper("TestHelperProcess")
	execCommand := func(arg ...string) *exec.Cmd {
		return execTestHelper.ExecCommand("openvpn", arg...)
	}
	execTestHelper.AddExecResult("", "", 0, 10000, "openvpn")

	process := NewCmdWrapper(testProcessPrefix, execCommand)
	processStarted := sync.WaitGroup{}
	processStarted.Add(1)

	processWaitExited := make(chan int, 1)
	processStopExited := make(chan int, 1)

	go func() {
		assert.NoError(t, process.Start([]string{}))
		processStarted.Done()
		process.Wait()
		processWaitExited <- 1
	}()
	processStarted.Wait()

	go func() {
		process.Stop()
		processStopExited <- 1
	}()

	select {
	case <-processWaitExited:
	case <-time.After(600 * time.Millisecond):
		assert.Fail(t, "CmdWrapper.Wait() didn't return in 600 miliseconds")
	}

	select {
	case <-processStopExited:
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "CmdWrapper.Stop() didn't return in 100 miliseconds")
	}
}

func TestWaitReturnsIfProcessDies(t *testing.T) {
	execTestHelper := NewExecCmdTestHelper("TestHelperProcess")
	execCommand := func(arg ...string) *exec.Cmd {
		return execTestHelper.ExecCommand("openvpn", arg...)
	}
	execTestHelper.AddExecResult("", "", 0, 100, "openvpn")

	process := NewCmdWrapper(testProcessPrefix, execCommand)
	processWaitExited := make(chan int, 1)

	go func() {
		process.Wait()
		processWaitExited <- 1
	}()

	assert.NoError(t, process.Start([]string{}))
	select {
	case <-processWaitExited:
	case <-time.After(3000 * time.Millisecond):
		assert.Fail(t, "CmdWrapper.Wait() didn't return on time")
	}
}
