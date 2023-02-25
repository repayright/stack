#!/usr/bin/env bash

scriptFile=$(realpath $0)
scriptDir=$(dirname $scriptFile)
randomFile=/tmp/$RANDOM
components=$(ls $scriptDir/../components)
for mod in $components; do
  echo "Run tests on ${mod}"
  pushd $scriptDir/../components/${mod} >/dev/null
  [[ -e Taskfile.yml ]] && task tests
  popd >/dev/null
done

echo "Run $cmdLine"
$cmdLine
