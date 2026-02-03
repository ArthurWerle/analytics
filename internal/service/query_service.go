package service

import (
	"analytics/external"
	"analytics/internal/db"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

const SYSTEM_PROMPT_TO_GET_QUERY = `
	I want you to, based on a database schema and a user questions,
	give me the exact Postgres query to get the data the user wants.

	Based on the following database schema:

	------------------------------
	-- PostgreSQL database schema --
	------------------------------

	create type transaction_type as enum ('income', 'expense');;

	TABLE transactions (
		id SERIAL PRIMARY KEY,
		migrated_id INTEGER,
		is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
		category_id INTEGER,
		amount DECIMAL(12, 2) NOT NULL,
		type transaction_type NOT NULL,
		description TEXT,
		date TIMESTAMP WITH TIME ZONE,
		frequency VARCHAR(50),
		start_date DATE,
		end_date DATE,
		created_by_id INTEGER NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	-------------------------------------
    Categories reference

	ID  NAME                DESCRIPTION
	1	Food             	All expenses with food exluding food bought to daily meals. This is meant for food like deliveries and stuff	
	2	Salary             	Salary	
	3	Health             	All expenses related to health	
	4	Needs             	All needs: internet, energy, water etcetera	
	5	Subscriptions       All subscriptions	
	6	Grocery             Grocery shop: food, house items, cleaning items etc	
	7	Car             	All expenses related to my car	
	8	Useless             All useless things I buy lol	
	10	Gifts             	all gifts	
	11	Electronics         All electronics	
	12	Education           College, courses, material etc	
	13	Help             	All transactions related to helping people	
	14	Misc             	Miscellaneous	
	15	Clothes             All clothes	
	16	Books             	all books	
	17	Transportation      Uber, bus, taxis, trains etc	
	18	Electrodomestics    Eletrodomésticos	
	19	Games             	all games	
	20	Pets             	All expenses related to pets	
	21	Travel             	All expenses related to travelling	
	22	Housing             Rent and all other housing expenses	
	23	Furniture           All furniture 	
	--------------------------------------------------------------------------------

	--------------------------------------
	-- PostgreSQL database reference end --
	--------------------------------------

	Instructions:
		- The types are only meant to store 2 different types: 'income' or 'expense'.
		- Note that the "Salary" category is of type "income", so it should never be used for questions for expenses. Always ensure to be using
		transactions from type "expense" for it.

	Query Requirements:
		- Use standard SQL syntax without any special formatting or markdown
		- Do not include any comments or explanations in the query
		- Use proper table aliases for better readability
		- Include appropriate JOIN conditions
		- Use proper date/time functions when dealing with timestamps
		- Consider using CTEs (WITH clauses) for complex queries
		- Always include proper ORDER BY clauses when the order matters
		- Use appropriate aggregate functions (SUM, COUNT, etc.) when needed
		- When using CTEs with aggregations, ensure proper handling of grouped and non-grouped columns
		- Avoid cross joins between aggregated and non-aggregated results
		- Create a query that will most likely be able to answer the question. Remember that after running
		the query, I will send all results back to an AI to analyze them, so some calculations can be made
		by the AI itself and don't need to be made in the query.

	Very important!
		Answer the user question ONLY with the SQL query to get what he wants. Please, answer ONLY with the query that
		will need to run in the database, nothing more. 

		I'm gonna get your answer and run it directly in the database, so nothing more than the SQL itself can be present
		in the answer. I don't want errors like ERROR: syntax error at or near "` + "```" + `" (SQLSTATE 42601)

	Example of good response:
		SELECT t.date, SUM(t.amount) as total_amount
		FROM transactions t
		WHERE t.date >= CURRENT_DATE - INTERVAL '1 month'
		GROUP BY t.date
		ORDER BY t.date DESC;

	Example of bad response:
		"` + "```sql" + `
		SELECT * FROM transactions;
		-- This is a comment that shouldn't be here
		"` + "```" + `"
`

const SYSTEM_PROMPT_TO_ANALYZE_RESULTS = `
	I will give you a user question and the results of the database search to answer it. Using the
	results from the query, you'll need to answer the user in the most concise and precise way.

	Guidelines for your response:
	1. Start with a direct answer to the user's question
	2. If the results are numerical, include the exact numbers
	3. If there are multiple results, summarize the key findings
	4. If the results are empty, explain what that means in the context of the question
	5. Keep the response focused and relevant to the original question
	6. Use proper formatting for numbers and dates
	7. If the results show trends or patterns, point them out
	8. If the query returned an error, explain what went wrong
	9. If the result is a money value, the currency is BRL (R$)

	The following are the question and the results: 
`

type QueryService struct {
	query      string
	userPrompt string
}

func (q *QueryService) GetQuery(prompt string) string {
	openAIService := &external.OpenAIService{}
	openAIService.GetOpenAIClient()
	q.query = openAIService.Ask(prompt)

	log.Printf("Query: %v", q.query)

	return q.query
}

func (q *QueryService) RunQuery() (pgx.Rows, error) {
	databaseService := &db.DatabaseService{}
	pool := databaseService.GetPool()

	rows, err := pool.Query(context.Background(), q.query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (q *QueryService) ConvertResult(rows pgx.Rows) (gin.H, error) {
	log.Printf("Starting ConvertResult function")

	if rows == nil {
		log.Printf("Error: rows is nil")
		return nil, fmt.Errorf("rows is nil")
	}

	fieldDescriptions := rows.FieldDescriptions()
	log.Printf("Field descriptions count: %d", len(fieldDescriptions))

	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = string(fd.Name)
		log.Printf("Column %d: %s", i, columns[i])
	}

	var results []map[string]interface{}
	rowCount := 0

	hasNext := rows.Next()
	log.Printf("First rows.Next() returned: %v", hasNext)

	if hasNext {
		values, err := rows.Values()
		if err != nil {
			log.Printf("Error getting row values: %v", err)
			return nil, err
		}
		log.Printf("Raw values for row %d: %+v", rowCount, values)
		log.Printf("Row %d values count: %d", rowCount, len(values))

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			log.Printf("Row %d, Column %s, Raw value: %v, Type: %T", rowCount, col, val, val)
			row[col] = val
		}
		results = append(results, row)
		rowCount++

		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				log.Printf("Error getting row values: %v", err)
				return nil, err
			}
			log.Printf("Raw values for row %d: %+v", rowCount, values)
			log.Printf("Row %d values count: %d", rowCount, len(values))

			row := make(map[string]interface{})
			for i, col := range columns {
				val := values[i]
				log.Printf("Row %d, Column %s, Raw value: %v, Type: %T", rowCount, col, val, val)
				row[col] = val
			}
			results = append(results, row)
			rowCount++
		}
	} else {
		log.Printf("No rows returned from query")
		row := make(map[string]interface{})
		for _, col := range columns {
			row[col] = nil
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	log.Printf("Total rows processed: %d", rowCount)
	log.Printf("Results array length: %d", len(results))
	log.Printf("Results content: %+v", results)

	response := gin.H{
		"query":   q.query,
		"columns": columns,
		"results": results,
	}

	log.Printf("Final response structure: query=%v, columns=%v, results length=%d",
		q.query, columns, len(results))

	return response, nil
}

func (q *QueryService) AnalyzeDatabase(userPrompt string) (string, error) {
	q.userPrompt = userPrompt
	prompt := SYSTEM_PROMPT_TO_GET_QUERY + userPrompt
	q.GetQuery(prompt)

	rows, err := q.RunQuery()
	if err != nil {
		log.Printf("Error running query: %v", err)
		return "", err
	}

	results, err := q.ConvertResult(rows)
	if err != nil {
		log.Printf("Error converting results: %v", err)
		return "", err
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Printf("Error marshaling results to JSON: %v", err)
		return "", err
	}

	openAIService := &external.OpenAIService{}
	openAIService.GetOpenAIClient()
	resultsPrompt := SYSTEM_PROMPT_TO_ANALYZE_RESULTS + q.userPrompt + string(jsonData)
	log.Printf("Results prompt: %v", resultsPrompt)

	response := openAIService.Ask(resultsPrompt)

	log.Printf("AI response: %v", response)

	return response, nil
}
