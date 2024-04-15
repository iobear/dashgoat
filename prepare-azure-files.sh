#!/usr/bin/env bash

FILES_FOR_DESTINATION="host.json local.settings.json .funcignore"

if [ -z $1 ]; then
  echo "Creating files"
  for AZFILE in $FILES_FOR_DESTINATION
  do
    cp deploy/azure-functions/$AZFILE .
  done

  mkdir dashgoat
  cp deploy/azure-functions/function.json dashgoat/
  exit 0
fi

if [ $1 = "clean" ]; then
  echo "deleting files"
  for AZFILE in $FILES_FOR_DESTINATION
  do
    rm "$AZFILE"
  done

  rm -rf dashgoat
  exit 0
fi

echo "No match on input"