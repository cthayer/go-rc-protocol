package rc_protocol

import (
    "encoding/json"
)

type Message struct {
    Id int `json:"id"`
    Command string `json:"command"`
    Options MessageOptions `json:"options"`
}

type MessageOptions struct {
    Cwd string `json:"cwd"`
    Timeout int `json:"timeout"`
    Env map[string]string `json:"env"`
}

func newMessage(jsonStr string) Message {
    msg := Message{}

    if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
        // log error
    }

    return msg
}
