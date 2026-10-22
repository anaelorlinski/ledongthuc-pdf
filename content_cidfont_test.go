package pdf

import (
	"strconv"
	"testing"
)

// identityFont returns a Type0 font with Identity-H encoding whose CIDFont
// has the given W array and default width, and whose ToUnicode CMap maps the
// CIDs 1, 2 and 3 to A, B and C.
func identityFont(w string, dw int) []string {
	cmap := "/CIDInit /ProcSet findresource begin\n12 dict begin\nbegincmap\n" +
		"1 begincodespacerange\n<0000> <FFFF>\nendcodespacerange\n" +
		"3 beginbfchar\n<0001> <0041>\n<0002> <0042>\n<0003> <0043>\nendbfchar\n" +
		"endcmap\nCMapName currentdict /CMap defineresource pop\nend\nend"
	return []string{
		"<< /Type /Font /Subtype /Type0 /BaseFont /Test /Encoding /Identity-H" +
			" /DescendantFonts [6 0 R] /ToUnicode 7 0 R >>",
		"<< /Type /Font /Subtype /CIDFontType2 /BaseFont /Test" +
			" /CIDSystemInfo << /Registry (Adobe) /Ordering (Identity) /Supplement 0 >>" +
			" /W " + w + " /DW " + strconv.Itoa(dw) + " >>",
		pdfStream("", cmap),
	}
}

// TestContentCIDFontWidths verifies that the glyphs of a two-byte Type0 font
// advance by the widths of its CIDFont. Content used to look up the widths
// of single bytes in the Type0 font, which has none, so every glyph of a
// string was placed at the same position.
func TestContentCIDFontWidths(t *testing.T) {
	for _, tc := range []struct {
		w    string
		dw   int
		want []float64 // X of A, B and C at size 10
	}{
		{"[1 [600 700 800]]", 1000, []float64{0, 6, 13}}, // consecutive widths
		{"[1 3 500]", 1000, []float64{0, 5, 10}},         // one width for a range
		{"[2 [700]]", 400, []float64{0, 4, 11}},          // DW for CIDs not in W
		{"[]", 0, []float64{0, 0, 0}},                    // explicit DW 0
	} {
		data := buildPagePDF("<< /Font << /F1 5 0 R >> >>",
			"BT /F1 10 Tf <000100020003> Tj ET",
			identityFont(tc.w, tc.dw)...)
		c := firstPageContent(t, data)
		if got := contentString(c); got != "ABC" {
			t.Fatalf("W %s: text = %q, want %q", tc.w, got, "ABC")
		}
		for i, x := range tc.want {
			if c.Text[i].X != x {
				t.Errorf("W %s: glyph %q at x = %v, want %v", tc.w, c.Text[i].S, c.Text[i].X, x)
			}
		}
	}
}
