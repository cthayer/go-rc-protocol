package rc_protocol

import (
    "encoding/json"
)

type Response interface {}

type response struct {
    Id string `json:"id"`
    Stdout string `json:"stdout"`
    Stderr string `json:"stderr"`
    ExitCode int `json:"exitCode"`
    Signal string `json:"signal"`
}

func newResponse(jsonStr string) Response {
    resp := response{ExitCode:-1}

    if err := json.Unmarshal([]byte(jsonStr), resp); err != nil {
        // log error
    }

    return resp
}
