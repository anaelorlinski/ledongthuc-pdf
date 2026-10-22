# ledongthuc/pdf — patch stack `patchstack-1`

**This branch is a patch stack, not a development fork.** It is upstream
[ledongthuc/pdf](https://github.com/ledongthuc/pdf) `master` at `5e46805`
(2026-10-20) with the patches below applied on top, one commit per patch, in
this order. A patch stack branch is never rewritten: when upstream moves or
patches get merged, a new one is built (`patchstack-2`, `patchstack-3`, ...).

| # | Patch | Branch (from upstream `master`) | Upstream status |
|---|---|---|---|
| 1 | fix: skip inline image data in content streams | `fix-content-inline-image` | overlaps open PR #89 |
| 2 | fix: Content returns the text of form XObjects | `fix-content-form-xobjects` | PR #97 |
| 3 | fix: Q restores the text encoding of the restored font | `fix-content-font-state` | not proposed yet |
| 4 | fix: apply word spacing in Content | `fix-content-word-spacing` | not proposed yet |
| 5 | fix: use CIDFont widths for Type0 fonts in Content | `fix-content-cid-widths` | not proposed yet |
| 6 | feat: report the size and direction of rotated text | `feat-content-text-angle` | not proposed yet |
| 7 | Fork only: text extraction heuristics in Content | — | not for upstream |
| 8 | This README section | — | not for upstream |

The module path is unchanged (`github.com/ledongthuc/pdf`); use it with this line in `go.mod`, then `go mod tidy` pins it to a commit:

```
replace github.com/ledongthuc/pdf => github.com/anaelorlinski/ledongthuc-pdf patchstack-1
```

Building the next stack: branch from the new upstream `master`, cherry-pick the
patches that are still needed (drop the merged ones), then compare with
`git range-diff master...patchstack-1 master...patchstack-2`.

---

# PDF Reader

[![Built with WeBuild](https://raw.githubusercontent.com/webuild-community/badge/master/svg/WeBuild.svg)](https://webuild.community)

A simple Go library which enables reading PDF files. Forked from https://github.com/rsc/pdf

Features
  - Get plain text content (without format)
  - Get Content (including all font and formatting information)

## Install:

`go get -u github.com/ledongthuc/pdf`

## Examples:

 - Check in examples/ folder


## Read plain text

```golang
package main

import (
	"bytes"
	"fmt"

	"github.com/ledongthuc/pdf"
)

func main() {
	pdf.DebugOn = true

	f, r, err := pdf.Open("./pdf_test.pdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		panic(err)
	}
	buf.ReadFrom(b)
	content := buf.String()
	fmt.Println(content)
}
```

## Read all text with styles from PDF

```golang
package main

import (
	"fmt"

	"github.com/ledongthuc/pdf"
)

func main() {
	f, r, err := pdf.Open("./pdf_test.pdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sentences, err := r.GetStyledTexts()
	if err != nil {
		panic(err)
	}

	// Print all sentences
	for _, sentence := range sentences {
		fmt.Printf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n",
			sentence.Font,
			sentence.FontSize,
			sentence.X,
			sentence.Y,
			sentence.S)
	}
}
```


## Read text grouped by rows

```golang
package main

import (
	"fmt"
	"os"

	"github.com/ledongthuc/pdf"
)

func main() {
	content, err := readPdf(os.Args[1]) // Read local pdf file
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
	return
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() || p.V.Key("Contents").Kind() == pdf.Null {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
		    println(">>>> row: ", row.Position)
		    for _, word := range row.Content {
		        fmt.Println(word.S)
		    }
		}
	}
	return "", nil
}
```

## Demo
![Run example](https://i.gyazo.com/01fbc539e9872593e0ff6bac7e954e6d.gif)
