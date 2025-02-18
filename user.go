/*
Copyright (c) 2014 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipmi

import "strings"

const MaxUsernameLen = 16

// GetUserNameRequest per section 22.29
type GetUserNameRequest struct {
	UserID byte
}

// GetUserNameRequest per section 22.29
type GetUserNameResponse struct {
	CompletionCode
	Username string
}

// SetUserNameRequest per section 22.29
type SetUserNameRequest struct {
	UserID   byte
	Username string
}

// SetUserNameRequest per section 22.29
type SetUserNameResponse struct {
	CompletionCode
}

type GetUserSummaryRequest struct {
	ChannelNumber byte
	UserID        byte
}

type GetUserSummaryResponse struct {
	CompletionCode
	MaxUsers         byte
	CurrEnabledUsers byte
	FixedNameUsers   byte
	ChannelAccess    byte
}

type SetUserPassRequest struct {
	UserID byte
	Pass   []byte
}

type SetUserPassResponse struct {
	CompletionCode
}

type EnableUserRequest struct {
	UserID byte
}

type EnableUserResponse struct {
	CompletionCode
}

type SetUserAccessRequest struct {
	AccessOptions    byte
	UserID           byte
	UserLimits       byte
	UserSessionLimit byte
}

type SetUserAccessResponse struct {
	CompletionCode
}

func (r *SetUserPassRequest) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 18)
	buf[0] = r.UserID
	buf[1] = 0x02
	copy(buf[2:], r.Pass)
	return buf, nil
}

func (r *SetUserPassRequest) UnmarshalBinary(buf []byte) error {
	if len(buf) < 18 {
		return ErrShortPacket
	}
	if len(buf) > 18 {
		return ErrLongPacket
	}
	r.UserID = buf[0]
	r.Pass = buf[2:]
	return nil
}

func (r *EnableUserRequest) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 2)
	buf[0] = r.UserID
	buf[1] = 0x01
	return buf, nil
}

func (r *EnableUserRequest) UnmarshalBinary(buf []byte) error {
	if len(buf) < 1 {
		return ErrShortPacket
	}
	if len(buf) > 1 {
		return ErrLongPacket
	}
	r.UserID = buf[0]
	return nil
}

func (r *GetUserNameResponse) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 1+MaxUsernameLen)
	buf[0] = byte(r.CompletionCode)
	copy(buf[1:], r.Username)
	return buf, nil
}

func (r *GetUserNameResponse) UnmarshalBinary(buf []byte) error {
	if len(buf) < 1+MaxUsernameLen {
		return ErrShortPacket
	}
	r.CompletionCode = CompletionCode(buf[0])
	r.Username = strings.Trim(string(buf[1:]), "\000")
	return nil
}

func (r *SetUserNameRequest) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 1+MaxUsernameLen)
	buf[0] = r.UserID
	copy(buf[1:], r.Username)
	return buf, nil
}

func (r *SetUserNameRequest) UnmarshalBinary(buf []byte) error {
	if len(buf) > 1+MaxUsernameLen {
		return ErrLongPacket
	}
	r.UserID = buf[0]
	r.Username = string(buf[1:])
	return nil
}
