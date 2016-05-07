package geom

// MultiPoint is a holder for multiple related points.
type MultiPoint []Point

// Bounds gives the rectangular extents of the MultiPoint.
func (mp MultiPoint) Bounds() *Bounds {
	b := NewBounds()
	for _, p := range mp {
		b.extendPoint(p)
	}
	return b
}

// Within calculates whether all of the points in mp are within poly.
func (mp MultiPoint) Within(poly Polygonal) bool {
	for _, p := range mp {
		if !pointInPolygonal(p, poly) {
			return false
		}
	}
	return true
}
