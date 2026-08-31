package state

// Every pipeline artifact carries a retrieval frontmatter block: the seven
// fields validate-artifact checks (SKILL.md step 5) and the retrieval
// corpus indexes on (epic 59, write-time enrichment). A typed artifact
// without it would be strictly worse for retrieval than the markdown one it
// replaces, so the state carries what the block needs and the renderer
// emits it.
//
// Three of the seven are derived from fields the state already has —
// feature, bounded context, and touched files — so only the four an agent
// must supply are modelled here.

// Retrieval is the part of the frontmatter an agent provides. Every field
// is optional: a feature that references no ADR should say so with an empty
// list, not by inventing one.
type Retrieval struct {
	DomainTerms []string `json:"domainTerms,omitempty" jsonschema:"description=Canonical terms from DOMAIN_DICTIONARY.md used by this feature"`
	IssueRefs   []string `json:"issueRefs,omitempty" jsonschema:"description=Ticket or issue references, e.g. PROJ-123"`
	LinkedADRs  []string `json:"linkedAdrs,omitempty" jsonschema:"description=Repo-relative paths to referenced ADRs"`
	LinkedKIs   []string `json:"linkedKis,omitempty" jsonschema:"description=Repo-relative paths to referenced Knowledge Items"`
}

// frontmatter is the fully resolved block, derived fields included.
type frontmatter struct {
	feature        string
	boundedContext string
	domainTerms    []string
	filesTouched   []string
	issueRefs      []string
	linkedADRs     []string
	linkedKIs      []string
}

func (a AnalysisState) frontmatter() frontmatter {
	return newFrontmatter(a.Feature, a.BoundedContext.Owning, a.Retrieval, affectedPaths(a.AffectedComponents))
}

func (a ArchitectureState) frontmatter() frontmatter {
	return newFrontmatter(a.Feature, a.BoundedContext.Owning, a.Retrieval, placementPackages(a.ComponentPlacement))
}

func newFrontmatter(feature, boundedContext string, retrieval Retrieval, filesTouched []string) frontmatter {
	return frontmatter{
		feature: feature, boundedContext: boundedContext,
		domainTerms: retrieval.DomainTerms, filesTouched: filesTouched,
		issueRefs: retrieval.IssueRefs, linkedADRs: retrieval.LinkedADRs, linkedKIs: retrieval.LinkedKIs,
	}
}

// affectedPaths derives files_touched from the components the analyst named.
func affectedPaths(components []AffectedComponent) []string {
	paths := make([]string, 0, len(components))
	for _, component := range components {
		paths = append(paths, component.Path)
	}
	return paths
}

// placementPackages derives files_touched for architecture notes from where
// the architect placed components — the closest thing that document has to
// a file list.
func placementPackages(placements []ComponentPlacement) []string {
	packages := make([]string, 0, len(placements))
	for _, placement := range placements {
		packages = append(packages, placement.Package)
	}
	return packages
}
