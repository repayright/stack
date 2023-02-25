#!/usr/bin/env bash

scriptFile=$(realpath $0)
scriptDir=$(dirname $scriptFile)
randomFile=/tmp/$RANDOM
components=$(ls $scriptDir/../components)
for mod in $components; do
  echo "cd components/${mod} && go mod tidy" >> $randomFile
done

parallel < $randomFile
