package main

import (
    "encoding/json"
)

const kvStoreKey = "disabled_users"

func (p *Plugin) loadDisabledUsers() map[string]bool {
    data, appErr := p.API.KVGet(kvStoreKey)
    if appErr != nil || data == nil {
        return make(map[string]bool)
    }

    var result map[string]bool
    if err := json.Unmarshal(data, &result); err != nil {
        return make(map[string]bool)
    }
    return result
}

func (p *Plugin) saveDisabledUsers(disabled map[string]bool) {
    data, err := json.Marshal(disabled)
    if err != nil {
        return
    }
    p.API.KVSet(kvStoreKey, data)
}

func (p *Plugin) isDisabled(userID string) bool {
    disabled := p.loadDisabledUsers()
    return disabled[userID]
}

func (p *Plugin) setDisabled(userID string, value bool) {
    disabled := p.loadDisabledUsers()
    if value {
        disabled[userID] = true
    } else {
        delete(disabled, userID)
    }
    p.saveDisabledUsers(disabled)
}
