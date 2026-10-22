package pdf

import (
	"math"
	"testing"
)

// TestContentTextAngle verifies the size, width and direction of rotated
// text. FontSize used to be the first element of the text rendering matrix,
// which is 0 for text rotated by 90 degrees.
func TestContentTextAngle(t *testing.T) {
	for _, tc := range []struct {
		content string
		size    float64
		angle   float64
	}{
		{"BT /F1 12 Tf 1 0 0 1 100 100 Tm (a) Tj ET", 12, 0},
		{"BT /F1 12 Tf 0 1 -1 0 100 100 Tm (a) Tj ET", 12, math.Pi / 2},    // bottom to top
		{"BT /F1 12 Tf -1 0 0 -1 100 100 Tm (a) Tj ET", 12, math.Pi},       // upside down
		{"BT /F1 12 Tf 2 0 0 2 100 100 Tm (a) Tj ET", 24, 0},               // scaled by Tm
		{"q 0 -1 1 0 0 792 cm BT /F1 10 Tf (a) Tj ET Q", 10, -math.Pi / 2}, // rotated by the CTM
	} {
		data := buildPagePDF("<< /Font << /F1 5 0 R >> >>", tc.content, helveticaFont)
		c := firstPageContent(t, data)
		if len(c.Text) != 1 {
			t.Fatalf("%q: got %d glyphs, want 1", tc.content, len(c.Text))
		}
		got := c.Text[0]
		if math.Abs(got.FontSize-tc.size) > 1e-9 {
			t.Errorf("%q: FontSize = %v, want %v", tc.content, got.FontSize, tc.size)
		}
		if want := 0.5 * tc.size; math.Abs(got.W-want) > 1e-9 {
			t.Errorf("%q: W = %v, want %v", tc.content, got.W, want)
		}
		if math.Abs(got.Angle-tc.angle) > 1e-9 {
			t.Errorf("%q: Angle = %v, want %v", tc.content, got.Angle, tc.angle)
		}
	}
}
