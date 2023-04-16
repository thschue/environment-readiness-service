package server

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"html/template"
	"net/http"
	"time"
)

//go:embed index.html
var indexTemplate string

//go:embed assets/*
var assets embed.FS

type Config struct {
	Port string `json:"port"`
}

type InfraReadiness struct {
	Ready        bool                   `json:"infraReady"`
	Applications map[string]Application `json:"applications"`
}

var infra = InfraReadiness{
	Ready: false,
}

type Application struct {
	Name      string              `json:"name"`
	Workloads map[string]Workload `json:"workloads"`
	Version   string              `json:"version"`
	Ready     bool                `json:"ready"`
}

type Workload struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Ready   bool   `json:"ready"`
}

type StatusPost struct {
	Ready bool `json:"ready"`
}

type AppRegistration struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Ready   bool   `json:"ready"`
}

type WorkloadRegistration struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Application string `json:"application"`
	Ready       bool   `json:"ready"`
}

type RegistrationPost struct {
	Type     string               `json:"type"`
	App      AppRegistration      `json:"app"`
	Workload WorkloadRegistration `json:"workload"`
}

type EventPost struct {
	Timestamp       string `json:"timestamp"`
	App             string `json:"app"`
	Workload        string `json:"workload"`
	AppVersion      string `json:"appVersion"`
	WorkloadVersion string `json:"workloadVersion"`
	Result          string `json:"result"`
}

var events = []EventPost{}

func (s *Config) healthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func (s *Config) registrationHandler(w http.ResponseWriter, r *http.Request) {
	if infra.Applications == nil {
		infra.Applications = make(map[string]Application)
	}
	var p RegistrationPost
	err := json.NewDecoder(r.Body).Decode(&p)
	fmt.Println(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch p.Type {
	case "app":
		for k, v := range infra.Applications {
			fmt.Println(k, p.App.Name)
			if k == p.App.Name {
				v.Version = p.App.Version
				v.Ready = p.App.Ready
				infra.Applications[k] = v
				return
			}
		}
		infra.Applications[p.App.Name] = Application{
			Name:    p.App.Name,
			Version: p.App.Version,
			Ready:   p.App.Ready,
		}
	case "workload":
		for k, v := range infra.Applications {
			if k == p.Workload.Application {
				if v.Workloads == nil {
					v.Workloads = make(map[string]Workload)
				}
				v.Workloads[p.Workload.Name] = Workload{
					Name:    p.Workload.Name,
					Version: p.Workload.Version,
					Ready:   p.Workload.Ready,
				}
				return
			}
		}
		fmt.Println("Creating new app")
		infra.Applications[p.Workload.Application] = Application{
			Name: p.Workload.Application,
			Workloads: map[string]Workload{
				p.Workload.Name: {
					Name:    p.Workload.Name,
					Version: p.Workload.Version,
					Ready:   p.Workload.Ready,
				},
			},
		}
	}
}

func (s *Config) readinessHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		js, err := json.MarshalIndent(infra, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	case "POST":
		var p StatusPost
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		infra.Ready = p.Ready
	}
}

func (s *Config) eventHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		js, err := json.MarshalIndent(events, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	case "POST":
		var p EventPost
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p.Timestamp = time.Now().Format(time.RFC3339)
		events = append(events, p)
	}
}

func (s *Config) webHandler(w http.ResponseWriter, r *http.Request) {
	var status string
	if infra.Ready {
		status = "Ready"
	} else {
		status = "Not Ready"
	}

	data := struct {
		InfrastructureState string
		WelcomeMessage      string
	}{
		InfrastructureState: status,
		WelcomeMessage:      "Hello KubeCon",
	}

	tpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (s *Config) Serve() error {
	assetsFs := http.FileServer(http.FS(assets))

	mux := http.NewServeMux()

	mux.HandleFunc("/infraReadiness", s.readinessHandler)
	mux.HandleFunc("/register", s.registrationHandler)
	mux.HandleFunc("/healthz", s.healthzHandler)
	mux.HandleFunc("/event", s.eventHandler)
	mux.HandleFunc("/web", s.webHandler)
	mux.Handle("/assets/", assetsFs)
	color.Green("Starting server on port %s", s.Port)
	err := http.ListenAndServe(":"+s.Port, mux)
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		return err
	}
	return nil
}
