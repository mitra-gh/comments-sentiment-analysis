package main

import (
	"comments-sentiment-analysis/internal/analyzer"
	"comments-sentiment-analysis/internal/csvparser"
	"comments-sentiment-analysis/internal/report"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize OpenAI client
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY not set in environment variables")
	}
	client := openai.NewClient(apiKey)

	// Parse CSV file
	comments, err := csvparser.ParseCSV("data/comments.csv")
	if err != nil {
		log.Fatalf("Error parsing CSV: %v", err)
	}

	// Prepare comments for analysis
	var texts []string
	for _, comment := range comments {
		texts = append(texts, comment.Text)
	}

	// Analyze comments (consider batching if you have many comments)
	fmt.Println("Analyzing comments...")
	results, err := analyzer.AnalyzeComments(client, texts)
	if err != nil {
		log.Fatalf("Error analyzing comments: %v", err)
	}

	// Generate report
	fmt.Println("Generating report...")
	err = report.GenerateReport(results, "outputs/report.json")
	if err != nil {
		log.Fatalf("Error generating report: %v", err)
	}

	fmt.Println("Analysis complete! Report saved to outputs/report.json")
}