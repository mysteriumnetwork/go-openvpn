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
	"testing"
	"time"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/tunnel"

	"github.com/mysteriumnetwork/go-openvpn/openvpn/config"
	"github.com/stretchr/testify/assert"
)

// TestHelperProcess_Openvpn IS ESENTIAL FOR CMD MOCKING - DO NOT DELETE
func TestHelperProcess_Openvpn(t *testing.T) {
	RunTestExecOpenvpn()
}

func TestOpenvpnProcessStartsAndStopsSuccessfully(t *testing.T) {
	execTestHelper := NewExecCmdTestHelper("TestHelperProcess_Openvpn")
	execCommand := func(arg ...string) *exec.Cmd {
		return execTestHelper.ExecCommand("openvpn", arg...)
	}
	execTestHelper.AddExecResult("", "", 0, 0, "openvpn")
	process := newProcess(&tunnel.NoopSetup{}, &config.GenericConfig{}, execCommand)

	err := process.Start()
	assert.NoError(t, err)

	time.Sleep(200 * time.Millisecond)

	process.Stop()

	err = process.Wait()
	assert.NoError(t, err)
}

func TestOpenvpnProcessStartReportsErrorIfCmdWrapperDiesTooEarly(t *testing.T) {
	execTestHelper := NewExecCmdTestHelper("TestHelperProcess")
	execTestHelper.AddExecResult("", "", 1, 0, "openvpn")
	execCommand := func(arg ...string) *exec.Cmd {
		return execTestHelper.ExecCommand("openvpn", arg...)
	}
	process := newProcess(&tunnel.NoopSetup{}, &config.GenericConfig{}, execCommand)

	err := process.Start()
	assert.Error(t, err)
}
