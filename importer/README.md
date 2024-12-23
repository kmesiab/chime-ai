# Importer

[![Chime Logo](https://upload.wikimedia.org/wikipedia/commons/thumb/f/f6/Chime_company_logo.svg/2880px-Chime_company_logo.svg.png)](https://www.chime.com)

## Download Your Statements

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

## Converting Your Statements

This program automatically converts your Chime PDF statements into
plain text using `pdftotext`. It then parses the data and imports
it into a SQLite database for analysis.

If `pdftotext` is not installed, the script will notify you.

### Installing `pdftotext`

For most Linux distributions:

```bash
sudo apt-get install poppler-utils
```

On macOS (via Homebrew):

```bash
brew install poppler
```

On Windows:

- Download and install Poppler from [https://blog.alivate.com.au/poppler-windows/](https://blog.alivate.com.au/poppler-windows/).<!-- markdownlint-disable-line MD013 -->
- Add the `bin` folder containing `pdftotext.exe` to your PATH.

## Building the Importer

**Install Go**:

Make sure you have Go installed. You can download it from [https://go.dev/dl/](https://go.dev/dl/).<!-- markdownlint-disable-line MD013 -->

**Clone the Repository**:

```bash
git clone <repository-url>
cd chime-ai/importer
```

**Build the Program**:

```bash
go build -o importer
```

---

## Running the Importer

Once built, you can run the program and point it to the directory
containing your Chime statements:

```bash
./importer -dir /path/to/your/statements
```

The program will:

1. Convert all PDF files in the specified directory into text files.
2. Parse the text files and extract transaction data.
3. Store the transactions in a SQLite database (`transactions.db`).
4. Clean up all generated `.txt` files once processing is complete.

---

## Output

- The transactions will be stored in a SQLite database named
`transactions.db`.
- Each transaction includes:
  - Date
  - Description
  - Transaction Type (e.g., Deposit, Withdrawal, Fee, etc.)
  - Amount
  - Net Amount
  - Settlement Date

---

## Example Usage

### Input Directory Structure

```plaintext
/path/to/your/statements/
├── Your_Name_Checking_eStatement (1).pdf
├── Your_Name_Checking_eStatement (2).pdf
├── Your_Name_Checking_eStatement (3).pdf
```

### Command

```bash
./importer -dir /path/to/your/statements
```

### Example Output

```plaintext
Found 3 PDF files for conversion.
Converted Your_Name_Checking_eStatement (1).pdf to Your_Name_Checking_eStatement (1).txt<!-- markdownlint-disable-line MD013 -->
Converted Your_Name_Checking_eStatement (2).pdf to Your_Name_Checking_eStatement (2).txt<!-- markdownlint-disable-line MD013 -->
Converted Your_Name_Checking_eStatement (3).pdf to Your_Name_Checking_eStatement (3).txt<!-- markdownlint-disable-line MD013 -->
Processing Your_Name_Checking_eStatement (1).txt...
Processing Your_Name_Checking_eStatement (2).txt...
Processing Your_Name_Checking_eStatement (3).txt...
Inserted 30 transactions from Your_Name_Checking_eStatement (1).txt
Inserted 28 transactions from Your_Name_Checking_eStatement (2).txt
Inserted 32 transactions from Your_Name_Checking_eStatement (3).txt
Deleted Your_Name_Checking_eStatement (1).txt
Deleted Your_Name_Checking_eStatement (2).txt
Deleted Your_Name_Checking_eStatement (3).txt
All files processed and cleaned up successfully!
```

---

## Notes

- Ensure that your PDF statements are formatted properly and contain structured text.
- The program will skip duplicate transactions to avoid redundant entries.
- The SQLite database can be queried using tools like DB Browser for SQLite or
programmatically with any library supporting SQLite.
