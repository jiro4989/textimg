#!/bin/bash
# vim: tw=0 nowrap:

readonly CMD=./bin/textimg
readonly OUTDIR=testdata/out
test_count=0
err_count=0

if [ "$TEXTIMG_EMOJI_DIR" = "" ]; then
  export TEXTIMG_EMOJI_DIR=/usr/local/src/noto-emoji/png/128
fi

if [ "$TEXTIMG_FONT_FILE" = "" ]; then
  export TEXTIMG_FONT_FILE=/usr/share/fonts/TTF/HackGen-Regular.ttf
fi

# Symbolaフォント指定
if [ "$TEXTIMG_EMOJI_FONT_FILE" = "" ]; then
  export TEXTIMG_EMOJI_FONT_FILE=/usr/share/fonts/truetype/symbola/Symbola.ttf
fi

# 色のANSIエスケープシーケンス定数 {{{

readonly COLOR_RESET="\x1b[0m"
readonly COLOR_FG_BLACK="\x1b[30m"
readonly COLOR_FG_RED="\x1b[31m"
readonly COLOR_FG_GREEN="\x1b[32m"
readonly COLOR_FG_YELLOW="\x1b[33m"
readonly COLOR_FG_BLUE="\x1b[34m"
readonly COLOR_FG_MAGENTA="\x1b[35m"
readonly COLOR_FG_CYAN="\x1b[36m"
readonly COLOR_FG_WHITE="\x1b[37m"
readonly COLOR_BG_BLACK="\x1b[40m"
readonly COLOR_BG_RED="\x1b[41m"
readonly COLOR_BG_GREEN="\x1b[42m"
readonly COLOR_BG_YELLOW="\x1b[43m"
readonly COLOR_BG_BLUE="\x1b[44m"
readonly COLOR_BG_MAGENTA="\x1b[45m"
readonly COLOR_BG_CYAN="\x1b[46m"
readonly COLOR_BG_WHITE="\x1b[47m"

#}}}

# ユーティリティ関数 {{{

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

suite() { echo -e "$(f_blue [Suite]) $1"; };
info() { echo -e "  $(f_green "[OK]") $1"          ; };
err()  { echo -e "  $(f_red   "[NG]") $1" ; };

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
  local desc
  local inputstr
  local outfile
  local exitcode

  desc="$1"
  inputstr="$2"
  outfile="$OUTDIR/$3"
  echo -e "$inputstr" | $CMD ${opts[@]} -o "$outfile"

  exitcode=$?
  test_count=$((test_count + 1))
  if [ "$exitcode" -eq 0 ]; then
    info "$desc"
  else
    err "$desc"
    err_count=$((err_count + 1))
  fi
}

#}}}

# ==============================================================================
#
#     ここからテスト開始
#
# ==============================================================================

make build || { echo "$(f_red Failed to build application)"; exit 1; };
mkdir -p testdata/out

# Test: ANSIカラー{{{

suite "ANSI color tests"

for color in black red green yellow blue magenta cyan white; do
  run_test "Foreground ANSI color ($color)" "$(f_$color $color)" ansi_f_$color.png
  run_test "Background ANSI color ($color)" "$(b_$color $color)" ansi_b_$color.png
done

for t in f b; do
  run_test "Switch color on 1 line ($t)" "$(${t}_red red)$(${t}_green green)$(${t}_blue blue)" ansi_${t}_rgb.png
  run_test "Japanese text ($t)" "$(${t}_red 赤あか)$(${t}_green 緑みどり)$(${t}_blue 青あお)" ansi_${t}_rgb_ja.png
  run_test "Multiline text ($t)" "$(repeat 10 $(${t}_red red)$(${t}_green green)$(${t}_blue blue)\\n)" ansi_${t}_multiline.png
  run_test "Half width and Full width characters ($t)" "$(${t}_red 赤RR)\n$(${t}_green 緑GG)" ansi_${t}_full_half.png
done
run_test "No reset color" "\x1b[31mRed\x1b[32mGreen\x1b[34mBlue\x1b[0m" ansi_f_rgb2.png
run_test "Grep color" "$(echo TestAbcTest | grep --color=always Abc)" grep.png
run_test "No escape sequence" no_color no_color.png
run_test "CLI background option" "$(echo -e "あいうえおかきくけこ" | sed -E 's/[^　]/\x1b[31m&\x1b[0m/g')" ansi_f_set_bg.png

# 背景色をRGBA指定するテスト
colors=(30 31 32 33 34 35 36 37)
i=0
while read -r line; do
  echo -e "$line" | sed -E 's/.*/\x1b['"${colors[$((i%8))]}"'m&\x1b[m/g'
  i=$((i+1))
done <<< "$(seq 8 | xargs -I@ echo TEST)" | $CMD -b 50,100,12,255 -o $OUTDIR/ansi_f_bgopt_rgba.png

run_test "Output JPG" "$(f_red RedJPG)" ansi_f_red.jpg
run_test "Output GIF" "$(f_red RedGIF)" ansi_f_red.gif

# 引数から指定
$CMD "$(f_red RedArgs)" -o $OUTDIR/ansi_f_red_args.png

# 全体の文字色を変更
$CMD "Normal$(f_red Red)Normal" --foreground green -o $OUTDIR/ansi_f_changefg.png
$CMD "Normal$(f_red Red)Normal" --foreground 255,255,0,255 -o $OUTDIR/ansi_f_changefg2.png
$CMD "Normal$(f_red Red)Normal" --foreground 0,0,0,0 -o $OUTDIR/ansi_f_changefg3.png

run_test "Reverse color" "\x1b[31;42mRedGreen\x1b[7mRedGreen" ansi_fb_reverse.png

echo -e '\x1b[31mText\x1b[0m
\x1b[32mText\x1b[0m
\x1b[33mText\x1b[0m
\x1b[34mText\x1b[0m
\x1b[35mText\x1b[0m
\x1b[36mText\x1b[0m
\x1b[37mText\x1b[0m
\x1b[41mText\x1b[0m
\x1b[42mText\x1b[0m
\x1b[43mText\x1b[0m
\x1b[44mText\x1b[0m
\x1b[45mText\x1b[0m
\x1b[46mText\x1b[0m
\x1b[47mText\x1b[0m' | $CMD -a -o $OUTDIR/ansi_fb_anime_1line.gif

echo -e '\x1b[31mText\x1b[0m
\x1b[32mText\x1b[0m
\x1b[41mText\x1b[0m
\x1b[42mText\x1b[0m' | $CMD -a -d 60 -l 2 -o $OUTDIR/ansi_fb_anime_2line.gif

#}}}

# Test: 拡張256色 {{{

suite "Extension 256 color tests"

echo_256rainbow() {
  seq 0 255 | while read i; do
    echo -ne "\x1b[$1;5;${i}m$(printf %03d $i)"
    if [ $(((i+1) % 16)) -eq 0 ]; then
      echo
    fi
  done
}

run_test "Foreground rainbow" "$(echo_256rainbow 38)" ext256_f_rainbow.png
run_test "Background rainbow" "$(echo_256rainbow 48)" ext256_b_rainbow.png

#}}}

# Test: 拡張256色(RGB) {{{

suite "Extension 256 color (RGB) tests"

echo_rgb_gradation() {
  seq 0 255 | while read i; do
    echo -ne "\x1b[$1;2;${i};0;0m$(printf %03d $i)"
    if [ $(((i+1) % 16)) -eq 0 ]; then
      echo
    fi
  done
}

run_test "Foreground gradation" "$(echo_rgb_gradation 38)" extrgb_f_gradation.png
run_test "Background gradation" "$(echo_rgb_gradation 48)" extrgb_b_gradation.png

#}}}

# Test: スライドアニメーション {{{

suite "Slide animation tests"

seq 5 | $CMD -l 3 -S -o $OUTDIR/slide_3_1.gif
seq 5 | $CMD -l 3 -SW 2 -o $OUTDIR/slide_3_2.gif
echo -e '\x1b[31mText\x1b[0m
\x1b[32mText\x1b[0m
\x1b[33mText\x1b[0m
\x1b[34mText\x1b[0m
\x1b[35mText\x1b[0m
\x1b[36mText\x1b[0m
\x1b[37mText\x1b[0m
\x1b[41mText\x1b[0m
\x1b[42mText\x1b[0m
\x1b[43mText\x1b[0m
\x1b[44mText\x1b[0m
\x1b[45mText\x1b[0m
\x1b[46mText\x1b[0m
\x1b[47mText\x1b[0m' | $CMD -l 5 -S -o $OUTDIR/slide_5_1_rainbow.gif

echo -e '\x1b[31mText\x1b[0m
\x1b[32mText\x1b[0m
\x1b[33mText\x1b[0m
\x1b[34mText\x1b[0m
\x1b[35mText\x1b[0m
\x1b[36mText\x1b[0m
\x1b[37mText\x1b[0m
\x1b[41mText\x1b[0m
\x1b[42mText\x1b[0m
\x1b[43mText\x1b[0m
\x1b[44mText\x1b[0m
\x1b[45mText\x1b[0m
\x1b[46mText\x1b[0m
\x1b[47mText\x1b[0m' | $CMD -l 5 -SE -o $OUTDIR/slide_5_1_rainbow_forever.gif

#}}}

# Test: 絵文字 {{{

suite "Emoji"

run_test "Draw 1 line emoji" "あ😃a👍！👀ん👄" emoji1.png

run_test "Draw 2 line emoji" "あ😃い👍う👀え👄
😃い👍う👀え👄あ" emoji2.png

run_test "Draw 3 line emoji" "ab😃cd👍ef👀gh👄
😃12👍34👀5a👄あ
😃a👍b👀c👄dabcd" emoji3.png

# Symbola fontを使う指定
opts=(-i)

run_test "(Symbola)Draw 1 line emoji " "あ😃a👍！👀ん👄" emoji1_symbola.png

run_test "(Symbola)Draw 2 line emoji" "あ😃い👍う👀え👄
😃い👍う👀え👄あ" emoji2_symbola.png

run_test "(Symbola)Draw 3 line emoji" "ab😃cd👍ef👀gh👄
😃12👍34👀5a👄あ
😃a👍b👀c👄dabcd" emoji3_symbola.png

#}}}

if [ "$err_count" -lt 1 ]; then
  echo -e "$(f_green Success:) [$test_count/$test_count] tests passed"
  exit 0
else
  echo -e "$(f_red Failure:) [$err_count/$test_count] tests don't passed"
  exit 1
fi
