FROM golang:1.14.0 AS builder

#
WORKDIR /workspace
COPY . .

# Build
#
RUN CGO_ENABLED=0 go build -o bin/ipgeolocator cmd/*

# --------------------------------------------------------------

FROM scratch

COPY --from=builder /workspace/bin/* /

ENTRYPOINT [ "/ipgeolocator" ]
