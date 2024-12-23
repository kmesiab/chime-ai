# 🤑 Chime AI 🤖

![Golang](https://img.shields.io/badge/Go-00add8.svg?labelColor=171e21&style=for-the-badge&logo=go)

![Build](https://github.com/kmesiab/chime-ai/actions/workflows/go-build.yml/badge.svg)
![Lint](https://github.com/kmesiab/chime-ai/actions/workflows/go-lint.yml/badge.svg)
![Test](https://github.com/kmesiab/chime-ai/actions/workflows/go-test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/kmesiab/chime-ai)](https://goreportcard.com/report/github.com/kmesiab/equilibria)

**Import your Chime bank statements in PDF format, then chat about them with AI.**

_Note: This is a hacky personal project. Don't rely on this for important things..._

---

## Usage

[![Chime Logo](https://upload.wikimedia.org/wikipedia/commons/thumb/f/f6/Chime_company_logo.svg/2880px-Chime_company_logo.svg.png)](https://www.chime.com)

### Download Your Statements

1. Log into your Chime account and navigate to your Documents:
   [https://app.chime.com/settings/documents](https://app.chime.com/settings/documents)
2. Click on the `PDF` link on the right.
3. Download **every statement you can!**

Place all your statements in a single folder and name them like:

```text
Your_Name_Checking_eStatement (1).pdf
Your_Name_Checking_eStatement (2).pdf
```

---

## Importing Your Statements

The importer parses your Chime statements into a SQLite database for
analysis. For details on setting up and running the importer,
see the [README in the `importer` folder](./importer/README.md).

---

## Running the App

This app allows you to interact with your imported Chime transactions
using AI. Once you've set up your `transactions.db` using the importer:

1. **Start the AI Interface** (Coming Soon):

   ```bash
   ./chime-ai
   ```

2. **Ask Questions About Your Transactions**:
    - "How much did I spend on dining last month?"
    - "What are my recurring subscriptions?"
    - "Show me my largest transactions in the past year."

### Example Response

```text
Prompt: Where do I spend most of my money?  Give me the top 10 places in October

Response:

Based on the provided transaction information for October, 
here are the top 10 places where you spent the most money:

1. **Transfer to Credit Builder**: $1,167.47
2. **Transfer to Chime Savings Account**: $1,166.32
3. **Foursight**: $631.02
4. **Transfer to Chase**: $525.00
5. **Progressive**: $361.38
6. **Clear Me**: $259.00
7. **Verizon**: $245.47
8. **Amazon**: $171.93
9. **Les Schwab Tire Centers**: $157.00
10. **Alaska Air**: $153.99

Your highest spending was on transfers and savings, which 
constitutes financial management actions rather than direct 
purchasing of goods or services from typical vendors or service 
providers; this is something to keep in mind when analyzing your 
expenses critically for budget adjustments or identifying areas 
where you could potentially reduce spending if necessary.
```

---

## Roadmap

- [x] Import Chime bank statements into SQLite
- [x] Clean up parsed text files
- [x] Implement DB repository functions
- [x] Add AI chat functionality for transaction queries
- [x] Add tool calling and open ended query building agent
- [ ] Improve error handling and data validation
- [ ] Support for additional bank formats
