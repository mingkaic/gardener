package gardener

import (
	"github.com/google/uuid"
	"gopkg.in/fatih/set.v0"
	"math/rand"
	"math"
)

//// ====== Structures ======

// SiteContent ...
// Is the pointer content of SiteNode
// to persist values after frequent value conversion
type SiteContent struct {
	Depth uint
	Hostname, LinkPath, FullLink string
	Page  *HTMLNode
	Refs  []*TreeNode
}

// SiteInfo ...
// Aggregates Site data for generated site
// Serves as the Ad Hoc information (expected output) when testing web crawlers
type SiteInfo struct {
	Pages    PageMap
	MaxDepth uint
	HostList []string
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

//// ====== Global ======

const protocol = "http"

//// ====== Public ======

//// Members for SiteNode

// NewChild ...
// Make a new Site and add as reference
func (this SiteNode) NewChild() *TreeNode {
	var hostname string
	nHosts := uint(len(this.Info.HostList))
	if nHosts == 0 {
		hostname = protocol + "://" + uuid.New().String() + ".com"
		this.Info.HostList = []string{hostname}
	} else {
		idx := uint(math.Abs(rand.NormFloat64() * float64(nHosts) / 2))
		if idx >= nHosts {
			// add new hostname
			hostname = protocol + "://" + uuid.New().String() + ".com"
			this.Info.HostList = append(this.Info.HostList, hostname)
		} else {
			// use existing hostname
			hostname = this.Info.HostList[idx]
		}
	}

	linkPath := uuid.New().String()
	s := SiteNode{
		&SiteContent{Hostname: hostname, LinkPath: linkPath, FullLink: hostname + "/" + linkPath},
		this.Info,
	}
	this.Info.Pages[s.FullLink] = s.SiteContent

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
	hostname := protocol + "://" + uuid.New().String() + ".com"
	linkPath := uuid.New().String()
	fullLink := hostname + "/" + linkPath
	content := &SiteContent{Hostname: hostname, LinkPath: linkPath, FullLink: fullLink}
	site := SiteNode{
		content,
		&SiteInfo{Pages: PageMap{fullLink: content}, HostList: []string{hostname}},
	}
	var tOrigin TreeNode = site
	RandGraph(&tOrigin, uint(nSites-1))

	// build up Depth, and Page
	visited := set.NewNonTS()
	var siteTraversal func(*TreeNode, uint) string
	siteTraversal = func(node *TreeNode, depth uint) string {
		page := (*node).(SiteNode)
		if site.Info.MaxDepth < depth {
			site.Info.MaxDepth = depth
		}
		if !visited.Has(page.FullLink) {
			visited.Add(page.FullLink)
			page.Depth = depth

			var links set.Interface = set.NewNonTS()
			for _, ref := range page.Refs {
				links.Add(siteTraversal(ref, depth+1))
			}
			// number of elements on a page should be significantly higher than number of potential links
			nElems := nSites*uint(2) + uint(gen.Intn(91))
			page.Page = GeneratePage(nElems, links)
			if page.Page == nil {
				panic("Generated nil page")
			}
		}
		return page.FullLink
	}
	siteTraversal(&tOrigin, 0)

	return &site
}
