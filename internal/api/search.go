package api

import (
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
	query := c.QueryParam("q")
	kitchenID := c.QueryParam("kitchenId")

	// Remove leading and trailing whitespace from the query,
	query = strings.Trim(query, " ")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing query parameter 'q'")
	}

	recipes, err := recipes.SearchRecipe(ctx, a.db, recipes.SearchRecipeInput{
		Query:     query,
		KitchenID: kitchenID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not search for recipes").SetInternal(err)
	}

	return c.JSON(http.StatusOK, recipes)
}
