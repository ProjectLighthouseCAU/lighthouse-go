package colors

func HslToRgb(h, s, l float64) (r, g, b float64) {

	if s == 0 {
		return l, l, l
	}

	var v1, v2 float64
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - (s * l)
	}

	v1 = 2*l - v2

	_r := hueToRGB(v1, v2, h+(1.0/3.0))
	_g := hueToRGB(v1, v2, h)
	_b := hueToRGB(v1, v2, h-(1.0/3.0))

	return _r, _g, _b
}

func hueToRGB(v1, v2, h float64) float64 {
	if h < 0 {
		h += 1
	}
	if h > 1 {
		h -= 1
	}
	switch {
	case 6*h < 1:
		return (v1 + (v2-v1)*6*h)
	case 2*h < 1:
		return v2
	case 3*h < 2:
		return v1 + (v2-v1)*((2.0/3.0)-h)*6
	}
	return v1
}
