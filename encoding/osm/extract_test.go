package osm

import (
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/ctessum/geom"
)

func TestCountTags(t *testing.T) {
	f, err := os.Open("testdata/honolulu_hawaii.osm.pbf")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	tags, err := CountTags(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 19918 {
		t.Errorf("Wrong number of tags %d", len(tags))
	}

	tags2 := tags.Filter(func(t *TagCount) bool {
		return t.Key == "highway" && t.Value == "residential"
	})
	tableWant := [][]string{
		[]string{"Key", "Value", "Total", "Node", "Closed way", "Open way", "Relation"},
		[]string{"highway", "residential", "6839", "0", "55", "6784", "0"}}
	tableHave := tags2.Table()
	if !reflect.DeepEqual(tableWant, tableHave) {
		t.Error("tables don't match")
	}
	if (*tags2)[0].DominantType() != ClosedWay {
		t.Errorf("dominant type should be %d but is %d", ClosedWay, (*tags2)[0].DominantType())
	}
}

func TestExtractTag_Point(t *testing.T) {
	f, err := os.Open("testdata/honolulu_hawaii.osm.pbf")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := ExtractTag(f, "natural", "tree")
	if err != nil {
		t.Fatal(err)
	}
	geomTags, err := data.Geom()
	if err != nil {
		t.Fatal(err)
	}
	if len(geomTags) != 588 {
		t.Errorf("have %d objects, want 588", len(geomTags))
	}
	minx := math.Inf(1)
	miny := math.Inf(1)
	for _, have := range geomTags {
		p := have.Geom.(geom.Point)
		minx = math.Min(minx, p.X)
		miny = math.Min(miny, p.Y)
	}
	const (
		wantx = -158.1244373
		wanty = 21.265047600000003
	)
	if minx != wantx {
		t.Errorf("minimum x value: have %g, want %g", minx, wantx)
	}
	if miny != wanty {
		t.Errorf("minimum y value: have %g, want %g", miny, wanty)
	}
}

func TestExtractTag_Line(t *testing.T) {
	f, err := os.Open("testdata/honolulu_hawaii.osm.pbf")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := ExtractTag(f, "trail_visibility", "bad")
	if err != nil {
		t.Fatal(err)
	}
	geomTags, err := data.Geom()
	if err != nil {
		t.Fatal(err)
	}
	if len(geomTags) != 1 {
		t.Errorf("have %d objects, want 1", len(geomTags))
	}
	want := &GeomTags{
		Geom: geom.LineString{
			geom.Point{X: -157.8260688, Y: 21.404186000000003},
			geom.Point{X: -157.8258194, Y: 21.403686500000003},
		},
		Tags: map[string]string{
			"highway":          "path",
			"surface":          "dirt",
			"trail_visibility": "bad",
			"access":           "private"},
	}
	have := geomTags[0]
	if !reflect.DeepEqual(want, have) {
		t.Errorf("have %#v, want %#v", have, want)
	}
}

func TestExtractTag_Polygon(t *testing.T) {
	f, err := os.Open("testdata/honolulu_hawaii.osm.pbf")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := ExtractTag(f, "name", "Napili Tower")
	if err != nil {
		t.Fatal(err)
	}
	geomTags, err := data.Geom()
	if err != nil {
		t.Fatal(err)
	}
	if len(geomTags) != 1 {
		t.Errorf("have %d objects, want 1", len(geomTags))
	}
	want := &GeomTags{
		Geom: geom.Polygon{
			[]geom.Point{
				geom.Point{X: -157.82454280000002, Y: 21.2800456},
				geom.Point{X: -157.8245124, Y: 21.280018000000002},
				geom.Point{X: -157.8245062, Y: 21.280024},
				geom.Point{X: -157.8244428, Y: 21.279966400000003},
				geom.Point{X: -157.8243979, Y: 21.2800093},
				geom.Point{X: -157.82441010000002, Y: 21.2800204},
				geom.Point{X: -157.8243701, Y: 21.2800587},
				geom.Point{X: -157.82436380000001, Y: 21.2800531},
				geom.Point{X: -157.82432680000002, Y: 21.280088600000003},
				geom.Point{X: -157.8243377, Y: 21.2800985},
				geom.Point{X: -157.8242916, Y: 21.280142700000003},
				geom.Point{X: -157.8242788, Y: 21.280131100000002},
				geom.Point{X: -157.8242406, Y: 21.280167900000002},
				geom.Point{X: -157.82431160000002, Y: 21.280232100000003},
				geom.Point{X: -157.8243443, Y: 21.280200800000003},
				geom.Point{X: -157.8243622, Y: 21.280217},
				geom.Point{X: -157.82442550000002, Y: 21.2801563},
				geom.Point{X: -157.82441540000002, Y: 21.280147200000002},
				geom.Point{X: -157.8244885, Y: 21.2800771},
				geom.Point{X: -157.8244995, Y: 21.280087100000003},
				geom.Point{X: -157.82454280000002, Y: 21.2800456}},
		},
		Tags: map[string]string{
			"addr:city":        "Honolulu",
			"addr:state":       "HI",
			"addr:street":      "Nahua Street",
			"addr:postcode":    "96815",
			"addr:housenumber": "451",
			"name":             "Napili Tower",
			"website":          "http://www.napilitowers.com/",
			"building":         "apartments"},
	}
	have := geomTags[0]
	if !reflect.DeepEqual(want, have) {
		t.Errorf("have %#v, want %#v", have, want)
	}
}

func TestExtractTag_MultiLineString(t *testing.T) {
	f, err := os.Open("testdata/honolulu_hawaii.osm.pbf")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := ExtractTag(f, "wikipedia", "en:Pearl City, Hawaii")
	if err != nil {
		t.Fatal(err)
	}
	geomTags, err := data.Geom()
	if err != nil {
		t.Fatal(err)
	}
	if len(geomTags) != 1 {
		t.Errorf("have %d objects, want 1", len(geomTags))
	}
	switch geomTags[0].Geom.(type) {
	case geom.MultiLineString:
	default:
		t.Errorf("should be a MultiLineString")
	}
}

func TestExtractTag_RelationPolygon(t *testing.T) {
	f, err := os.Open("testdata/honolulu_hawaii.osm.pbf")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := ExtractTag(f, "start_date", "1974")
	if err != nil {
		t.Fatal(err)
	}
	geomTags, err := data.Geom()
	if err != nil {
		t.Fatal(err)
	}
	if len(geomTags) != 1 {
		t.Errorf("have %d objects, want 1", len(geomTags))
	}
	switch typ := geomTags[0].Geom.(type) {
	case geom.Polygon:
	default:
		t.Errorf("should be a Polygon, instead is %#v", typ)
	}
}
