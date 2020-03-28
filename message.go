package rc_protocol

import (
    "encoding/json"
)

type Message interface {}

type message struct {
    Id int `json:"id"`
    Command string `json:"command"`
    Options map[string]interface{} `json:"options"`
}

func newMessage(jsonStr string) Message {
    msg := message{}

    if err := json.Unmarshal([]byte(jsonStr), msg); err != nil {
        // log error
    }

    return msg
}
