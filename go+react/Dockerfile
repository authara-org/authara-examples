FROM node:20-alpine AS frontend-builder
WORKDIR /app

COPY web/package.json web/package-lock.json* ./
RUN npm install

COPY web .
RUN npm run build

FROM golang:1.25 AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /app/dist ./web/dist

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/testapp ./main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=backend-builder /out/testapp /app/testapp
COPY --from=backend-builder /app/web/dist /app/web/dist

EXPOSE 8080
ENTRYPOINT ["/app/testapp"]
