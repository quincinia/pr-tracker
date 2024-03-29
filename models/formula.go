// formula defines methods for calculating the point distribution among a tournament's attendees
package models

import "pr-tracker/requests"

var (
	up_bonus    int = 5
	att_bonus   int = 10
	first_bonus int = 10
	br_bonus    int = 5
)

// For one-off calculations where you don't need to calculate the point distribution for all attendees
func (t *Tournament) PointsFromPlacement(placement int) (points int) {
	index := placementIndex(placement)
	if index+1 > t.UniquePlacings {
		return -1
	}

	points += up_bonus * (t.UniquePlacings - 1 - index)
	points += att_bonus

	if placement == 1 {
		points += first_bonus
	}

	if placement == 2 && t.BracketReset {
		points += br_bonus
	}

	if t.Tier != nil {
		points *= t.Tier.Multiplier
	}

	return
}

// Generates a mapping from a player's final placement to the points that placement earns
func (t *Tournament) PointMap() (pm map[int]int) {
	placings := placingsList(t.UniquePlacings)
	pm = make(map[int]int)

	for i, p := range placings {
		// Each unique placement earns you 5 points (unless you got last)
		// eg. in a tournament of 5 players, the point distribution looks like this:
		// 1st 2nd 3rd 4th 5th
		// 20p 15p 10p  5p  0p
		pm[p] = up_bonus * (t.UniquePlacings - 1 - i)

		// Each attendee gets a bonus just for showing up
		// Attendance bonus set to 10 points
		pm[p] += att_bonus
	}

	// The player in 1st place gets a bonus for winning
	pm[1] += first_bonus

	// The player who made a bracket reset, but didn't end up winning also gets a bonus
	if t.BracketReset {
		pm[2] += br_bonus
	}

	// After the raw scores are calculated, the tournament multiplier will be applied
	if t.Tier != nil {
		for _, p := range placings {
			pm[p] *= t.Tier.Multiplier
		}
	}

	return
}

// Generates a list of the actual unique placings from the known number of unique placings
func placingsList(up int) (list []int) {
	if up < 1 {
		panic("unsupported number of unique placings")
	}

	var (
		placing int = 3
		shared  int = 1
		update  bool
	)
	list = make([]int, up)

	// Panic check ensures length is at least 1
	list[0] = 1

	if up == 1 {
		return
	}

	// The first two elements don't follow the pattern, so manually putting the values in
	list[1] = 2

	// Pattern holds for all values > 2
	for i := 2; i < up; i++ {
		// Add the placing number
		list[i] = placing

		// Move to the next placing number
		placing += shared

		// The number of shared placings doubles every 2 unique placings
		if update {
			shared *= 2
		}

		// Using a bool to track every second placement
		update = !update
	}

	return
}

// Returns the number of placements that are "ahead" of the given one
// Can also be thought of as the placement's "index" in the placement list (see above)
// If the placement is invalid, then returns -1
func placementIndex(placement int) (index int) {
	if placement < 1 {
		return -1
	}

	if placement <= 4 {
		return placement - 1
	}

	exp, value := requests.TruncLog2(placement)

	if placement == value+1 {
		return 2 * exp
	}

	if placement == (3*value/2)+1 {
		return 2*exp + 1
	}

	return -1
}
