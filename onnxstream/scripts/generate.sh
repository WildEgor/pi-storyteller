#!/bin/bash

START_TIMESTAMP=$(date "+%s")

RESULT_DIR="/app/out"
NEG="ugly, tiling, poorly drawn hands, poorly drawn feet, poorly drawn face, out of frame, extra limbs, disfigured, deformed, body out of frame, bad anatomy, watermark, signature, cut off, low contrast, underexposed, overexposed, bad art, beginner, amateur, distorted face, blurry, draft, grainy"
STEPS=28
PROMPT="comic photo of (yennefer from vengerberg:1.1) . graphic illustration, comic art, graphic novel art, vibrant, highly detailed"
MODELS="/app/weights"

while [ $# -gt 0 ]; do
  case "$1" in
    -s)
      STEPS=$2
      shift 2
      ;;
    -p)
      PROMPT="$2"
      shift 2
      ;;
    -mp)
      MODELS="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

RPI=""
if [ "$(uname -m)" = "armv7l" ]; then
    RPI="--rpi-lowmem"
fi

if [ "$(uname -m)" = "aarch64" ]; then
    RPI="--rpi"
fi

/app/OnnxStream/src/build/sd --turbo --models-path $MODELS --steps $STEPS $RPI --prompt "$PROMPT" --neg-prompt "$NEG" --output "${RESULT_DIR}/result.png"

END_TIMESTAMP=$(date "+%s")

echo "${START_TIMESTAMP} ${END_TIMESTAMP}"