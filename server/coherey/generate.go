package coherey

import (
	"fmt"

	"github.com/cohere-ai/cohere-go"
	"github.com/davecgh/go-spew/spew"
)

type CohereService interface {
	GenerateConspiracyTheories(topic string, length uint) (string, error)
}

type cohereService struct {
	cohere *cohere.Client
}

func NewService(cohere *cohere.Client) CohereService {
	return &cohereService{
		cohere: cohere,
	}
}

func (service *cohereService) GenerateConspiracyTheories(topic string, length uint) (string, error) {

	defaultPrompt := "This is an idea for a conspiracy theory about " + topic + ":"

	response, err := service.cohere.Generate(cohere.GenerateOptions{
		Prompt:    defaultPrompt,
		MaxTokens: length,
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	spew.Dump(response.Generations)

	return response.Generations[0].Text, nil

}
