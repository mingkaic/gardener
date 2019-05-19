package gardener

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/net/html"
	"gopkg.in/eapache/queue.v1"
)

// =============================================
//                    Globals
// =============================================

const N_TESTS = 100

var gard *Gardener

// =============================================
//                    Tests
// =============================================

func TestMain(m *testing.M) {
	gard = New()
	retCode := 0
	for i := 0; i < N_TESTS && retCode == 0; i++ { // repeat all tests because of randomness
		retCode = m.Run()
	}
	os.Exit(retCode)
}

// TestHTMLValid ...
// Ensures generated page is equal to tree parsed by html package
func TestHTMLValid(t *testing.T) {
	page := gard.GeneratePage(100, nil)
	htmlTxt := ToHTML(page)

	var rc io.ReadCloser = &MockRC{bytes.NewBufferString(htmlTxt)}
	root, err := html.Parse(rc)
	panicCheck(err)

	// expect root is equivalent to page
	treeCheck(page, root,
		func(msg string, args ...interface{}) {
			t.Errorf(msg, args...)
		})
}

// TestPageLinks ...
// Ensures input links are found in the generated page
func TestPageLinks(t *testing.T) {
	var links = make(map[string]struct{})
	nLinks := 35 + gard.Intn(25)
	for i := 0; i < nLinks; i++ {
		links[uuid.New().String()] = struct{}{}
	}
	site := gard.GeneratePage(150, links)
	htmlTxt := ToHTML(site)
	if site.Info.nRemaining != 0 {
		t.Errorf("%d remaining links", site.Info.nRemaining)
	}
	if len(links) > 0 {
		t.Errorf("%d links remaining (not used)", len(links))
	}

	for link := range links {
		lookup := fmt.Sprintf("href=\"%s\"", link)
		if !strings.Contains(htmlTxt, lookup) {
			t.Errorf("missing link: %s", link)
		}
	}
}

// TestSiteValid ...
// Ensures generated site makes sense
func TestSiteValid(t *testing.T) {
	site := gard.GenerateSite(20)
	// perform a breadth first traversal on site
	q := queue.New()
	visited := make(map[*PageNode]struct{})
	q.Add(site)
	for q.Length() > 0 {
		curr := q.Remove().(*SiteNode)
		visited[curr.PageNode] = struct{}{}
		if page, ok := site.Info.Pages[curr.FullLink]; ok {
			if page != curr.PageNode {
				t.Errorf("page with %s link is not current page", curr.FullLink)
			}
		} else {
			t.Errorf("link %s not found", curr.FullLink)
		}

		for _, ref := range curr.Refs {
			sRef := (*ref).(SiteNode)
			if _, ok := visited[sRef.PageNode]; !ok {
				q.Add(&sRef)
			}
		}
	}
}

func TestSitePageLinked(t *testing.T) {
	site := gard.GenerateSite(20)
	for _, page := range site.Info.Pages {
		if page.Page == nil {
			t.Errorf("cannot find page at link %s", page.FullLink)
		} else {
			htmlTxt := ToHTML(page.Page)
			for _, ref := range page.Refs {
				rNode := (*ref).(SiteNode)
				lookup := fmt.Sprintf("href=\"%s\"", rNode.FullLink)
				if !strings.Contains(htmlTxt, lookup) {
					t.Errorf("site %s missing link: %s", page.FullLink, rNode.FullLink)
				}
			}
		}
	}
}

// =============================================
//                    Private
// =============================================

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
