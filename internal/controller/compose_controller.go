package controller

import (
	"encoding/json"
	"net/http"

	"github.com/fpswan/anycraft-backend/internal/model"
	"github.com/fpswan/anycraft-backend/internal/service"
)

type ComposeController struct {
	svc *service.ComposeService
}

func NewComposeController(svc *service.ComposeService) *ComposeController {
	return &ComposeController{svc: svc}
}

func (c *ComposeController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/compose/base-elements", c.handleBaseElements)
	mux.HandleFunc("/api/v1/compose/combine", c.handleCombine)
	mux.HandleFunc("/api/v1/compose/challenges", c.handleChallenges)
}

func (c *ComposeController) handleBaseElements(w http.ResponseWriter, r *http.Request) {
	gameCode := r.URL.Query().Get("game_code")
	if gameCode == "" {
		http.Error(w, "missing game_code", http.StatusBadRequest)
		return
	}
	list, err := c.svc.GetBaseElements(r.Context(), gameCode)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, model.BaseElementsResponse{OK: true, Items: list})
}

func (c *ComposeController) handleCombine(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var req model.CombineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	resp := c.svc.Combine(r.Context(), req)
	writeJSON(w, resp)
}

func (c *ComposeController) handleChallenges(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	var req model.ChallengesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	resp := c.svc.GetChallenges(r.Context(), req)
	writeJSON(w, resp)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	_ = enc.Encode(v)
}
