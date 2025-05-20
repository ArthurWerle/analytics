package service

import (
	"context"
	"log"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const SYSTEM_PROMPT = `
	I want you to, based on a database schema and a user questions,
	give me the exact Postgres query to get the data the user wants.

	Based on the following database schema:

	------------------------------
	-- PostgreSQL database dump --
	------------------------------

	
	CREATE TABLE types (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);


	CREATE TABLE categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		color VARCHAR(50),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);


	CREATE TABLE transactions (
		id SERIAL PRIMARY KEY,
		category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
		amount DECIMAL(12, 2) NOT NULL,
		type_id INTEGER REFERENCES types(id) ON DELETE SET NULL,
		description TEXT,
		date TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE recurring_transactions (
		id SERIAL PRIMARY KEY,
		category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
		amount DECIMAL(12, 2) NOT NULL,
		type_id INTEGER REFERENCES types(id) ON DELETE SET NULL,
		description TEXT,
		frequency VARCHAR(50) NOT NULL,
		start_date DATE NOT NULL,
		end_date DATE,
		last_occurrence DATE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	

	--------------------------------------
	-- PostgreSQL database referece end --
	--------------------------------------

	Instrcutions:

		- 'recurring_transactions' is meant to store recurrent transactions. For example: subscriptions, monthly fees, yearly payments etc.
		- transactions is just for common transactions.
		- To get the total value spent for a month, you'll need to check both tables, because some values might be recurrent.
		- The types table is only meant to store 2 different types: 'income' or 'expense'.


	Very important!

		Answer the user question ONLY with the SQL query to get what he wants. Please, answer ONLY with the query that
		will need to run in the database, nothing more. 

		I'm gonna get your answer and run it directly in the database, so nothing more than the SQL itself can be present
		in the answer, not even things like '''sql.
`

func GetQuery(prompt string) string {
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Fatalf("OPENAI_API_KEY is not set")
	}

	client := openai.NewClient(
		option.WithAPIKey(openAIKey),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(SYSTEM_PROMPT),
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		panic(err.Error())
	}

	return chatCompletion.Choices[0].Message.Content
}
