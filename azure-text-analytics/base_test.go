package function

var mockTextAnalyticsService *MockTextAnalyticsService

func createLanguageResult() *LanguageResult {
	return &LanguageResult{
		Name:        "English",
		ISO6391Name: "en",
		Score:       1.0,
	}
}
