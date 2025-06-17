package api

import (
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

type WebSearchResult struct {
	Name         string  `json:"name"`
	Cover        string  `json:"cover"`
	ReviewCount  int     `json:"reviewCount"`
	ReviewRating float64 `json:"reviewRating"`
	URL          string  `json:"url"`
	Source       string  `json:"source"`
}

func (a *App) WebSearch(c echo.Context) error {
	ctx := c.Request().Context()
	query := c.QueryParam("q")

	// Remove leading and trailing whitespace from the query,
	query = strings.Trim(query, " ")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing query parameter 'q'")
	}

	// append " recipes" if the query does not end with "recipes" or "recipe"
	if !strings.HasSuffix(query, "recipes") && !strings.HasSuffix(query, "recipe") {
		query += " recipe"
	}

	searchResults, err := a.searchClient.Search(ctx, query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not search for recipes").SetInternal(err)
	}

	var results []WebSearchResult
	for _, r := range searchResults.Web.Results {
		if r.Subtype == "recipe" {
			parsedURL, err := url.Parse(r.URL)
			if err != nil {
				log.Error().Err(err).Msg("could not parse URL")
				continue
			}

			results = append(results, WebSearchResult{
				Name:         r.Recipe.Title,
				Cover:        r.Thumbnail.Src,
				ReviewCount:  r.Recipe.Rating.Count,
				ReviewRating: r.Recipe.Rating.Value,
				URL:          r.URL,
				Source:       parsedURL.Hostname(),
			})
		}
	}

	// Order the results by review count and then review rating.
	sort.Slice(results, func(i, j int) bool {
		return float64(results[i].ReviewCount) > float64(results[j].ReviewCount)
	})

	return c.JSON(http.StatusOK, results)
}

func (a *App) RecipeSearch(c echo.Context) error {
	ctx := c.Request().Context()

	input, err := prepareSearchRecipeInput(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	recipes, err := recipes.SearchRecipe(ctx, a.db, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not search for recipes").SetInternal(err)
	}

	return c.JSON(http.StatusOK, recipes)
}

func prepareSearchRecipeInput(c echo.Context) (recipes.SearchRecipeInput, error) {
	input := recipes.SearchRecipeInput{}

	// Required: Query parameter 'q'
	input.Query = strings.Trim(c.QueryParam("q"), " ")
	if input.Query == "" {
		return input, errors.New("missing query parameter 'q'")
	}

	// Optional: String query parameters
	input.KitchenID = null.NewString(c.QueryParam("kitchenId"), c.QueryParam("kitchenId") != "")
	input.Course = null.NewString(c.QueryParam("course"), c.QueryParam("course") != "")
	input.Class = null.NewString(c.QueryParam("class"), c.QueryParam("class") != "")
	input.Cuisine = null.NewString(c.QueryParam("cuisine"), c.QueryParam("cuisine") != "")
	input.OrderBy = null.NewString(c.QueryParam("sort"), c.QueryParam("sort") != "")

	// Optional: Query parameter 'difficulty', is int
	difficulty, err := web.ParseIntQueryParam(c, "difficulty", 0)
	if err != nil {
		return input, errors.New("invalid value for query parameter 'difficulty'")
	}
	input.MaxDifficulty = difficulty

	// Optional: Query parameter 'rating', is int
	rating, err := web.ParseIntQueryParam(c, "rating", 0)
	if err != nil {
		return input, errors.New("invalid value for query parameter 'rating'")
	}
	input.MinRating = rating

	// Optional: Query parameter 'time', is int
	time, err := web.ParseIntQueryParam(c, "time", 0)
	if err != nil {
		return input, errors.New("invalid value for query parameter 'time'")
	}
	input.MaxTime = time

	// Optional: Query parameter 'limit', is uint64, default to 20
	limit, err := web.ParseUintQueryParam(c, "limit", 20)
	if err != nil {
		return input, err
	}
	input.Limit = limit

	// Optional: Query parameter 'offset', is uint64, default to 0
	offset, err := web.ParseUintQueryParam(c, "offset", 0)
	if err != nil {
		return input, err
	}
	input.Offset = offset

	return input, nil
}

type RecipeSearchFiltersResponse struct {
	Courses    []string `json:"courses"`
	Classes    []string `json:"classes"`
	Difficulty []int    `json:"difficulty"`
	Rating     []int    `json:"rating"`
	Sort       []string `json:"sort"`
}

func (a *App) RecipeSearchFilters(c echo.Context) error {
	filters := RecipeSearchFiltersResponse{
		Courses:    stringMapKeys(recipes.ValidCourses),
		Classes:    stringMapKeys(recipes.ValidClasses),
		Difficulty: []int{1, 2, 3, 4, 5},
		Rating:     []int{1, 2, 3, 4, 5},
		Sort:       stringMapKeys(recipes.ValidSort),
	}

	return c.JSON(http.StatusOK, filters)
}

func stringMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
