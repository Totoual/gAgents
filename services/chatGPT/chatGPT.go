package chatgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	gAgents "github.com/totoual/gAgents/agent"
	pb "github.com/totoual/gAgents/protos/registry"
)

const (
	gptEndpoint = "https://api.openai.com/v1/chat/completions"
)

type GPTClient struct {
	ApiKey          string
	model           string
	eventDispatcher *gAgents.EventDispatcher
	topics          []string
}

func NewGPTClient(ed *gAgents.EventDispatcher, model string, se, te gAgents.EventType, t []string) *GPTClient {
	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists {
		log.Panicf("Please provide an api key!")
	}
	client := &GPTClient{
		ApiKey:          apiKey,
		model:           model,
		eventDispatcher: ed,
		topics:          t,
	}
	client.eventDispatcher.Subscribe(se, client.HandleSearchEvent)
	client.eventDispatcher.Subscribe(te, client.HandleSuggestionEvent)

	return client
}

func (client *GPTClient) GenerateText(prompt string) (string, error) {
	// Create a messages array with a system instruction followed by a user prompt
	messages := []map[string]interface{}{
		{
			"role": "system",
			"content": `You are an AI that strictly conforms to responses in JSON formatted strings. 
			 Your responses consist of valid JSON syntax, with no other comments, explainations, reasoninng, or dialogue not consisting of valid JSON.`,
		},
		{
			"role":    "user",
			"content": prompt,
		},
	}

	// Update the request payload to use 'messages' instead of 'prompt'
	requestPayload := map[string]interface{}{
		"model":    client.model,
		"messages": messages,
		// "max_tokens": 150,  // max_tokens is not typically used with the messages parameter
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

	fmt.Println(responsePayload)

	// Update the way you extract the response text to account for the 'choices' structure
	if choices, ok := responsePayload["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					fmt.Printf("The message content is : %s", content)
					return content, nil
				}
			}
		}
	}

	return "", errors.New("failed to extract text from response")
}

func (client *GPTClient) HandleSuggestionEvent(event gAgents.Event) {
	fmt.Println("HandleSuggestionEvent invoked")
	payload, ok := event.Payload.(*pb.AgentRegistration)
	if !ok {
		log.Println("Failed to convert eventPayload to *pb.AgentRegistration type")
		return
	}

	suggestionPrompt := fmt.Sprintf(`
	Capabilities:  %s
	Tags: %s
	Existing Topics: %s

	Please provide only the list of the most relevant topics for the agent to s
	ubscribe to from the list of existing topics provided, based on the capabilities and tags. Use "relevant_topics" as key for the JSON
	`, payload.Capabilities, payload.Tags, client.topics)

	response, err := client.GenerateText(suggestionPrompt)
	if err != nil {
		fmt.Printf("Error when trying to get data from GPT")
	}

	var result map[string][]string
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The choices we get are: \n %s", response)
	event.ResponseChan <- result["relevant_topics"]
}

func (client *GPTClient) HandleSearchEvent(event gAgents.Event) {
	fmt.Printf("Invoke Search Event!")
	payload, ok := event.Payload.(*pb.SearchMessage)
	if !ok {
		log.Println("Failed to convert event.Payload to *pb.SearchMessage")
		return
	}
	searchPrompt := fmt.Sprintf(`
	Description: %s
	Existing Topics: %s
	Please analyze the description to identify and extract any information regarding the object, 
	characteristics, category, price range, intended use, deadline or time frame, preferred brands/makers, 
	material preferences, additional notes, and location preferences. 
	Then, suggest the most relevant topics to broadcast the message. 
	Provide the information in a JSON format with the specified schema:

	Schema:
		{
		"object": "string",
		"characteristics": [
			"string"
		],
		"category": "string",
		"price_range": float,
		"intended_use": "string",
		"material_preferences": ["string"],
		"relevant_topics": ["string"]
		}

	`, payload.Description, client.topics)

	response, err := client.GenerateText(searchPrompt)
	if err != nil {
		fmt.Printf("Error when trying to get data from GPT")
	}

	var sr SearchResult
	err = json.Unmarshal([]byte(response), &sr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The choices we get are: \n %s", response)
	event.ResponseChan <- sr
}
