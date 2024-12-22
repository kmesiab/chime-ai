# 🤑 Chime AI 🤖

Import your Chime bank statements in PDF format, then chat 
about them with AI.

Note: This is a hacky personal project.  Don't rely on this for important 
things...

## Usage

### Convert The Statements

Download your chime PDF bank statements into
a single folder, so they are named like:

```text
Jane_Doe_Checking_eStatement (1).pdf
Jane_Doe_Checking_eStatement (2).pdf
...
Jane_Doe_Checking_eStatement (21).pdf
```

Use the script in `./bin/convert.sh` to convert all
your statements into a parsable text file.

Move all the output `.txt` files to `./importer/files` and 
execute `go run main.go`

This will create `transactions.db` tinysql database with
your chime transactions in it.
