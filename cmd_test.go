// 実行可能バイナリをビルドして、バイナリのオプションやファイル生成の挙動をテス
// トする。

package main

// func TestANSIColorExt(t *testing.T) {
// 	tests := []struct {
// 		desc   string
// 		inCode int
// 		inFmt  string
// 		outFn  string
// 	}{
// 		{desc: "前景色256色", inCode: 38, inFmt: "\x1b[%d;5;%dm%03d%s", outFn: "anci_color_ext_256_foreground.png"},
// 		{desc: "背景色256色", inCode: 48, inFmt: "\x1b[%d;5;%dm%03d%s", outFn: "anci_color_ext_256_background.png"},
// 		{desc: "前景色RGBグラデーション色", inCode: 38, inFmt: "\x1b[%d;2;%d;0;0m%03d%s", outFn: "anci_color_ext_rgb_foreground.png"},
// 		{desc: "背景色RGBグラデーション色", inCode: 48, inFmt: "\x1b[%d;2;%d;0;0m%03d%s", outFn: "anci_color_ext_rgb_background.png"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.desc, func(t *testing.T) {
// 			var text string
// 			for i := 0; i < 256; i++ {
// 				es := fmt.Sprintf(tt.inFmt, tt.inCode, i, i, ansiColorReset)
// 				text += es
// 				if (i+1)%16 == 0 {
// 					text += "\n"
// 				}
// 			}
// 			cmd := fmt.Sprintf("echo -e '%s' | %s -o %s/%s", text, bin, outDir, tt.outFn)
// 			err := exec.Command("bash", "-c", cmd).Run()
// 			assert.Nil(t, err, tt.desc)
// 		})
// 	}
// }
//
// func TestSlideAnimation(t *testing.T) {
// 	format := "seq 5 | %s -S -l 3 %s -o %s/%s"
//
// 	var cmd string
// 	var err error
// 	cmd = fmt.Sprintf(format, bin, "", outDir, "slide_animation_line_3_width_default.gif")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "行3でスライド幅デフォルトのGIFアニメ")
//
// 	cmd = fmt.Sprintf(format, bin, "-W 2", outDir, "slide_animation_line_3_width_2.gif")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "行3でスライド幅2のGIFアニメ")
// }
//
// func TestSlideAnimationRainbow(t *testing.T) {
// 	format := "echo -e '%s' | %s -S -l 5 %s -o %s/%s"
//
// 	var cmd string
// 	var err error
// 	cmd = fmt.Sprintf(format, fgBgLine, bin, "", outDir, "slide_animation_rainbow.gif")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "スライドアニメ5行")
//
// 	cmd = fmt.Sprintf(format, fgBgLine, bin, "-E", outDir, "slide_animation_rainbow_forever.gif")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "スライドアニメ5行を無限に")
// }
//
// func TestEmojiZ(t *testing.T) {
// 	format := "echo -e '%s' | %s -z -o %s/%s"
//
// 	var cmd string
// 	var err error
// 	cmd = fmt.Sprintf(format, "ab😃cd👍ef👀gh👄\n😃12👍34👀5a👄あ\n😃a👍b👀c👄dabcd", bin, outDir, "emoji_z.png")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "絵文字Z")
// }
//
// func TestError(t *testing.T) {
// 	var cmd string
// 	var err error
//
// 	cmd = fmt.Sprintf(`echo -n "" | %s -o %s/%s`, bin, outDir, "empty.png")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.NotNil(t, err, "空の標準入力はエラーを返す")
//
// 	cmd = fmt.Sprintf(`%s "" -o %s/%s`, bin, outDir, "empty.png")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.NotNil(t, err, "空の引数はエラーを返す")
//
// 	cmd = fmt.Sprintf(`%s "%s" -o %s/%s`, "\n\n\n", bin, outDir, "empty.png")
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.NotNil(t, err, "改行のみの入力はエラーを返す")
// }
//
// func TestNoRedirectNoPipe(t *testing.T) {
// 	var cmd string
// 	var err error
//
// 	const msg = "test"
//
// 	// Note:
// 	//   go testで実行する時は端末ではないためか、絶対にnilになってしまう
// 	//   端末からコマンドを実行する時では動作確認できているためコメントアウト
// 	// err = exec.Command(bin, msg).Run()
// 	// assert.NotNil(t, err, "出力先を未指定の時は異常終了")
//
// 	cmd = fmt.Sprintf(`%s %s > /dev/null`, bin, msg)
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "出力先(リダイレクト)が指定されている時は正常終了")
//
// 	cmd = fmt.Sprintf(`%s %s | cat > /dev/null`, bin, msg)
// 	err = exec.Command("bash", "-c", cmd).Run()
// 	assert.Nil(t, err, "出力先(パイプ)が指定されている時は正常終了")
// }
