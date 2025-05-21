package analyzer

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Sentiment string

const (
	Positive Sentiment = "positive"
	Negative Sentiment = "negative"
	Neutral  Sentiment = "neutral"
)

type AnalysisResult struct {
	Text      string
	Sentiment Sentiment
	MainTopic string
}

func AnalyzeComments(client *openai.Client, comments []string) ([]AnalysisResult, error) {
	var results []AnalysisResult

	for _, comment := range comments {
		if strings.TrimSpace(comment) == "" {
			continue
		}

		// Analyze sentiment
		sentiment, err := analyzeSentiment(client, comment)
		if err != nil {
			return nil, err
		}

		// Extract main topic
		topic, err := extractMainTopic(client, comment)
		if err != nil {
			return nil, err
		}

		results = append(results, AnalysisResult{
			Text:      comment,
			Sentiment: sentiment,
			MainTopic: topic,
		})
	}

	return results, nil
}

func analyzeSentiment(client *openai.Client, text string) (Sentiment, error) {
	prompt := fmt.Sprintf(`تعیین لحن متن زیر به فارسی:
"%s"

لطفاً فقط یکی از گزینه‌های زیر را انتخاب کنید:
- positive
- negative
- neutral`, text)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0,
		},
	)

	if err != nil {
		return "", err
	}

	response := strings.ToLower(strings.TrimSpace(resp.Choices[0].Message.Content))
	switch response {
	case "positive":
		return Positive, nil
	case "negative":
		return Negative, nil
	default:
		return Neutral, nil
	}
}

func extractMainTopic(client *openai.Client, text string) (string, error) {
	prompt := fmt.Sprintf(`موضوع اصلی این نظر فارسی را در یک کلمه یا عبارت کوتاه خلاصه کنید:
"%s"`, text)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.2,
		},
	)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}