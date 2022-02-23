package resize

// ComplementWidthHeight は width, height の片方が 0 の時、サイズを調整する。
func ComplementWidthHeight(x, y, w, h int) (int, int) {
	if w == 0 {
		hh := y
		d := float64(h) / float64(hh)
		w = int(float64(x) * d)
		return w, h
	}
	if h == 0 {
		ww := x
		d := float64(w) / float64(ww)
		h = int(float64(y) * d)
		return w, h
	}
	return w, h
}
