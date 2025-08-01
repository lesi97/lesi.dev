package bungie_store

import "strconv"

func getKillCounts(objectivesPerPlug map[string][]plugObjective) (int, string) {
	categories := map[string][]int{
		"PVP":    {38912240, 231101171, 3244015567, 2285636663},
		"PVE":    {2240097604, 231101171, 3915764595, 3915764594, 3915764593, 1187045864, 1690059054, 16638393, 1124054883, 3624435060, 16638392, 2617715132, 2617715133, 2302094943, 905869860},
		"Trials": {3915764595},
	}

	for label, ids := range categories {
		for _, id := range ids {
			key := strconv.Itoa(id)
			if objectives, ok := objectivesPerPlug[key]; ok && len(objectives) > 0 {
				return objectives[0].Progress, label
			}
		}
	}

	return 0, ""
}
