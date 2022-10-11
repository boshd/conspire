package coherey

import (
	"encoding/json"
	"net/http"
)

type CohereHandler interface {
	HandleGenerateConspiracyTheories(w http.ResponseWriter, r *http.Request) error
}

type cohereHandler struct {
	cohereService CohereService
}

func NewHandler(cohereService CohereService) CohereHandler {
	return &cohereHandler{
		cohereService: cohereService,
	}
}

func (handler *cohereHandler) HandleGenerateConspiracyTheories(w http.ResponseWriter, r *http.Request) error {

	var (
		callbackRequest struct {
			Topic  string `json:"topic"`
			Length uint   `json:"length"`
		}
	)

	if err := json.NewDecoder(r.Body).Decode(&callbackRequest); err != nil {
		http.Error(w, "bad request :/", 400)
		return nil
	}

	theories, err := handler.cohereService.GenerateConspiracyTheories(callbackRequest.Topic, callbackRequest.Length)
	if err != nil {
		http.Error(w, "something went wrong", 500)
		return err
	}

	w.Write([]byte(theories))
	return nil
}
