package rc_protocol

import "sync"

var once sync.Once
var proto protocol

type RCProtocol interface {
	Message(json string) Message
	Response(json string) Response
	Auth
}

type protocol struct {
	auth Auth
}

func (p protocol) CreateSig(name string, keyDir string) (string, error) {
	return p.auth.CreateSig(name, keyDir)
}

func (p protocol) CheckSig(header string, certDir string) (bool, error) {
	return p.auth.CheckSig(header, certDir)
}

func (p protocol) GetHeaderName() string {
	return p.auth.GetHeaderName()
}

func (p protocol) ParseHeader(header string) []string {
	return p.auth.ParseHeader(header)
}

func (p protocol) Message(json string) Message  {
	return newMessage(json)
}

func (p protocol) Response(json string) Response  {
	return newResponse(json)
}

func NewRCProtocol() RCProtocol {
	once.Do(func () {
		proto = protocol{}
		proto.auth = newAuth()
	})

	return proto
}
