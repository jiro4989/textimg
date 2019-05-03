#!/bin/bash

curl "https://jonasjacek.github.io/colors/" \
  | grep style \
  | sed -r 's@.*<td>([0-9]+)</td>.*<td>(rgb[^<]+)</td>.*@\1:\2@g' \
  | sed -re 's@rgb\(@color.RGBA{@g' -e 's/\)/,255},/g'
