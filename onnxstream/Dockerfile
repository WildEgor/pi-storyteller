# FROM arm32v7/debian:latest as base
FROM ubuntu:jammy as base
RUN apt update && \
    apt install -y git git-lfs python3 cmake build-essential curl coreutils && \
    git lfs install

FROM base as XNNPACK
WORKDIR /XNNPACK
RUN curl -L -o /XNNPACK/xnnpack.tar.gz https://api.github.com/repos/google/XNNPACK/tarball && \
    tar -xzf /XNNPACK/xnnpack.tar.gz -C /XNNPACK --strip-components=1 && \
    rm /XNNPACK/xnnpack.tar.gz
RUN cd /XNNPACK && \
    if [ -d .git ]; then \
        git checkout 1c8ee1b68f3a3e0847ec3c53c186c5909fa3fbd3; \
    else \
        echo "No .git directory found. Skipping git checkout."; \
    fi
WORKDIR /XNNPACK/build
RUN cmake -DXNNPACK_BUILD_TESTS=OFF -DXNNPACK_BUILD_BENCHMARKS=OFF .. && \
    cmake --build . --config Release

FROM base as OnnxStream
COPY --from=XNNPACK /XNNPACK /XNNPACK
WORKDIR /OnnxStream
RUN git clone --verbose https://github.com/vitoplantamura/OnnxStream.git .
WORKDIR /OnnxStream/src/build
RUN cmake -DMAX_SPEED=ON -DXNNPACK_DIR=/XNNPACK .. && \
    cmake --build . -j --config Release

# lightweight model for testing
FROM alpine:3.18 as weights
RUN apk update && \
    apk add git git-lfs && \
    git lfs install && \
    mkdir -p /weights && \
    git clone --depth=1 https://huggingface.co/AeroX2/stable-diffusion-xl-turbo-1.0-onnxstream /weights

FROM base as app
# COPY --from=weights /weights /app/weights
COPY --from=OnnxStream /OnnxStream/src/build/sd /app/OnnxStream/src/build/sd
ENTRYPOINT [ "/app/scripts/generate.sh" ]