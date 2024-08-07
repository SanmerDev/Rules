#!/bin/bash

for file in $(ls);
do
  if [ -f "$file" ] && [[ "$file" == *.json ]]; then
    echo "Compiling $file"
    sing-box rule-set compile $file
  fi
done
