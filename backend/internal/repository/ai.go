package repository

import "context"

type AIRepository interface {
	SendPrompt(ctx context.Context, prompt string) (string, error)
	CloseClient() error
}

var aiRepoImpl AIRepository

func SetAIRepository(repo AIRepository) {
	aiRepoImpl = repo
}

func SendPrompt(ctx context.Context, prompt string) (string, error) {
	return aiRepoImpl.SendPrompt(ctx, prompt)
}

func CloseClient() error {
	return aiRepoImpl.CloseClient()
}
