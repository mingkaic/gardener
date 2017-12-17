package gardener

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
	"testing"
)

//// ====== Structures ======

type mockRC struct {
	*bytes.Buffer
}

//// ====== Tests ======

func TestHTMLValid(t *testing.T) {
	site := GenerateSite(nil)
	hChild := (*site.Children[0]).(HTMLNode)
	htmlTxt := ToHTML(&hChild)

	var rc io.ReadCloser = &mockRC{bytes.NewBufferString(htmlTxt)}
	root, err := html.Parse(rc)
	if err != nil {
		panic(err)
	}
	//root := stew.New(rc)

	// expect root is equivalent to site
	treeCheck(site, root,
		func(msg string, args ...interface{}) {
			t.Errorf(msg, args...)
		})
}

//// ====== Test Utility ======

//// Member of mockRC

// close mock readcloser
func (rc *mockRC) Close() (err error) {
	return
}

//// Check Utility

func treeCheck(expect *HTMLNode, got *html.Node, errCheck func(msg string, args ...interface{})) {
	if expect.Tag != got.Data {
		errCheck("@<%d> expected %s, got %s", expect.Pos, expect.Tag, got.Data)
	}

	expectN := len(expect.Children)
	gotN := 0
	for child := got.FirstChild; child != nil; child = child.NextSibling {
		gotN++
	}
	if expectN != gotN {
		errCheck("@<%d %s> expected %d children, got %d children",
			expect.Pos, expect.Tag, expectN, gotN)
	} else {
		i := 0
		for child := got.FirstChild; child != nil; child = child.NextSibling {
			eChild := (*expect.Children[i]).(HTMLNode)
			treeCheck(&eChild, child, errCheck)
			i++
		}
	}
}
