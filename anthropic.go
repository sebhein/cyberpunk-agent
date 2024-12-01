package main

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func DebugLogger(req *http.Request, next option.MiddlewareNext) (res *http.Response, err error) {
	// Before the request

	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}
	log.Println("%+v", req.Header)
	log.Println("%+v", string(body[:]))

	// Forward the request to the next handler
	return next(req)
}

func createClient() *anthropic.Client {
	return anthropic.NewClient(
	// option.WithMiddleware(DebugLogger),
	// option.WithHeader("anthropic-beta", anthropic.AnthropicBetaPDFs2024_09_25),
	// option.WithHeaderAdd("anthropic-beta", anthropic.AnthropicBetaPromptCaching2024_07_31),
	)
}

func queryRules(client *anthropic.Client, query string) {
	// This is too expensive / requires too many tokens to make it usable
	// Need proper RAG

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	bytes, err := os.ReadFile(cwd + "/pdfs/cpr-rules.pdf")
	if err != nil {
		panic(err.Error())
	}

	b64PDF := base64.StdEncoding.EncodeToString(bytes)

	message, err := client.Beta.Messages.New(
		context.TODO(), anthropic.BetaMessageNewParams{
			Betas: anthropic.F([]string{
				anthropic.AnthropicBetaPDFs2024_09_25,
				anthropic.AnthropicBetaPromptCaching2024_07_31,
			}),
			System: anthropic.F([]anthropic.BetaTextBlockParam{
				{
					Text: anthropic.F("You are an AI 'agent' from the CyberPunk world, your responses should be concise and robotic. Do not be nice or kind. Only draw on information from CyberPunk RED."),
					Type: anthropic.F(anthropic.BetaTextBlockParamTypeText),
				},
			}),
			Model:     anthropic.F(anthropic.ModelClaude3_5Sonnet20241022),
			MaxTokens: anthropic.F(int64(1024)),
			Messages: anthropic.F([]anthropic.BetaMessageParam{
				{
					Role: anthropic.F(anthropic.BetaMessageParamRoleUser),
					Content: anthropic.F([]anthropic.BetaContentBlockParamUnion{
						anthropic.BetaBase64PDFBlockParam{
							Type: anthropic.F(anthropic.BetaBase64PDFBlockTypeDocument),
							Source: anthropic.F(anthropic.BetaBase64PDFSourceParam{
								Data:      anthropic.F(b64PDF),
								MediaType: anthropic.F(anthropic.BetaBase64PDFSourceMediaTypeApplicationPDF),
								Type:      anthropic.F(anthropic.BetaBase64PDFSourceTypeBase64),
							}),
							CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
								Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
							}),
						},
						anthropic.BetaTextBlockParam{
							Text: anthropic.F(query),
							Type: anthropic.F(anthropic.BetaTextBlockParamTypeText),
						},
					}),
				},
			},
			),
		},
	)

	if err != nil {
		panic(err.Error())
	}

	for idx := 0; idx < len(message.Content); idx++ {
		log.Println("%+v", message.Content[idx].Text)
	}

}

func queryAgent(client *anthropic.Client, query string) {
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model: anthropic.F(anthropic.ModelClaude3_5Sonnet20241022),
		System: anthropic.F([]anthropic.TextBlockParam{
			{
				Text: anthropic.F("You are an AI 'agent' from the CyberPunk world, your responses should be concise and robotic. Do not be nice or kind. Only draw on information from CyberPunk RED."),
				Type: anthropic.F(anthropic.TextBlockParamTypeText),
			},
		}),
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(query)),
		}),
	})

	if err != nil {
		panic(err.Error())
	}

	for idx := 0; idx < len(message.Content); idx++ {
		log.Println("%+v", message.Content[idx].Text)
	}
}
