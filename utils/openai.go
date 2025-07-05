package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"soulprint-backend/config"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client     *openai.Client
	httpClient *http.Client
	useLocal   bool
	localURL   string
	localModel string
}

// Local model request structures (for Ollama/local APIs)
type LocalModelRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type LocalModelResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func NewOpenAIClient() *OpenAIClient {
	client := &OpenAIClient{
		httpClient: &http.Client{},
		useLocal:   config.AppConfig.UseLocalModel,
		localURL:   config.AppConfig.LocalModelURL,
		localModel: config.AppConfig.LocalModelName,
	}

	// Debug logging
	log.Printf("DEBUG: UseLocalModel=%t, LocalURL=%s, LocalModel=%s",
		config.AppConfig.UseLocalModel, config.AppConfig.LocalModelURL, config.AppConfig.LocalModelName)

	if !config.AppConfig.UseLocalModel {
		client.client = openai.NewClient(config.AppConfig.OpenAIAPIKey)
	}

	return client
}

func (oai *OpenAIClient) GenerateReflection(journalContent, reflectionType string) (string, error) {
	log.Printf("DEBUG: GenerateReflection called with useLocal=%t", oai.useLocal)

	if oai.useLocal {
		log.Printf("DEBUG: Using local model for reflection")
		return oai.generateLocalReflection(journalContent, reflectionType)
	}

	log.Printf("DEBUG: Using OpenAI model for reflection")
	if config.AppConfig.OpenAIAPIKey == "" {
		return "AI reflection unavailable - API key not configured", nil
	}

	prompt := oai.buildPrompt(journalContent, reflectionType)

	resp, err := oai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: config.AppConfig.OpenAIModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a thoughtful journal reflection assistant. Provide insightful, empathetic, and constructive reflections on journal entries.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   500,
			Temperature: 0.7,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate reflection: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no reflection generated")
	}

	return resp.Choices[0].Message.Content, nil
}

func (oai *OpenAIClient) buildPrompt(content, reflectionType string) string {
	switch reflectionType {
	case "summary":
		return fmt.Sprintf("Please provide a concise summary of this journal entry, highlighting the main themes and emotions:\n\n%s", content)
	case "analysis":
		return fmt.Sprintf("Please provide a thoughtful analysis of this journal entry, identifying patterns, emotions, and potential insights for personal growth:\n\n%s", content)
	default: // "insight"
		return fmt.Sprintf("Please provide a thoughtful reflection on this journal entry, offering gentle insights and perspectives that might help with self-understanding and growth:\n\n%s", content)
	}
}

func (oai *OpenAIClient) ExtractKeywords(content string) ([]string, error) {
	if oai.useLocal {
		return oai.extractLocalKeywords(content)
	}

	if config.AppConfig.OpenAIAPIKey == "" {
		return []string{}, nil
	}

	prompt := fmt.Sprintf("Extract 3-5 key themes or keywords from this journal entry. Return only the keywords separated by commas:\n\n%s", content)

	resp, err := oai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: config.AppConfig.OpenAIModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   50,
			Temperature: 0.3,
		},
	)

	if err != nil {
		return []string{}, err
	}

	if len(resp.Choices) == 0 {
		return []string{}, nil
	}

	// Parse comma-separated keywords
	keywords := strings.Split(strings.TrimSpace(resp.Choices[0].Message.Content), ",")
	for i, keyword := range keywords {
		keywords[i] = strings.TrimSpace(keyword)
	}
	return keywords, nil
}

// Local model methods
func (oai *OpenAIClient) generateLocalReflection(journalContent, reflectionType string) (string, error) {
	systemPrompt := "You are a thoughtful journal reflection assistant. Provide insightful, empathetic, and constructive reflections on journal entries."
	userPrompt := oai.buildPrompt(journalContent, reflectionType)

	fullPrompt := fmt.Sprintf("%s\n\nUser: %s\n\nAssistant:", systemPrompt, userPrompt)

	return oai.callLocalModel(fullPrompt)
}

func (oai *OpenAIClient) extractLocalKeywords(content string) ([]string, error) {
	prompt := fmt.Sprintf("Extract 3-5 key themes or keywords from this journal entry. Return only the keywords separated by commas:\n\n%s\n\nKeywords:", content)

	response, err := oai.callLocalModel(prompt)
	if err != nil {
		return []string{}, err
	}

	// Parse comma-separated keywords
	keywords := strings.Split(strings.TrimSpace(response), ",")
	for i, keyword := range keywords {
		keywords[i] = strings.TrimSpace(keyword)
	}
	return keywords, nil
}

func (oai *OpenAIClient) callLocalModel(prompt string) (string, error) {
	// Support for Ollama API format
	reqBody := LocalModelRequest{
		Model:  oai.localModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Try Ollama API format first
	url := fmt.Sprintf("%s/api/generate", oai.localURL)
	resp, err := oai.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call local model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("local model API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var localResp LocalModelResponse
	if err := json.Unmarshal(body, &localResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return localResp.Response, nil
}
