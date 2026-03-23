FROM golang:1.25 AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 \
  GOOS=$TARGETOS \
  GOARCH=$TARGETARCH \
  go build -o /out/testapp ./main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=build /out/testapp /app/testapp

EXPOSE 8080
ENTRYPOINT ["/app/testapp"]
