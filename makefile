build:
	echo "Building KubeLinterBot"
	go build cmd/kube-linter-bot/kube-linter-bot.go cmd/kube-linter-bot/server.go cmd/kube-linter-bot/utilities.go