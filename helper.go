package main

import (
	"strconv"
)

// MustStringToInt converts a string to an int
func (a *App) MustStringToInt(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		Log.Logger.Warn().Str("input", s).Str("error", err.Error()).Msg("Unable to convert string to integer.")
	}
	return i
}

// MustStringToFloat converts a string to a float64
func (a *App) MustStringToFloat(s string) float64 {
	if s == "" {
		return 0.0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		Log.Logger.Warn().Str("input", s).Str("error", err.Error()).Msg("Unable to convert string to float.")
	}
	return f
}
