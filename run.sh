#!/bin/bash

repeat() {
  seq $1 | xargs -I@ echo -ne "$2"
}

make build
echo -e "\x1b[30mUNKO\x1b[0m" | ./bin/txtimg out/black.png
echo -e "\x1b[31mUNKO\x1b[0m" | ./bin/txtimg out/red.png
echo -e "\x1b[32mUNKO\x1b[0m" | ./bin/txtimg out/green.png
echo -e "\x1b[33mUNKO\x1b[0m" | ./bin/txtimg out/yellow.png
echo -e "\x1b[34mUNKO\x1b[0m" | ./bin/txtimg out/blue.png
echo -e "\x1b[35mUNKO\x1b[0m" | ./bin/txtimg out/magenta.png
echo -e "\x1b[36mUNKO\x1b[0m" | ./bin/txtimg out/syan.png
echo -e "\x1b[37mUNKO\x1b[0m" | ./bin/txtimg out/white.png
echo -e "\x1b[31mRed\x1b[32mGreen\x1b[34mBlue\x1b[0m" | ./bin/txtimg out/rgb.png
echo -e "$(repeat 10 "\x1b[31mR\x1b[32mG\x1b[34mB")\x1b[0m" | ./bin/txtimg out/rgb2.png
echo -e "$(repeat 10 "\x1b[31m赤\x1b[32m緑\x1b[34m青")\x1b[0m" | ./bin/txtimg out/rgb_ja.png
echo -e "$(repeat 10 "\x1b[31mあか\x1b[32mみどり\x1b[34mあお")\x1b[0m" | ./bin/txtimg out/rgb2_ja.png
echo ---------
echo -e "$(repeat 5 "\x1b[31m赤\x1b[32m緑\x1b[34m青\n")" | ./bin/txtimg out/rgblf_ja.png
echo -e "$(repeat 3 "\x1b[31m赤\x1b[32m緑\x1b[34m青\n\x1b[31mRR\x1b[32mGG\x1b[34mBB\n")" | ./bin/txtimg out/rgb2lf_ja.png
echo Test | grep --color=always Te | ./bin/txtimg out/grep.png
