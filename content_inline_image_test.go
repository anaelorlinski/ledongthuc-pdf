package pdf

import "testing"

// TestContentInlineImage verifies that the text after an inline image is
// returned. The image data between ID and EI is binary; lexed as tokens it
// produced unbalanced strings and dictionaries that swallowed the rest of
// the content stream.
func TestContentInlineImage(t *testing.T) {
	for _, data := range []string{
		"\x00\xff) (",    // string delimiters
		"<<\x00 E I EI2", // dictionary start, E and I not forming the EI operator
	} {
		content := "BT /F1 10 Tf (A) Tj ET BI /W 4 /H 1 /BPC 8 /CS /G ID " + data + " EI BT /F1 10 Tf (B) Tj ET"
		pdf := buildPagePDF("<< /Font << /F1 5 0 R >> >>", content, helveticaFont)
		if got := contentString(firstPageContent(t, pdf)); got != "AB" {
			t.Errorf("image data %q: text = %q, want %q", data, got, "AB")
		}
	}
}
