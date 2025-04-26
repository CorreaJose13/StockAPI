package ai

import (
	"context"
	"errors"
	"fmt"

	"github.com/CorreaJose13/StockAPI/config"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	errNoCandidates = errors.New("no candidates found")
	errNoResponse   = errors.New("no response received from Gemini")
)

type GeminiRepository struct {
	gClient *genai.Client
}

func NewGeminiRepository(ctx context.Context, cfg *config.Config) (*GeminiRepository, error) {

	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiRepository{
		gClient: client,
	}, nil
}

func (g *GeminiRepository) CloseClient() error {
	return g.gClient.Close()
}

func extractFirstCandidate(candidates []*genai.Candidate) (string, error) {
	for _, candidate := range candidates {
		content := candidate.Content
		if content == nil || len(content.Parts) == 0 {
			continue
		}

		return fmt.Sprintf("%v", content.Parts[0]), nil
	}

	return "", errNoCandidates
}

func (g *GeminiRepository) SendPrompt(ctx context.Context, prompt string) (string, error) {

	model := g.gClient.GenerativeModel("gemini-2.0-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if resp == nil {
		return "", errNoResponse
	}

	candidates := resp.Candidates

	return extractFirstCandidate(candidates)
}
