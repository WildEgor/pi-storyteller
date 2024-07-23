#!/bin/bash

STEPS=5

SOURCES_COUNT=2
TITLE_9GAG=$(curl -s http://9gagrss.com/feed/ | xmllint --nocdata --xpath '/rss/channel/item[1]/title/text()' -)
TITLE_REDDIT=$(curl -s 'https://www.reddit.com/r/nottheonion.rss' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/115.0' -H 'Accept: text/html,applcation/xml;q=0.9,image/avif,image/webp,*/*;q=0.8' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'Upgrade-Insecure-Requests: 1' -H 'Sec-Fetch-Dest: document' -H 'Sec-Fetch-Mode: navigate' -H 'Sec-Fetch-Site: none' -H 'Sec-Fetch-User: ?1' -H 'Connection: keep-alive' -H 'TE: trailers' | sed 's/<title>/\n<title>/g' | sed 's/<\/title>/<\/title>\n/g' | grep -Po '<title>.*</title>' | sed 's/<title>//g' | sed 's/<\/title>//g' | sed -n '2p')

set -e

ARR[0]=$TITLE_9GAG
ARR[1]=$TITLE_REDDIT
RAND_TITLE=$[ $RANDOM % $SOURCES_COUNT ]

PREFIX_ARR[0]=""
PREFIX_ARR[1]="Realistic style scene with "
PREFIX_ARR[2]="Anime scene with "
PREFIX_ARR[3]="Gothic painting of "
PREFIX_ARR[4]="Documentary style photography of "
PREFIX_ARR[5]="Selfie with "
PREFIX_ARR[6]="Futuristic style photo of "
RAND_PREFIX=$[ $RANDOM % 7 ]

SUFFIX_ARR[0]="."
SUFFIX_ARR[1]=". Concept art, detail, sharp focus."
SUFFIX_ARR[2]=". Calm, realistic, volumetric Lighting."
SUFFIX_ARR[3]=". Clear definition, unique and one-of-a-kind piece."
SUFFIX_ARR[4]=". Gloomy, dramatic, stunning, dreamy."
SUFFIX_ARR[5]=". Anime, cartoon."
RAND_SUFFIX=$[ $RANDOM % 6 ]

NEG="ugly, tiling, poorly drawn hands, poorly drawn feet, poorly drawn face, out of frame, extra limbs, disfigured, deformed, body out of frame, bad anatomy, watermark, signature, cut off, low contrast, underexposed, overexposed, bad art, beginner, amateur, distorted face, blurry, draft, grainy"

TITLE=${ARR[$RAND_TITLE]}
SUFFIX=${SUFFIX_ARR[$RAND_SUFFIX]}
PREFIX=${PREFIX_ARR[$RAND_PREFIX]}
TIMESTAMP=$(date "+%s")
RESULT_DIR="/app/out/result-${TIMESTAMP}"
mkdir -p $RESULT_DIR
touch "${RESULT_DIR}/info.txt"
PROMPT=${PREFIX}${TITLE}${SUFFIX}

START=$(date +%s);
/app/OnnxStream/src/build/sd --turbo --models-path /app/weights --steps $STEPS --neg-prompt "$NEG" --prompt "$PROMPT" && cp result.png "${RESULT_DIR}/result.png"
END=$(date +%s);
echo $((END-START)) | awk '{print int($1/60)":"int($1%60)}'

echo "title: ${TITLE} \\n promt: ${PROMPT} \\n neg: ${NEG}" >> "${RESULT_DIR}/info.txt"

/app/bin/templater -v TITLE=$TITLE -v IMAGE=$RESULT_DIR/result.png -v TEXT=$PROMPT -t /app/templates/base.html -o $RESULT_DIR/result.html