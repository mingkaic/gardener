package gardener

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
	"gopkg.in/fatih/set.v0"
)

// =============================================
//                    Declarations
// =============================================

// PageNode ...
// Is the pointer content of SiteNode
// to persist values after frequent value conversion
type PageNode struct {
	Depth                        uint
	Hostname, LinkPath, FullLink string
	Page                         *HTMLNode
	Refs                         []*TreeNode
	BackLinks                    []*PageNode
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
	*PageNode
	Info *SiteInfo
}

// PageMap ...
// Associates links to a single HTML page
type PageMap map[string]*PageNode

// =============================================
//                    Globals
// =============================================

const protocol = "http"

// =============================================
//                    Public
// =============================================

//// Gardener Extension

// GeneratePage ...
// Randomly generates a Website graph
func (this Gardener) GenerateSite(nSites uint) *SiteNode {
	hostname := protocol + "://" + uuid.New().String() + ".com"
	linkPath := uuid.New().String()
	fullLink := hostname + "/" + linkPath
	content := &PageNode{Hostname: hostname, LinkPath: linkPath, FullLink: fullLink}
	site := SiteNode{
		content,
		&SiteInfo{Pages: PageMap{fullLink: content}, HostList: []string{hostname}},
	}
	var tOrigin TreeNode = site
	this.RandGraph(&tOrigin, uint(nSites-1))

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
			nElems := nSites*uint(2) + uint(this.Intn(91))
			page.Page = this.GeneratePage(nElems, links)
			if page.Page == nil {
				panic("Generated nil page")
			}
		}
		return page.FullLink
	}
	siteTraversal(&tOrigin, 0)

	return &site
}

//// Members for SiteNode

// NewChild ...
// Make a new Site and add as reference
func (this SiteNode) NewChild(gen *rand.Rand) *TreeNode {
	var hostname string
	nHosts := uint(len(this.Info.HostList))
	if nHosts == 0 {
		hostname = protocol + "://" + uuid.New().String() + ".com"
		this.Info.HostList = []string{hostname}
	} else {
		idx := uint(math.Abs(gen.NormFloat64() * float64(nHosts) / 2))
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
		&PageNode{
			Hostname:  hostname,
			LinkPath:  linkPath,
			FullLink:  hostname + "/" + linkPath,
			BackLinks: []*PageNode{this.PageNode},
		},
		this.Info,
	}
	this.Info.Pages[s.FullLink] = s.PageNode

	var out TreeNode = s
	this.AddChild(&out)
	return &out
}

// AddChild ...
// Add an existing SiteNode as reference
func (this SiteNode) AddChild(child *TreeNode) {
	sChild := (*child).(SiteNode)
	sChild.BackLinks = append(sChild.BackLinks, this.PageNode)
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
		has = has || sThis.PageNode == sChild.PageNode
	}
	return has
}
