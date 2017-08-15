package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pragkent/slackwork/sink"
)

func (s *Server) registerHandlers(r *mux.Router) {
	r.HandleFunc("/hooks/{name:[a-zA-Z0-9]+}/{secret:[a-zA-Z0-9]+}", s.HookHandler)
}

func (s *Server) HookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	hookName, hookSecret := vars["name"], vars["secret"]

	log.Printf("Hook requested: %v", hookName)

	hook := s.c.Slack.GetHook(hookName)
	if hook == nil {
		log.Printf("Hook not found: %v", hookName)
		http.Error(w, "Hook not found", 404)
		return
	}

	if !hook.Auth(hookSecret) {
		log.Printf("Hook secret invalid: %v %v", hookName, hookSecret)
		http.Error(w, "Hook access denied", 403)
		return
	}

	payload, err := s.parsePayLoad(r)
	if err != nil {
		log.Printf("Hook payload invalid %v: %v", hookName, err)
		http.Error(w, "Hook payload error", 400)
		return
	}

	if err := s.sc.Dispatch(hookName, payload); err != nil {
		http.Error(w, "Hook internal error", 500)
		return
	}
}

func (s *Server) parsePayLoad(r *http.Request) (*sink.Payload, error) {
	var err error
	var rawPayload []byte
	ct := r.Header.Get("Content-Type")
	if ct == "application/json" {
		rawPayload, err = s.getJsonPayload(r)
		if err != nil {
			return nil, err
		}
	} else {
		rawPayload = s.getFormPayload(r)
	}

	var payload sink.Payload
	err = json.Unmarshal(rawPayload, &payload)
	if err != nil {
		return nil, fmt.Errorf("Parse payload error: %v", err)
	}

	return &payload, nil
}

func (s *Server) getJsonPayload(r *http.Request) ([]byte, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Get json payload error: %v", err)
	}

	return b, nil
}

func (s *Server) getFormPayload(r *http.Request) []byte {
	r.ParseForm()
	return []byte(r.PostFormValue("payload"))
}
