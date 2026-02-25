package config

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func StartRestServer(listenAddr string, port int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		templateList := make([]string, 0)
		for name := range templates {
			templateList = append(templateList, name)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(templateList)
	})

	addr := fmt.Sprintf("%s:%d", listenAddr, port)
	http.ListenAndServe(addr, mux)
}
