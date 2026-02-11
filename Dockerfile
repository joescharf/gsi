FROM alpine:latest
RUN apk add --no-cache ca-certificates git
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/gsi /usr/local/bin/gsi
ENTRYPOINT ["gsi"]
