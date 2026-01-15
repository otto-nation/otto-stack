package services

// CategoryDisplayInfo maps categories to display information
var CategoryDisplayInfo = map[string]struct {
	Name string
	Icon string
}{
	CategoryDatabase:      {"Database", "📊"},
	CategoryCache:         {"Cache", "💾"},
	CategoryMessaging:     {"Messaging", "📨"},
	CategoryObservability: {"Observability", "🔍"},
	CategoryCloud:         {"Cloud", "☁️"},
}
