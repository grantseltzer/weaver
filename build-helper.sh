#!/bin/bash

listOfSOs=`ldd ./bin/weaver | awk 'NF == 4 { print $3 }; NF == 2 { print $1 }'`

tempBuild=$1

mkdir -p $tempBuild

for soPath in $listOfSOs
do
  dir=`dirname $soPath`
  file=`basename $soPath`
  mkdir -p "$tempBuild/$dir"
  cp $soPath "$tempBuild/$dir"
done
