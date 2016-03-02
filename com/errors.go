package com
import "errors"

const (
	E_Success = 0
	E_Fail = 1
	E_Unknown = 2
)

const (
	E_AgentAccountNotExist = 1000
	E_AgentAccountExist = 1001
	E_AgentPasswordIncorrect = 1002
)

const (
	ErrClientTimeout = errors.New("client timeout")
)