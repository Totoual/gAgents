package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	gAgents "github.com/totoual/gAgents/agent"
	pb "github.com/totoual/gAgents/protos/registry"
)

const (
	gptEndpoint = "https://api.openai.com/v1/engines/gpt-4/completions"
)

type GPTClient struct {
	ApiKey          string
	eventDispatcher *gAgents.EventDispatcher
}

func NewGPTClient(ed *gAgents.EventDispatcher, se, te gAgents.EventType) *GPTClient {
	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists {
		log.Panicf("Please provide an api key!")
	}
	client := &GPTClient{
		ApiKey:          apiKey,
		eventDispatcher: ed,
	}
	client.eventDispatcher.Subscribe(se, client.HandleSearchEvent)
	client.eventDispatcher.Subscribe(te, client.HandleSuggestionEvent)

	return client
}

func (client *GPTClient) GenerateText(prompt string) (string, error) {
	requestPayload := map[string]interface{}{
		"prompt":     prompt,
		"max_tokens": 150,
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", gptEndpoint, bytes.NewBuffer(jsonData))

	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+client.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var responsePayload map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&responsePayload)
	if err != nil {
		return "", err
	}
	return responsePayload["choices"].([]interface{})[0].(map[string]interface{})["text"].(string), nil
}

func (client *GPTClient) HandleSuggestionEvent(event gAgents.Event) {
	fmt.Println("HandleSuggestionEvent invoked")
	payload, ok := event.Payload.(*pb.AgentRegistration)
	if !ok {
		log.Println("Failed to convert eventPayload to *pb.AgentRegistration type")
		return
	}
	fmt.Printf("Event payload is: %v \n", payload)
	fmt.Printf("Capabilities are: %v \n", payload.Capabilities)
	fmt.Printf("Tags are: %v \n", payload.Tags)
	// capabilities := "negotiation, retail"
	// tags := "clothing, modern-cloths, trousers, men, women, kids"
	// existingTopics := `["Homeware", "KidsToys", "ModernCloths", "Cloths", "ClothsLondon", "KidsItemsEurope"]`

	// suggestionPrompt := fmt.Sprintf(`
	// 	Capabilities:  %s
	// 	Location: London
	// 	Tags: %s
	// 	Existing Topics: %s

	// 	Please analyse the details of the agent and extract any information on the capabilities, location, and tags. Then suggest the most relevant topics for the agent to subscribe to from the list of existing topics provided.
	// `, capabilities, tags, existingTopics)

}

func (client *GPTClient) HandleSearchEvent(event gAgents.Event) {
	//TBA
}
