package function

// structs

type AnalyzeRequest struct {
	Language string
	Text     string
}

type AnalyzeResult struct {
	Language       string   `json:"language"`
	SentimentScore float64  `json:"sentiment_score"`
	KeyPhrases     []string `json:"key_phrases"`
}

type AnalyzeResponse struct {
	Result       *AnalyzeResult
	ErrorMessage string
	StatusCode   int
}

type LanguageResult struct {
	Name        string  `json:"name"`
	ISO6391Name string  `json:"iso6391name"`
	Score       float64 `json:"score"`
}

type ErrorContext struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error *ErrorContext `json:"error"`
}

// interfaces

type TextAnalyticsService interface {
	FetchLanguage(text string) (*LanguageResult, error)
	FetchKeyPhrases(req *AnalyzeRequest) ([]string, error)
	FetchSentiment(req *AnalyzeRequest) (float64, error)
}
