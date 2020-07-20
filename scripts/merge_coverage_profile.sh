#!/bin/bash

# https://songmu.jp/riji/entry/2015-01-15-goveralls-multi-package.html

set -e

cleanup() {
  retval=$?
  if [ $tmpprof != "" ] && [ -f $tmpprof ]; then
    rm -f $tmpprof
  fi
  exit $retval
}
trap cleanup INT QUIT TERM EXIT

# メインの処理
prof=${1:-"coverage.txt"}
echo "mode: count" > $prof
pkgroot=$(go list)
for pkg in $(go list ./...); do
  dir=$(echo $pkg | sed "s,$pkgroot,.,")
  tmpprof=$dir/profile.tmp
  go test -covermode=count -coverprofile=$tmpprof $pkg
  if [ -f $tmpprof ]; then
    cat $tmpprof | tail -n +2 >> $prof
    rm $tmpprof
  fi
done
