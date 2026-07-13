package cli

type Item struct {
	Label    string
	Disabled bool
}

func itemLabels(items []Item) []string {
	labels := make([]string, len(items))
	for i, item := range items {
		if item.Disabled {
			labels[i] = item.Label + " (уже есть)"
		} else {
			labels[i] = item.Label
		}
	}
	return labels
}
