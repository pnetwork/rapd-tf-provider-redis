# Build the manager binary
FROM cr.pentium.network/golang:1.21-bookworm as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

# Download Terraform binary with specific platform
RUN curl -s https://packagecloud.io/install/repositories/opentofu/tofu/script.deb.sh?any=true | bash \
    && apt-get install tofu=1.6.0

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY internal/ internal/

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o mavis-auto-po cmd/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM cr-preview.pentium.network/distroless/static-debian12:debug-nonroot

WORKDIR /home/nonroot

# Copy GO & Terraform binary from builder stage
COPY --from=builder --chown=65532:65532 /workspace/mavis-auto-po /home/nonroot/mavis-auto-po
COPY --from=builder --chown=65532:65532 /bin/tofu /home/nonroot/.tofu/bin/tofu

# Add Terraform binary to Path
ENV PATH="/home/nonroot/.tofu/bin:${PATH}"

USER 65532:65532

ENTRYPOINT ["/home/nonroot/mavis-auto-po"]
