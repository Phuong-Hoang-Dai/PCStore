package model

// Paging controls pagination for search queries.
type Paging struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Process normalises paging values into safe bounds.
func (p *Paging) Process() {
	if p.Limit <= 0 {
		p.Limit = DefaultLimit
	}
	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}

// SearchRequest holds the parameters used to build an Elasticsearch query.
// An empty Query with no filters matches all products.
// type SearchRequest struct {
// 	Query    string // full-text query over name, brand and description
// 	Brand    string // exact brand filter
// 	Type     string // exact product type filter (standard/composite)
// 	CateID   int    // exact category id filter
// 	MinPrice float64
// 	MaxPrice float64
// 	Paging   Paging
// }

type SearchRequest struct {
	State       RequestState `json:"state"`
	QueryConfig QueryConfig  `json:"queryConfig"`
}

type RequestState struct {
	Current        int           `json:"current"`
	ResultsPerPage int           `json:"resultsPerPage"`
	SearchTerm     string        `json:"searchTerm"`
	Sort           SortValue     `json:"sort"`
	Filter         []FilterValue `json:"filter"`
}
type FilterValue struct {
	Name        string   `json:"name"`
	FilterValue []string `json:"filterValue"`
	Type        string   `json:"type"`
	Persistent  bool     `json:"persistent"`
}

type SortValue struct {
	Sort          string `json:"sort"`
	SortDirection int    `json:"sortDirection"`
}

type SearchResult struct {
	Total    int64         `json:"totalResults"`
	Products []Product     `json:"results"`
	Facets   FacetResponse `json:"facet"`
}

type AutoCompleteResult struct {
	Total    int64     `json:"totalResults"`
	Products []Product `json:"autocompletedResults"`
}

type FacetResponse struct {
	Value []FacetValueRes `json:"facetValue"`
	Range []FacetRangeRes `json:"facetRange"`
}

type FacetRangeRes struct {
	Range Range `json:"range"`
	Count int   `json:"count"`
}

type FacetValueRes struct {
	Field string `json:"field"`
	Count int    `json:"count"`
}

type QueryConfig struct {
	Facet FacetRange `json:"facet"`
}

type FacetRange struct {
	Range []Range `json:"range"`
}

type Range struct {
	From int `json:"from,omitempty"`
	To   int `json:"to,omitempty"`
}
