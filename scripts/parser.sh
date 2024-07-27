#!/bin/bash

# Define the sources and their corresponding commands
declare -A sources
sources=(
    ["CNN"]="curl -s http://rss.cnn.com/rss/cnn_latest.rss | xmllint --nocdata --xpath 'string(/rss/channel/item[1]/description)' -; curl -s http://rss.cnn.com/rss/cnn_latest.rss | xmllint --nocdata --xpath 'string(/rss/channel/item[1]/link)' -"
    ["BBC"]="curl -s https://feeds.bbci.co.uk/news/rss.xml | xmllint --nocdata --xpath 'string(/rss/channel/item[1]/description)' -; curl -s https://feeds.bbci.co.uk/news/rss.xml | xmllint --nocdata --xpath 'string(/rss/channel/item[1]/link)' -"
    ["Reuters"]="curl -s https://ir.thomsonreuters.com/rss/news-releases.xml | xmllint --nocdata --xpath 'string(/rss/channel/item[1]/description)' -; curl -s https://ir.thomsonreuters.com/rss/news-releases.xml | xmllint --nocdata --xpath 'string(/rss/channel/item[1]/link)' -"
)

# Get a random key from the sources array
source_key=$(shuf -e "${!sources[@]}" -n 1)

# Separate the commands for text and link
text_command=$(echo "${sources[$source_key]}" | cut -d';' -f1)
link_command=$(echo "${sources[$source_key]}" | cut -d';' -f2)

# Execute the commands to get the text and link
text=$(eval "$text_command")
link=$(eval "$link_command")

# Check if the text or link is empty and handle errors
if [ -z "$text" ] || [ -z "$link" ]; then
    echo "Error: Unable to fetch data from $source_key"
    exit 1
fi

# Prepare JSON output
json_output=$(jq -n --arg source "$source_key" --arg text "$text" --arg link "$link" \
    '{source: $source, text: $text, link: $link}')

# Print the JSON output
echo "$json_output"
