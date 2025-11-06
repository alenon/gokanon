//go:build !race

package interactive

// isRaceEnabled is false when tests are run without -race flag
const isRaceEnabled = false
