package openvpn3

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestSessionInitFailsForInvalidProfile(t *testing.T) {
	session := NewSession(&fmtLogger{})
	session.Start("abc", Credentials{})
	err := session.Wait()
	assert.Equal(t, ErrInitFailed, err)
}

func TestSessionConnectFailsForInvalidRemote(t *testing.T) {
	session := NewSession(&fmtLogger{})
	session.Start("remote localhost 1111", Credentials{})
	err := session.Wait()
	assert.Equal(t, ErrConnectFailed, err)

}

type fmtLogger struct {
}

func (l *fmtLogger) Log(text string) {
	fmt.Println(text)
}
