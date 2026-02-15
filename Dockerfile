# -------- Stage 1: Build --------
FROM golang:1.24.5-bookworm
ENV CGO_ENABLED=1
ENV GOOS=linux

WORKDIR /app
COPY . .

RUN go get
RUN go build -o bin .

ENTRYPOINT ["/app/bin"]