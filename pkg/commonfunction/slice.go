package commonfunction

func CommonElement(list1, list2 []string) []string {

	var matches []string

	for _, item1 := range list1 {
		for _, item2 := range list2 {
			if item1 == item2 {
				matches = append(matches, item1)
			}
		}
	}

	return matches

}
