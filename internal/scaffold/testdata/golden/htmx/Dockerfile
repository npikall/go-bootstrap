FROM golang:1.27.4-alpine AS build
RUN apk add --no-cache libstdc++ libgcc
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go run "github.com/a-h/templ/cmd/templ@$(go list -m github.com/a-h/templ | cut -d' ' -f2)" generate
ARG TARGETARCH
RUN wget -qO /usr/local/bin/tailwindcss "https://github.com/tailwindlabs/tailwindcss/releases/download/v4.3.5/tailwindcss-linux-$(echo "${TARGETARCH:-amd64}" | sed 's/amd64/x64/')-musl" \
    && chmod +x /usr/local/bin/tailwindcss \
    && tailwindcss -i web/styles/input.css -o web/assets/css/output.css --minify
RUN CGO_ENABLED=0 go build -o main ./cmd/api

FROM alpine:3.22 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE 8080
CMD ["./main"]
