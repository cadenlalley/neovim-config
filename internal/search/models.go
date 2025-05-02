package search

//
// Models for Brave API, this is a subset of what's available.
// https://api-dashboard.search.brave.com/app/documentation/web-search/responses#MixedResponse
//

type SearchResult struct {
	Query Query    `json:"query"`
	Web   WebItems `json:"web"`
	Type  string   `json:"type"`
}

type WebItems struct {
	Type    string    `json:"type"`
	Results []WebItem `json:"results"`
}

type Query struct {
	Original string `json:"original"`
}

type WebItem struct {
	Title          string    `json:"title"`
	URL            string    `json:"url"`
	PageAge        string    `json:"page_age"`
	Language       string    `json:"language"`
	FamilyFriendly bool      `json:"family_friendly"`
	Type           string    `json:"type"`
	Subtype        string    `json:"subtype"`
	Thumbnail      Thumbnail `json:"thumbnail"`
	Recipe         Recipe    `json:"recipe"`
}

type Thumbnail struct {
	Src      string `json:"src"`
	Original string `json:"original"`
	Logo     bool   `json:"logo"`
}

type Recipe struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Rating      RecipeRating `json:"rating"`
}

type RecipeRating struct {
	Value float64 `json:"ratingValue"`
	Count int     `json:"reviewCount"`
}
