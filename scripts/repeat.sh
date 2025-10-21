#!/bin/bash
# repeat.sh - Repeat a character N times
# Usage: repeat.sh <character> <count>

if [ $# -ne 2 ]; then
    echo "Usage: $0 <character> <count>"
    exit 1
fi

char="$1"
count="$2"

# Repeat the character
for ((i=0; i<count; i++)); do
    printf "%s" "$char"
done
printf "\n"
