FROM golang:1.26.2-bookworm

WORKDIR /workspace
ENV GOTOOLCHAIN=local

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build ./... && go build -o /usr/local/bin/recipe-planner ./cmd/server

EXPOSE 8080
CMD ["/usr/local/bin/recipe-planner"]
