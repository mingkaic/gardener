package gardener

import (
	"github.com/google/uuid"
	"gopkg.in/fatih/set.v0"
)

//// ====== Structures ======

// SiteContent ...
// Is the pointer content of SiteNode
// to persist values after frequent value conversion
type SiteContent struct {
	Depth      uint
	Link, Page string
	Refs       []*TreeNode
}

// SiteInfo ...
// Aggregates Site data for generated site
// Serves as the Ad Hoc information (expected output) when testing web crawlers
type SiteInfo struct {
	Pages PageMap
}

// SiteNode ...
// Is a TreeNode implementation and
// an atomic element of a site
type SiteNode struct {
	*SiteContent
	Info *SiteInfo
}

// PageMap ...
// Associates links to a single HTML page
type PageMap map[string]*SiteContent

//// ====== Public ======

//// Members for SiteNode

// NewChild ...
// Make a new Site and add as reference
func (this SiteNode) NewChild() *TreeNode {
	link := uuid.New().String()
	s := SiteNode{
		&SiteContent{Link: link},
		this.Info,
	}
	this.Info.Pages[link] = s.SiteContent

	var out TreeNode = s
	this.AddChild(&out)
	return &out
}

// AddChild ...
// Add an existing SiteNode as reference
func (this SiteNode) AddChild(child *TreeNode) {
	this.Refs = append(this.Refs, child)
}

// HasChild ...
// Check if this already have SiteNode as a reference
func (this SiteNode) HasChild(child *TreeNode) bool {
	has := false
	n := len(this.Refs)
	sChild := (*child).(SiteNode)
	for i := 0; i < n && !has; i++ {
		sThis := (*this.Refs[i]).(SiteNode)
		has = has || sThis.SiteContent == sChild.SiteContent
	}
	return has
}

//// Core Functions

// GeneratePage ...
// Randomly generates a Website graph
func GenerateSite(nSites uint) *SiteNode {
	link := uuid.New().String()
	site := SiteNode{
		&SiteContent{Link: link},
		&SiteInfo{Pages: make(PageMap)},
	}
	site.Info.Pages[link] = site.SiteContent
	var tOrigin TreeNode = site
	RandGraph(&tOrigin, uint(nSites - 1))

	// build up Depth, and Page
	visited := set.NewNonTS()
	var siteTraversal func(*TreeNode, uint) string
	siteTraversal = func(node *TreeNode, depth uint) string {
		page := (*node).(SiteNode)
		if !visited.Has(page.Link) {
			page.Depth = depth

			var links set.Interface = set.NewNonTS()
			for _, ref := range page.Refs {
				links.Add(siteTraversal(ref, depth+1))
			}
			// number of elements on a page should be significantly higher than number of potential links
			nElems := nSites * uint(2) + uint(gen.Intn(91))
			htmlPage := GeneratePage(nElems, links)
			page.Page = ToHTML(htmlPage)
		}
		return page.Link
	}

	return &site
}
