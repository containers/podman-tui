package system

func unique(evtypes []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range evtypes {
		if _, value := keys[entry]; !value {
			keys[entry] = true

			list = append(list, entry)
		}
	}

	return list
}
