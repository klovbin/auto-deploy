package cli

type Item struct {
	Label        string
	Disabled     bool
	DisabledText string
}

func itemLabels(items []Item) []string {
	labels := make([]string, len(items))
	for i, item := range items {
		if item.Disabled {
			text := item.DisabledText
			if text == "" {
				text = "недоступно"
			}
			labels[i] = item.Label + " (" + text + ")"
		} else {
			labels[i] = item.Label
		}
	}
	return labels
}
