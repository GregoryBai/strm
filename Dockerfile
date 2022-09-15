FROM golang:latest
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./out/dist ./cmd/

# EXPOSE 8000 ?
# ENTRYPOINT ["build/out/dist"] ?

CMD ./out/dist
