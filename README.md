# 🤑 Chime AI 🤖

**Import your Chime bank statements in PDF format, then chat 
about them with AI.**

_Note: This is a hacky personal project.  Don't rely on this for important 
things..._

## Usage

[![Chime Logo](https://upload.wikimedia.org/wikipedia/commons/thumb/f/f6/Chime_company_logo.svg/2880px-Chime_company_logo.svg.png)](https://www.chime.com)

### Download Your Statements

1. Log into your Chime account and navigate to your Documents:
https://app.chime.com/settings/documents
2. Click on the `PDF` link on the right
3. Download EVERY statement you can!

Place all your statements in a single folder, and name them like:

```text
Your_Name_Checking_eStatement (1).pdf
Your_Name_Checking_eStatement (2).pdf
```
## Converting Your Statements

We'll use the `pdftotext` command to convert your statements into a 
parsable text file.

1. Install `poppler` on your Mac:

```bash
brew install poppler
```

2. Convert your statements into a parsable text file:

```bash
pdftotext -layout Your_Name_Checking_eStatement (1).pdf checking_1.txt
```

### Conversion Script

In the `./bin` folder, you will find a convenience script called `convert.
sh`.  You can modify this file to point to the folder where you saved your 
statements and the naming convention you used.  

Make note of the section 
`for i in {1..19}; do`
and ensure you are looping through all of your statements.

## Importing Your Statements

Move all the output `.txt` files to `./importer/files` and execute;
Note: The output files should be in the format: `checking_n.txt`

## Running The App

```
cd ./importer 
go run main.go
```

This will create `transactions.db,` a tinysql database with
your chime transactions in it.
