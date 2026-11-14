FROM golang:1.27.4-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main ./cmd/api

FROM alpine:3.22 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE 8080
CMD ["./main"]

FROM node:22-alpine AS frontend_build
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

FROM node:22-alpine AS frontend
RUN npm install -g serve
COPY --from=frontend_build /frontend/dist /app/dist
EXPOSE 5173
CMD ["serve", "-s", "/app/dist", "-l", "5173"]
