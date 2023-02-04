package main

const (
	tgBotHost = "api.telegram.org"
)


func main() {
	tgClient = telegram.New(mustToken())

	// fetcher = fetcher.New()

	// processor = processor.New()

	// consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String(
		name: "token-bot-token", 
		value: "", 
		usage: "token for access to tg bot")

	flag.Parse()

	if *token == "" {
		log.Fatal(v...: "token is not specified")
	}

	return *token
}