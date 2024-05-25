package utils

func FindIntersection(arr1, arr2 []string) []string {
	intersection := make([]string, 0)
	set := make(map[string]bool)

	for _, num := range arr1 {
		set[num] = true
	}

	for _, num := range arr2 {
		if set[num] {
			intersection = append(intersection, num)
		}
	}

	return intersection
}
