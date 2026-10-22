package pdf

import "testing"

// TestContentFontRestoredByQ verifies that Q restores the text encoding
// together with the font. The encoding used to be kept outside the graphics
// state, so after "q /F2 Tf ... Q" the text of the restored font was still
// decoded with the encoding of F2.
func TestContentFontRestoredByQ(t *testing.T) {
	greek := "<< /Type /Font /Subtype /Type1 /BaseFont /Symbol /FirstChar 65 /LastChar 65 /Widths [500]" +
		" /Encoding << /Type /Encoding /Differences [65 /Alpha] >> >>"
	data := buildPagePDF("<< /Font << /F1 5 0 R /F2 6 0 R >> >>",
		"BT /F1 10 Tf q /F2 10 Tf (A) Tj Q (A) Tj ET",
		helveticaFont, greek)

	c := firstPageContent(t, data)
	if got, want := contentString(c), "\u0391A"; got != want {
		t.Fatalf("text = %q, want %q", got, want)
	}
	if c.Text[1].Font != "Helvetica" {
		t.Errorf("font after Q = %q, want %q", c.Text[1].Font, "Helvetica")
	}
}

// TestContentTextBeforeTf verifies that text shown before any Tf operator is
// still returned, undecoded.
func TestContentTextBeforeTf(t *testing.T) {
	data := buildPagePDF("<< >>", "BT (ab) Tj ET")
	if got := contentString(firstPageContent(t, data)); got != "ab" {
		t.Errorf("text = %q, want %q", got, "ab")
	}
}
