# 🤑 Chime AI 🤖

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

The importer parses your Chime statements into a SQLite database for analysis. For details on setting up and running the importer, see the [README in the `importer` folder](./importer/README.md).

---

## Running the App

This app allows you to interact with your imported Chime transactions using AI. Once you've set up your `transactions.db` using the importer:

1. **Start the AI Interface** (Coming Soon):
   ```bash
   ./chime-ai
   ```

2. **Ask Questions About Your Transactions**:
    - "How much did I spend on dining last month?"
    - "What are my recurring subscriptions?"
    - "Show me my largest transactions in the past year."

---

## Roadmap

- [x] Import Chime bank statements into SQLite
- [x] Clean up parsed text files
- [x] Implement DB repository functions
- [ ] Add AI chat functionality for transaction queries
- [ ] Add tool calling and open ended query building agent
- [ ] Improve error handling and data validation
- [ ] Support for additional bank formats

---

## License

This project is licensed under the MIT License. See the LICENSE file for details.

