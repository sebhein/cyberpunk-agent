package main

func main() {
	anthropicClient := createClient()
	// queryRules(anthropicClient, "How does movement work?")
	queryAgent(anthropicClient, "Respond in the affirmative if this is working")
}
