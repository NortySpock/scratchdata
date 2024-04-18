FROM golang:1.22.1 as builder
WORKDIR /build
# avoiding doing "COPY . ."
# as that could bake the "data" folder into the image if the operator uses the default config
VOLUME ["/data"]
VOLUME ["/config.yaml"]
COPY go.mod go.sum ./
COPY models ./
COPY pkg ./
RUN go mod download
RUN go build -o /scratchdata
RUN chmod +x /scratchdata
EXPOSE 8080
ENTRYPOINT ["/scratchdata", "config.yaml"]
