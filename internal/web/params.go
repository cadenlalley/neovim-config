package web

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ParseUintQueryParam parses a query parameter as uint64 with a default value
func ParseUintQueryParam(c echo.Context, paramName string, defaultValue uint64) (uint64, error) {
	if value := c.QueryParam(paramName); value != "" {
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return 0, errors.New("invalid value for query parameter '" + paramName + "'")
		}
		return parsed, nil
	}
	return defaultValue, nil
}

// ParseIntQueryParam parses a query parameter as int with a default value
func ParseIntQueryParam(c echo.Context, paramName string, defaultValue int) (int, error) {
	if value := c.QueryParam(paramName); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid value for query parameter '" + paramName + "'")
		}
		return parsed, nil
	}
	return defaultValue, nil
}
