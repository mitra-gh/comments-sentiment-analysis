package report

import (
	"encoding/json"
	"os"

	"comments-sentiment-analysis/internal/analyzer"
)

type SummaryReport struct {
	TotalComments   int `json:"total_comments"`
	PositiveCount  int `json:"positive_count"`
	NegativeCount  int `json:"negative_count"`
	NeutralCount   int `json:"neutral_count"`
	PositivePercentage float64 `json:"positive_percentage"`
	NegativePercentage float64 `json:"negative_percentage"`
	NeutralPercentage  float64 `json:"neutral_percentage"`
	Topics         map[string]TopicStats `json:"topics"`
}

type TopicStats struct {
	Count          int `json:"count"`
	PositiveCount int `json:"positive_count"`
	NegativeCount int `json:"negative_count"`
	NeutralCount  int `json:"neutral_count"`
}

func GenerateReport(results []analyzer.AnalysisResult, outputPath string) error {
	report := SummaryReport{
		Topics: make(map[string]TopicStats),
	}

	report.TotalComments = len(results)

	for _, result := range results {
		// Update sentiment counts
		switch result.Sentiment {
		case analyzer.Positive:
			report.PositiveCount++
		case analyzer.Negative:
			report.NegativeCount++
		case analyzer.Neutral:
			report.NeutralCount++
		}

		// Update topic stats
		topic := result.MainTopic
		if topic == "" {
			topic = "موضوع نامشخص"
		}

		stats := report.Topics[topic]
		stats.Count++

		switch result.Sentiment {
		case analyzer.Positive:
			stats.PositiveCount++
		case analyzer.Negative:
			stats.NegativeCount++
		case analyzer.Neutral:
			stats.NeutralCount++
		}

		report.Topics[topic] = stats
	}

	// Calculate percentages
	if report.TotalComments > 0 {
		report.PositivePercentage = float64(report.PositiveCount) / float64(report.TotalComments) * 100
		report.NegativePercentage = float64(report.NegativeCount) / float64(report.TotalComments) * 100
		report.NeutralPercentage = float64(report.NeutralCount) / float64(report.TotalComments) * 100
	}

	// Save to JSON file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}