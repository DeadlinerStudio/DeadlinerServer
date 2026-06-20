package state

import "time"

type LogicalVersion struct {
	TS  string `json:"ts"`
	Ctr int32  `json:"ctr"`
	Dev string `json:"dev"`
}

func CompareLogicalVersion(left, right LogicalVersion) (int, error) {
	leftTime, err := time.Parse(time.RFC3339, left.TS)
	if err != nil {
		return 0, err
	}

	rightTime, err := time.Parse(time.RFC3339, right.TS)
	if err != nil {
		return 0, err
	}

	switch {
	case leftTime.Before(rightTime):
		return -1, nil
	case leftTime.After(rightTime):
		return 1, nil
	case left.Ctr < right.Ctr:
		return -1, nil
	case left.Ctr > right.Ctr:
		return 1, nil
	case left.Dev < right.Dev:
		return -1, nil
	case left.Dev > right.Dev:
		return 1, nil
	default:
		return 0, nil
	}
}
