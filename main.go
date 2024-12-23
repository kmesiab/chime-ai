package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/sashabaranov/go-openai"

	"github.com/kmesiab/chime-ai/ai/tools/transactions"
	"github.com/kmesiab/chime-ai/database"
)

var memory []openai.ChatCompletionMessage

func main() {

	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var APIKEY = os.Getenv("OPENAI_API_KEY")
	if APIKEY == "" {
		log.Fatal("OPENAI_API_KEY environment variable not set")
	}

	var (
		sqlDB *sql.DB
		db    *gorm.DB
		err   error
	)

	// Get the database connection
	if db, err = database.GetDBConnection(); err != nil {
		log.Printf("Error connecting to database: %v\n", err)
		return
	}

	// Clean up the database connection
	if sqlDB, err = db.DB(); err != nil {
		log.Printf("Error getting generic database object: %v\n", err)
		return
	} else {
		defer sqlDB.Close()
	}

	repository := database.NewTransactionRepository(db)

	client := openai.NewClient(APIKEY)

	memory = []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a financial advisor and a SQL expert with access to a transaction
				history database via tools and can query it for more robust data and analysis. You use database results to make 
				informed responses to help the user.`,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "How has my spending changed month over month and give me a summary",
		},
	}

	completionRequest := openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: memory,
		Tools:    []openai.Tool{transactions.NewTool()},
	}

	resp, err := client.CreateChatCompletion(timeoutCtx, completionRequest)

	if err != nil {
		log.Printf("Error creating chat completion: %v\n", err)
		return
	}

	topChoice := resp.Choices[0]

	if len(topChoice.Message.ToolCalls) > 0 {

		analysis, err := processAnalysis(ctx, topChoice, client, repository)

		if err != nil || analysis == "" {
			fmt.Printf("An error occurred while processing the analysis:\n%v\n", err)
		}

		fmt.Println("Final response:")
		fmt.Println(analysis)
		fmt.Println("\n--------\nAnalysis:")

	} else {

		fmt.Println("No data required to make this analysis:")
		fmt.Println(topChoice.Message.Content)
	}
}

type ToolResponse struct {
	SQL string `json:"sql"`
}

func processAnalysis(ctx context.Context, topChoice openai.ChatCompletionChoice, client *openai.Client, repository *database.TransactionRepository) (string, error) {

	var (
		err                error
		completionResponse openai.ChatCompletionResponse
	)

	toolOutput, err := processToolCalls(topChoice, repository)

	if err != nil {
		return "", err
	}

	if strings.TrimSpace(toolOutput) != "" {

		memory = append(memory, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: toolOutput + " is there any further information you need to query for?",
		})

		completionRequest := openai.ChatCompletionRequest{
			Model:            openai.GPT4oMini,
			Messages:         memory,
			Temperature:      0.7,
			FrequencyPenalty: 0.7,
			Tools:            []openai.Tool{transactions.NewTool()},
		}

		if completionResponse, err = client.CreateChatCompletion(ctx, completionRequest); err != nil {

			return "", fmt.Errorf("Error creating chat completion: %v\n", err)
		}

		// If there is another tool call, recursively process it
		if completionResponse.Choices[0].FinishReason == "tool_call" {

			return processToolCalls(completionResponse.Choices[0], repository)
		}

	}

	return completionResponse.Choices[0].Message.Content, nil
}

func processToolCalls(topChoice openai.ChatCompletionChoice, repository *database.TransactionRepository) (string, error) {

	var (
		err           error
		output        []map[string]interface{}
		queryResponse []byte
		toolResults   = "The financial information you requested:\n\n"
	)

	for _, call := range topChoice.Message.ToolCalls {
		if call.Function.Name == "TransactionsTool" {

			var toolResponse ToolResponse
			if err = json.Unmarshal([]byte(call.Function.Arguments), &toolResponse); err != nil {
				return "", fmt.Errorf("Invalid tool arguments: %v\n", err)
			}

			fmt.Printf("Executing SQL query: %s\n", toolResponse.SQL)

			if output, err = repository.ExecuteRawQuery(toolResponse.SQL); err != nil {
				return "", fmt.Errorf("Error executing SQL query: %v\n", err)
			}

			if queryResponse, err = json.MarshalIndent(output, "", "   "); err != nil {

				return "", fmt.Errorf("Error marshaling query results: %v\n", err)
			}

			toolResults += string(queryResponse) + "\n\n"
		}
	}

	return toolResults, nil
}
