#!/bin/bash

## ファイル名に含まれる文字列から圧縮形式を変更する。

set -eu

readonly DIR=$1

if [[ "$DIR" =~ .*windows.* ]]; then
  zip -r $DIR.zip $DIR
else
  tar czf $DIR.tar.gz $DIR
fi
