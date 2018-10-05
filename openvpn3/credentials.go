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

package openvpn3

// #include <library.h>
import "C"

// UserCredentials represents the credentials structure
type UserCredentials struct {
	Username string
	Password string
}

func (credentials *UserCredentials) toPtr() (cCredentials C.user_credentials, unregister func()) {
	cUsername := newCharPointer(credentials.Username)
	cPassword := newCharPointer(credentials.Password)

	cCredentials = C.user_credentials{
		username: cUsername.Ptr,
		password: cPassword.Ptr,
	}
	unregister = func() {
		cUsername.delete()
		cPassword.delete()
	}
	return
}
