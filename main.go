package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
    "sync"
    "time"
)

var (
    mu             sync.RWMutex
    ethereumNodes  string
    pollingDuration = 10 * time.Second
)

func updateEthereumNodes(addressRecord string) {
    var buffer bytes.Buffer
    ipAddresses := []string{"127.0.0.1"} // Example IP addresses

    for _, ipAddress := range ipAddresses {
        resp, err := http.Get(fmt.Sprintf("http://%s:8080/enode", ipAddress))
        if err != nil {
            log.Printf("Error retrieving enode address from %s: %s", ipAddress, err)
            continue
        }
        defer resp.Body.Close()

        contents, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Printf("Error parsing response from %s: %s", ipAddress, err)
            continue
        }

        enodeAddress := strings.TrimSpace(string(contents))
        if buffer.Len() > 0 {
            buffer.WriteString(",")
        }
        buffer.WriteString(enodeAddress)
        log.Printf("%s with enode address %s", ipAddress, enodeAddress)
    }

    // Update list
    mu.Lock()
    defer mu.Unlock()
    ethereumNodes = buffer.String()
}
	
func startPollUpdateEthereumNodes(addressRecord string) {
    for {
        go updateEthereumNodes(addressRecord)
        <-time.After(pollingDuration)
    }
}

func main() {
    fmt.Println("Hello, World!")
    go startPollUpdateEthereumNodes("someAddressRecord")
    http.HandleFunc("/", webHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func webHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("handling request from %s", r.RemoteAddr)
    mu.RLock()
    defer mu.RUnlock()
    fmt.Fprintln(w, ethereumNodes)
}

// Removed the non-declaration statement outside function body
