package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type request struct { Text string `json:"text"` }
type response struct {
	Words int `json:"words"`
	Sentences int `json:"sentences"`
	Variety int `json:"variety"`
	Repeated []string `json:"repeated"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		engine := "unavailable"
		if _,err := os.Stat(analyzerPath()); err == nil { engine = "ready" }
		writeJSON(w, http.StatusOK, map[string]string{"status":"ok","engine":engine})
	})
	mux.HandleFunc("POST /api/analyze", analyze)
	mux.HandleFunc("GET /", func(w http.ResponseWriter,r *http.Request) {
		path := r.URL.Path
		if path == "/" { path = "/index.html" }
		if path != "/index.html" && path != "/app.css" && path != "/app.js" && path != "/assets/writing-desk-night.png" { http.NotFound(w,r); return }
		http.ServeFile(w,r,"."+path)
	})
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	server := &http.Server{Addr:":"+port, Handler:securityHeaders(mux), ReadHeaderTimeout:5*time.Second, ReadTimeout:10*time.Second, WriteTimeout:10*time.Second}
	log.Printf("Fluent Lab listening on http://localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}

func analyze(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w,r.Body,32<<10)
	defer r.Body.Close()
	var input request
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil { http.Error(w,"Invalid JSON or input too large",http.StatusBadRequest); return }
	if strings.TrimSpace(input.Text)=="" { http.Error(w,"Text is required",http.StatusBadRequest); return }
	ctx,cancel := context.WithTimeout(r.Context(),3*time.Second)
	defer cancel()
    cmd := exec.CommandContext(ctx,analyzerPath())
	cmd.Stdin = strings.NewReader(input.Text)
	data,err := cmd.Output()
	if err != nil {
		if errors.Is(ctx.Err(),context.DeadlineExceeded) { http.Error(w,"Analysis timed out",http.StatusGatewayTimeout); return }
		log.Printf("analyzer: %v",err)
		http.Error(w,"Analyzer unavailable; run cargo build --release in analyzer",http.StatusServiceUnavailable); return
	}
	var result response
	if err := json.Unmarshal(data,&result); err != nil { http.Error(w,"Analyzer returned invalid data",http.StatusBadGateway); return }
	writeJSON(w,http.StatusOK,result)
}

func writeJSON(w http.ResponseWriter,status int,value any) { w.Header().Set("Content-Type","application/json; charset=utf-8"); w.WriteHeader(status); _=json.NewEncoder(w).Encode(value) }
func analyzerPath() string { bin := filepath.Join("analyzer","target","release","fluent_analyzer"); if runtime.GOOS=="windows" { bin += ".exe" }; return bin }
func securityHeaders(next http.Handler) http.Handler { return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) { w.Header().Set("X-Content-Type-Options","nosniff"); w.Header().Set("Referrer-Policy","no-referrer"); w.Header().Set("X-Frame-Options","DENY"); if r.Method==http.MethodPost && r.URL.Path=="/api/analyze" && !strings.HasPrefix(r.Header.Get("Content-Type"),"application/json") { http.Error(w,"Content-Type must be application/json",http.StatusUnsupportedMediaType); return }; next.ServeHTTP(w,r) }) }

