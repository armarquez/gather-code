# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-gather-code"
LABEL REPO="https://github.com/armarquez/gather-code"

ENV PROJPATH=/go/src/github.com/armarquez/gather-code

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/armarquez/gather-code
WORKDIR /go/src/github.com/armarquez/gather-code

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/armarquez/gather-code"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/gather-code/bin

WORKDIR /opt/gather-code/bin

COPY --from=build-stage /go/src/github.com/armarquez/gather-code/bin/gather-code /opt/gather-code/bin/
RUN chmod +x /opt/gather-code/bin/gather-code

# Create appuser
RUN adduser -D -g '' gather-code
USER gather-code

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/gather-code/bin/gather-code"]
