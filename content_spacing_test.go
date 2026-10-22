package pdf

import "testing"

// TestContentWordSpacing verifies that the word spacing set by Tw and by the
// " operator moves the text after a space (single-byte code 32), and only
// after a space.
func TestContentWordSpacing(t *testing.T) {
	for _, tc := range []struct {
		content string
		want    []float64 // X of each glyph
	}{
		// Width 500 at size 10 advances 5 per glyph; the space adds Tw = 4.
		{"BT /F1 10 Tf 4 Tw (a b) Tj ET", []float64{0, 5, 14}},
		{"BT /F1 10 Tf 4 0 (a b) \" ET", []float64{0, 5, 14}},
		{"BT /F1 10 Tf 4 Tw (ab) Tj ET", []float64{0, 5}},
		// Horizontal scaling applies to the word spacing too.
		{"BT /F1 10 Tf 4 Tw 50 Tz (a b) Tj ET", []float64{0, 2.5, 7}},
	} {
		data := buildPagePDF("<< /Font << /F1 5 0 R >> >>", tc.content, helveticaFont)
		c := firstPageContent(t, data)
		if len(c.Text) != len(tc.want) {
			t.Errorf("%q: got %d glyphs, want %d", tc.content, len(c.Text), len(tc.want))
			continue
		}
		for i, x := range tc.want {
			if c.Text[i].X != x {
				t.Errorf("%q: glyph %d (%q) at x = %v, want %v", tc.content, i, c.Text[i].S, c.Text[i].X, x)
			}
		}
	}
}
