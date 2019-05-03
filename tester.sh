#!/bin/bash
# vim: tw=0 nowrap:

set -eu

repeat() {
  seq $1 | xargs -I@ echo -ne "$2"
}

make build
mkdir -p testdata/out
echo -e "\x1b[30mTESO\x1b[0m" | ./bin/textimg -o testdata/out/black.png
echo -e "\x1b[31mTESO\x1b[0m" | ./bin/textimg -o testdata/out/red.png
echo -e "\x1b[32mTESO\x1b[0m" | ./bin/textimg -o testdata/out/green.png
echo -e "\x1b[33mTESO\x1b[0m" | ./bin/textimg -o testdata/out/yellow.png
echo -e "\x1b[34mTESO\x1b[0m" | ./bin/textimg -o testdata/out/blue.png
echo -e "\x1b[35mTESO\x1b[0m" | ./bin/textimg -o testdata/out/magenta.png
echo -e "\x1b[36mTESO\x1b[0m" | ./bin/textimg -o testdata/out/syan.png
echo -e "\x1b[37mTESO\x1b[0m" | ./bin/textimg -o testdata/out/white.png
echo -e "\x1b[31mRed\x1b[32mGreen\x1b[34mBlue\x1b[0m" | ./bin/textimg -o testdata/out/rgb.png
echo -e "$(repeat 10 "\x1b[31mR\x1b[32mG\x1b[34mB")\x1b[0m" | ./bin/textimg -o testdata/out/rgb2.png
echo -e "$(repeat 10 "\x1b[31m赤\x1b[32m緑\x1b[34m青")\x1b[0m" | ./bin/textimg -o testdata/out/rgb_kan.png
echo -e "$(repeat 10 "\x1b[31mあか\x1b[32mみどり\x1b[34mあお")\x1b[0m" | ./bin/textimg -o testdata/out/rgb_hira.png
echo -e "$(repeat 5 "\x1b[31m赤\x1b[32m緑\x1b[34m青\n")" | ./bin/textimg -o testdata/out/rgb_kan_multiline.png
echo -e "$(repeat 3 "\x1b[31m赤\x1b[32m緑\x1b[34m青\n\x1b[31mRR\x1b[32mGG\x1b[34mBB\n")" | ./bin/textimg -o testdata/out/rgb_kan_multiline_halfandfull.png
echo TestAbcTest | grep --color=always Te | ./bin/textimg -o testdata/out/grep.png
echo normal | ./bin/textimg -o testdata/out/normal.png
echo -e "\x1b[31mRED\x1b[0mWhite" | ./bin/textimg -o testdata/out/red_white.png
echo -e "\x1b[31mRED\x1b[0mWhite" | ./bin/textimg -b red -o testdata/out/red_bg.png
echo -e "あいうえおかきくけこ" | sed -r 's/[^　]/\x1b[31m&\x1b[0m/g' | ./bin/textimg -b white -o testdata/out/hiragana.png
echo -e "\x1b[41mR\x1b[31mR\x1b[42mG\x1b[32mG\x1b[44mB\x1b[34mB\x1b[0m\nR\x1b[31mRG\x1b[32mGB\x1b[34mB\x1b[0m" | ./bin/textimg -b white -o testdata/out/bg.png

seq 0 255 | while read -r i; do
  echo -ne "\x1b[38;5;${i}m$(printf %03d $i)"
  if [ $(((i+1) % 50)) -eq 0 ]; then
    echo
  fi
done | ./bin/textimg -o testdata/out/256_fg.png

seq 0 255 | while read -r i; do
  echo -ne "\x1b[48;5;${i}m$(printf %03d $i)"
  if [ $(((i+1) % 50)) -eq 0 ]; then
    echo
  fi
done | ./bin/textimg -o testdata/out/256_bg.png

colors=(30 31 32 33 34 35 36 37)
i=0
while read -r line; do
  echo -e "$line" | sed -r 's/.*/\x1b['"${colors[$((i%8))]}"'m&\x1b[m/g'
  i=$((i+1))
done <<< "$(seq 8 | xargs -I@ echo TEST)" | ./bin/textimg -b 50,100,12,255 -o testdata/out/rainbow.png
