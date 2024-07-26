#!/bin/bash

# Define the sources and their corresponding commands
declare -A sources
sources=(
    ["CNN"]="curl -s http://rss.cnn.com/rss/cnn_latest.rss | xmllint --nocdata --xpath '/rss/channel/item[1]/description/text()' -"
)

# Get a random key from the sources array
source_key=$(shuf -e "${!sources[@]}" -n 1)

# Fetch the text from the selected source
title=$(eval "${sources[$source_key]}")

# Prepare JSON output
json_output=$(jq -n --arg source "$source_key" --arg title "$title" \
    '{source: $source, text: $title}')
  
# Print the JSON output
echo "$json_output"