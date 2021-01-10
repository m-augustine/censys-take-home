FROM golang:1.14.0 AS builder

#
WORKDIR /workspace
COPY . .

# Build
#
RUN CGO_ENABLED=0 go build -o bin/censys-take-home cmd/*

# --------------------------------------------------------------

FROM scratch

COPY --from=builder /workspace/bin/* /

ENTRYPOINT [ "/censys-take-home" ]
