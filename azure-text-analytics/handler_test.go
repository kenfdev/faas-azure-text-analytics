package function

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

func setupService() {
	mockTextAnalyticsService = &MockTextAnalyticsService{}
}

func TestHandle_FetchesLanguageIfNotPresentInRequest(t *testing.T) {
	// Arrange
	setupService()

	request := &AnalyzeRequest{
		Language: "",
		Text:     "Text to Request",
	}

	expectedLanguage := "en"
	mockTextAnalyticsService.On("FetchLanguage", request.Text).Return(&LanguageResult{ISO6391Name: expectedLanguage}, nil)

	expectedKeyPhrases := []string{"phrase1", "phrase2"}
	mockTextAnalyticsService.On("FetchKeyPhrases", mock.Anything).Return(expectedKeyPhrases, nil)

	expectedSentimentScore := 0.5
	mockTextAnalyticsService.On("FetchSentiment", mock.Anything).Return(expectedSentimentScore, nil)

	// Act
	response, _ := internalHandle(request, mockTextAnalyticsService)

	// Assert
	mockTextAnalyticsService.AssertExpectations(t)
	assert.Equal(t, expectedLanguage, response.Result.Language)
	assert.Equal(t, expectedKeyPhrases, response.Result.KeyPhrases)
	assert.Equal(t, expectedSentimentScore, response.Result.SentimentScore)
}

func TestHandle_DoesNotFetchLanguageIfLanguagePresentInRequest(t *testing.T) {
	// Arrange
	setupService()
	request := &AnalyzeRequest{
		Language: "en",
		Text:     "Text to Request",
	}

	expectedKeyPhrases := []string{"phrase1", "phrase2"}
	mockTextAnalyticsService.On("FetchKeyPhrases", mock.Anything).Return(expectedKeyPhrases, nil)

	expectedSentimentScore := 0.5
	mockTextAnalyticsService.On("FetchSentiment", mock.Anything).Return(expectedSentimentScore, nil)

	// Act
	response, _ := internalHandle(request, mockTextAnalyticsService)

	// Assert
	mockTextAnalyticsService.AssertNotCalled(t, "FetchLanguage", mock.Anything)
	assert.Equal(t, request.Language, response.Result.Language)
	assert.Equal(t, expectedKeyPhrases, response.Result.KeyPhrases)
	assert.Equal(t, expectedSentimentScore, response.Result.SentimentScore)
}
