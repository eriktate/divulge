FROM golang:1.14.2-buster AS builder
COPY ./ /opt/divulge
WORKDIR /opt/divulge
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o divulge cmd/server/main.go

FROM scratch
COPY --from=builder /opt/divulge/divulge divulge
ENV DIVULGE_HOST=0.0.0.0
EXPOSE 8080
ENTRYPOINT ["./divulge"]

