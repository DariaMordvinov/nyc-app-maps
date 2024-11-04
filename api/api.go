package api

import (
	"context"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenerateContent() (*genai.GenerateContentResponse, error) {
	ctx := context.Background()
	key := os.Getenv("GEMINI_KEY")
	if key != "" {
		client, err := genai.NewClient(ctx, option.WithAPIKey(key))
		if err != nil {
			return nil, err
		}
		defer client.Close()
	
		model := client.GenerativeModel("gemini-1.5-flash")
		model.ResponseMIMEType = "application/json"
		resp, err := model.GenerateContent(ctx, genai.Text("JSON array of the most iconic nyc buildings in the style of Colonial/Neo-Colonial. For each building coordinates, title, founded year and 2 paragraphs of the text (combined in one property)"))
		if err != nil {
			return nil, err
		}
	
		return resp, nil

	} else {
		return nil, nil
	}
}
