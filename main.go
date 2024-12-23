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

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a financial advisor and a SQL expert with access to a transaction
				history database via tools and can query it for more robust data and analysis. You use database results to make 
				informed responses to help the user.`,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Where do I spend most of my money?  Give me the top 10 places in October",
		},
	}

	completionRequest := openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: messages,
		Tools:    []openai.Tool{transactions.NewTool()},
	}

	resp, err := client.CreateChatCompletion(timeoutCtx, completionRequest)

	if err != nil {
		log.Printf("Error creating chat completion: %v\n", err)
		return
	}

	type ToolResponse struct {
		SQL string `json:"sql"`
	}

	topChoice := resp.Choices[0]
	if len(topChoice.Message.ToolCalls) > 0 {
		var toolResults = "Here is relevant financial information you requested:\n\n"
		for _, call := range topChoice.Message.ToolCalls {
			if call.Function.Name == "TransactionsTool" {

				var toolResponse ToolResponse
				if err = json.Unmarshal([]byte(call.Function.Arguments), &toolResponse); err != nil {
					log.Printf("Error unmarshalling JSON: %v\n", err)
					return
				}

				output, err := repository.ExecuteRawQuery(toolResponse.SQL)

				if err != nil {
					log.Printf("Error executing SQL query: %v\n", err)
					return
				}

				response, err := json.Marshal(output)

				if err != nil {
					log.Printf("Error marshaling JSON: %v\n", err)
					return
				}
				toolResults += string(response) + "\n\n"
			}
		}

		if strings.TrimSpace(toolResults) != "" {

			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: toolResults,
			})

			completionRequest := openai.ChatCompletionRequest{
				Model:            openai.GPT4o,
				Messages:         messages,
				Temperature:      0.98,
				FrequencyPenalty: 0.7,
			}

			resp, err := client.CreateChatCompletion(timeoutCtx, completionRequest)

			if err != nil {
				log.Printf("Error creating chat completion: %v\n", err)
				return
			}

			fmt.Println("Your financial analysis:")
			fmt.Println(resp.Choices[0].Message.Content)
		} else {
			fmt.Println("No data found for the given query.")
			fmt.Println(topChoice.Message.Content)
		}

	} else {
		fmt.Println("No data found for the given query. Please adjust your SQL query.")
		fmt.Println(topChoice.Message.Content)
	}
}
