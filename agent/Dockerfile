FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go test ./... && \
    go build -o /out/agent .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/agent /usr/local/bin/agent
ENV AGENT_SHARED_SECRET="" # set in compose
EXPOSE 8080
ENTRYPOINT ["agent"]
