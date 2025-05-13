package grid

type Bounds struct {
	Top    uint8
	Right  uint8
	Bottom uint8
	Left   uint8
}

func (b Bounds) Area() int {
	return int(absdiff(b.Top, b.Bottom) * absdiff(b.Left, b.Right))
}

func (bounds Bounds) GridKeys() []GridKey {
	keys := make([]GridKey, 0, bounds.Area())
	for i := bounds.Top; i <= bounds.Bottom; i++ {
		for j := bounds.Left; j <= bounds.Right; j++ {
			keys = append(keys, GridKey{i, j})
		}
	}
	return keys
}

func (b Bounds) InBounds(key GridKey) bool {
	return key.Line >= b.Top &&
		key.Line <= b.Bottom &&
		key.Beat >= b.Left &&
		key.Beat <= b.Right
}

func (b Bounds) Normalized() Bounds {
	return Bounds{Top: 0, Right: b.Right - b.Left, Bottom: b.Bottom - b.Top, Left: 0}
}

func (b Bounds) BottomRightFrom(key GridKey) GridKey {
	return GridKey{key.Line + b.Bottom, key.Beat + b.Right}
}

func (b Bounds) TopLeft() GridKey {
	return GridKey{b.Top, b.Left}
}

func absdiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	} else {
		return b - a
	}
}
