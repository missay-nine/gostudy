package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

const (
	defaultPort = 8080
)

//go:embed web
var content embed.FS

func main() {
	fsys, err := fs.Sub(content, "web")
	if err != nil {
		log.Fatalf("failed to load embedded assets: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(fsys)))

	addr := fmt.Sprintf("0.0.0.0:%d", defaultPort)
	log.Printf("Gomoku UI available at http://localhost:%d\n", defaultPort)
	log.Fatal(http.ListenAndServe(addr, mux))
}
