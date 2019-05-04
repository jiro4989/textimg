#!/bin/bash
# vim: tw=0 nowrap:

set -eu

readonly CMD=./bin/textimg
readonly OUTDIR=testdata/out

readonly COLOR_RESET="\x1b[0m"
readonly COLOR_FG_BLACK="\x1b[30m"
readonly COLOR_FG_RED="\X1b[31m"
readonly COLOR_FG_GREEN="\x1b[32m"
readonly COLOR_FG_YELLOW="\x1b[33m"
readonly COLOR_FG_BLUE="\x1b[34m"
readonly COLOR_FG_MAGENTA="\x1b[35m"
readonly COLOR_FG_CYAN="\x1b[36m"
readonly COLOR_FG_WHITE="\x1b[37m"
readonly COLOR_BG_BLACK="\x1b[40m"
readonly COLOR_BG_RED="\X1b[41m"
readonly COLOR_BG_GREEN="\x1b[42m"
readonly COLOR_BG_YELLOW="\x1b[43m"
readonly COLOR_BG_BLUE="\x1b[44m"
readonly COLOR_BG_MAGENTA="\x1b[45m"
readonly COLOR_BG_CYAN="\x1b[46m"
readonly COLOR_BG_WHITE="\x1b[47m"

## 色文字と文字列を合わせて出力し、色付けをリセットする。
##
## @param $1 色
## @param $2 出力する文字列
## @return 色付けされた文字列
echo_color_string() {
  echo -e "$1$2$COLOR_RESET"
}

f_black()   { echo_color_string "$COLOR_FG_BLACK"    "$1"; };
f_red()     { echo_color_string "$COLOR_FG_RED"      "$1"; };
f_green()   { echo_color_string "$COLOR_FG_GREEN"    "$1"; };
f_yellow()  { echo_color_string "$COLOR_FG_YELLOW"   "$1"; };
f_blue()    { echo_color_string "$COLOR_FG_BLUE"     "$1"; };
f_magenta() { echo_color_string "$COLOR_FG_MAGENTA"  "$1"; };
f_cyan()    { echo_color_string "$COLOR_FG_CYAN"     "$1"; };
f_white()   { echo_color_string "$COLOR_FG_WHITE"    "$1"; };

b_black()   { echo_color_string "$COLOR_BG_BLACK"    "$1"; };
b_red()     { echo_color_string "$COLOR_BG_RED"      "$1"; };
b_green()   { echo_color_string "$COLOR_BG_GREEN"    "$1"; };
b_yellow()  { echo_color_string "$COLOR_BG_YELLOW"   "$1"; };
b_blue()    { echo_color_string "$COLOR_BG_BLUE"     "$1"; };
b_magenta() { echo_color_string "$COLOR_BG_MAGENTA"  "$1"; };
b_cyan()    { echo_color_string "$COLOR_BG_CYAN"     "$1"; };
b_white()   { echo_color_string "$COLOR_BG_WHITE"    "$1"; };


## 指定の文字列を指定回数繰り返した文字列を1行出力する。
##
## @param $1 繰り返し回数
## @param $2 繰り返す文字列
## @return $1回繰り返した$2文字列
repeat() {
  seq $1 | xargs -I@ echo -ne "$2"
}

## 渡した文字列から画像を出力する。
## 生成された画像の正当性は目視確認する(爆)。
##
## @param $1 入力文字列
## @param $2 出力ファイル名。出力先はtestdata/out配下で固定
run_test() {
  echo -e "$1" | $CMD -o "$OUTDIR/$2"
}

# ==============================================================================
#
#     ここからテスト開始
#
# ==============================================================================

make build
mkdir -p testdata/out

# ------------------------------------------------------------------------------
#
#     ANSI color
#
# ------------------------------------------------------------------------------

# ANSI colorのテスト
for color in black red green yellow blue magenta cyan white; do
  run_test "$(f_$color $color)" ansi_f_$color.png
  run_test "$(b_$color $color)" ansi_b_$color.png
done

for t in f b; do
  # 1行で色の入れ替わりを試験
  run_test "$(${t}_red red)$(${t}_green green)$(${t}_blue blue)" ansi_${t}_rgb.png

  # 日本語のテスト
  run_test "$(${t}_red 赤あか)$(${t}_green 緑みどり)$(${t}_blue 青あお)" ansi_${t}_rgb_ja.png

  # 複数行出力のテスト
  run_test "$(repeat 10 \"$(${t}_red red)$(${t}_green green)$(${t}_blue blue)\n\")" ansi_${t}_multiline.png

  # 全角文字と半角文字の混在
  run_test "$(${t}_red 赤RR)\n$(${t}_green 緑GG)" ansi_${t}_full_half.png
done
run_test "\x1b[31mRed\x1b[32mGreen\x1b[34mBlue\x1b[0m" ansi_f_rgb2.png

# grepのテスト
echo TestAbcTest | grep --color=always Abc | $CMD -o $OUTDIR/grep.png

# 色の装飾なしのテスト
echo no_color | $CMD -o $OUTDIR/no_color.png

# 背景色指定有りのテスト
echo -e "あいうえおかきくけこ" | sed -r 's/[^　]/\x1b[31m&\x1b[0m/g' | $CMD -b white -o $OUTDIR/ansi_f_set_bg.png

# 背景色をRGBA指定するテスト
colors=(30 31 32 33 34 35 36 37)
i=0
while read -r line; do
  echo -e "$line" | sed -r 's/.*/\x1b['"${colors[$((i%8))]}"'m&\x1b[m/g'
  i=$((i+1))
done <<< "$(seq 8 | xargs -I@ echo TEST)" | $CMD -b 50,100,12,255 -o $OUTDIR/ansi_f_bgopt_rgba.png

# ------------------------------------------------------------------------------
#
#     Extension 256 color
#
# ------------------------------------------------------------------------------

# 拡張256色のテスト (foreground)
seq 0 255 | while read -r i; do
  echo -ne "\x1b[38;5;${i}m$(printf %03d $i)"
  if [ $(((i+1) % 50)) -eq 0 ]; then
    echo
  fi
done | $CMD -o $OUTDIR/ext256_f_rainbow.png

# 拡張256色のテスト (background)
seq 0 255 | while read -r i; do
  echo -ne "\x1b[48;5;${i}m$(printf %03d $i)"
  if [ $(((i+1) % 50)) -eq 0 ]; then
    echo
  fi
done | $CMD -o $OUTDIR/ext256_b_rainbow.png

# ------------------------------------------------------------------------------
#
#     Extension RGBA color
#
# ------------------------------------------------------------------------------

# 拡張RGB指定のテスト (foreground)
seq 0 255 | while read i; do
  echo -ne "\x1b[38;2;${i};0;0m$(printf %03d $i)"
  if [ $(((i+1) % 50)) -eq 0 ]; then
    echo
  fi
done | $CMD -o $OUTDIR/extrgb_f_gradation_red.png

# 拡張RGB指定のテスト (background)
seq 0 255 | while read i; do
  echo -ne "\x1b[48;2;0;$i;0m$(printf %03d $i)"
  if [ $(((i+1) % 50)) -eq 0 ]; then
    echo
  fi
done | $CMD -o $OUTDIR/extrgb_b_gradation_green.png
