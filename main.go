package main

import (
	"fmt"
	"net/http"
    "encoding/json"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Hello from Cisco Shipped Moooooo!</h1>\n")
    
    persons, err := runTestDB();
    if err != nil {
        http.Error(w, "Damn!. . .", 500)
        return
    }
    w.Header.Set('Content-type', 'application/json');
    err = json.NewEncoder(w).Encode(persons);
    if err != nil {
        http.Error(w, "Damn that was close!. . .", 500)
    }
}

func main() {
    http.HandleFunc("/", defaultHandler)
    fmt.Println("Example app listening at http://localhost:8888")
    http.ListenAndServe(":8888", nil)
}
