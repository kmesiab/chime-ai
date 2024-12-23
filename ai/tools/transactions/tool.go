package transactions

import (
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const ToolName = "TransactionsTool"
const ToolDescription = `Given the user's question, construct a sqlite query to retrieve a dataset to make
an informed response
		
		The table schema is:
		create table transactions
				(
					id          integer primary key,
					date        datetime,
					description text,
					type        text,
					amount      real,
					net_amount  real,
					settle_date datetime
				);

		Categories are:
	    Transfer
		Purchase
		Direct Debit
		Fee
		ATM Withdrawal
		Deposit
		Round Up

		Sample rows:
		8,2024-07-19 00:00:00+00:00,Islandadv.Whalewatch,Purchase,-274.18,-274.18,2024-07-20 00:00:00+00:00
		9,2024-07-19 00:00:00+00:00,Transfer from Chime Savings Account,Transfer,275,275,2024-07-19 00:00:00+00:00
		10,2024-07-19 00:00:00+00:00,"Supermaven, Inc.",Purchase,-10,-10,2024-07-20 00:00:00+00:00
		11,2024-07-19 00:00:00+00:00,"Notion Labs, Inc.",Purchase,-11.03,-11.03,2024-07-20 00:00:00+00:00
	
		Notes: 
		Descriptions can vary despite being the same merchant.  When constructing queries, consider
	    using flexible matching.
`

var toolParams = jsonschema.Definition{
	Type: jsonschema.Object,
	Properties: map[string]jsonschema.Definition{
		"sql": {
			Type:        jsonschema.String,
			Description: ToolDescription,
		},
	},
	Required: []string{"sql"},
}

var functionDefinition = openai.FunctionDefinition{
	Name:        ToolName,
	Description: ToolDescription,
	Strict:      false,
	Parameters:  toolParams,
}

func NewTool() openai.Tool {
	return openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &functionDefinition,
	}
}
