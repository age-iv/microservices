package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func GetJSON(url string, target interface{}) error {
    resp, err := httpClient.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("bad status: %d", resp.StatusCode)
    }
    return json.NewDecoder(resp.Body).Decode(target)
}

func PostJSON(url string, body interface{}, target interface{}) error {
    jsonBody, _ := json.Marshal(body)
    resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonBody))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("bad status: %d", resp.StatusCode)
    }
    if target != nil {
        return json.NewDecoder(resp.Body).Decode(target)
    }
    return nil
}
