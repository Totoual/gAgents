package chatgpt

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const gptEndpoint = "https://api.openai.com/v1/engines/gpt-4/completions"

type GPTClient struct {
	ApiKey string
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
