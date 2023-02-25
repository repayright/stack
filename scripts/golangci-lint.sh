#!/usr/bin/env bash

set -x;
set -e;

scriptFile=$(realpath $0)
scriptDir=$(dirname $scriptFile)
randomFile=/tmp/$RANDOM
components=$(ls $scriptDir/../components)
for mod in $components; do
  echo "cd components/${mod} && golangci-lint --allow-parallel-runners -v run --fix" >> $randomFile
done

parallel < $randomFile
