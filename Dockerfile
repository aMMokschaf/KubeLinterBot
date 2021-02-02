FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /cmd/kube-linter-bot/

# Copy and download dependency using go mod
#COPY go.mod .
#COPY go.sum .
#RUN go mod download

# Copy the code into the container
COPY . .

RUN go build -o kube-linter-bot ./cmd/kubelinterbot/kube-linter-bot.go

EXPOSE 4567

CMD [ "./kube-linter-bot" ]