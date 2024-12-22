#!/bin/bash

# Download your Chime bank statements as .pdf files into a single
# folder so they are named sequentially.  Such as:
# Jane_Doe_Checking_eStatement (1).pdf
# Jane_Doe_Checking_eStatement (2).pdf
#...
# Jane_Doe_Checking_eStatement (19).pdf
#
# Update the file name below to match yours and run the script.

# Directory where the PDFs are located
input_dir="."  # Replace with the path to your directory if not the current one

# Loop through 1 to 19
for i in {1..19}; do
    input_pdf="$input_dir/Jane_Doe_Checking_eStatement ($i).pdf"
    output_txt="$input_dir/checking_$i.txt"

    # Check if the PDF file exists
    if [[ -f "$input_pdf" ]]; then
        echo "Processing $input_pdf..."
        pdftotext -layout "$input_pdf" "$output_txt"
        echo "Output saved to $output_txt"
    else
        echo "File $input_pdf not found, skipping..."
    fi
done

echo "Processing completed."
