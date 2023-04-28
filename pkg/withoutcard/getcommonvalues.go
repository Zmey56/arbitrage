package withoutcard

func GetCommonValue(arr1, arr2 []string) []string {
	common := []string{}

	for _, v1 := range arr1 {
		for _, v2 := range arr2 {
			if v1 == v2 {
				common = append(common, v1)
				break
			}
		}
	}

	return common
}
