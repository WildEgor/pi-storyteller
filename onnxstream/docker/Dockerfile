FROM ubuntu:22.04 as deps
RUN apt update && \
    apt install -y git cmake build-essential curl

FROM deps as build
WORKDIR /onnx/XNNPACK
RUN git clone https://github.com/google/XNNPACK.git . && \
    git checkout 579de32260742a24166ecd13213d2e60af862675 && \
    mkdir build
WORKDIR /onnx/XNNPACK/build
RUN cmake -DXNNPACK_BUILD_TESTS=OFF -DXNNPACK_BUILD_BENCHMARKS=OFF .. && \
    cmake --build . --config Release
WORKDIR /onnx/OnnxStream
RUN git clone https://github.com/vitoplantamura/OnnxStream.git . && \
    mkdir -p /onnx/OnnxStream/src/build
WORKDIR /onnx/OnnxStream/src/build
RUN cmake -DMAX_SPEED=ON -DXNNPACK_DIR=/onnx/XNNPACK .. && \
    cmake --build . -j --config Release

# lightweight model for testing
FROM alpine:3.18 as weights
RUN apk update && \
    apk add git git-lfs && \
    git lfs install && \
    mkdir -p /weights && \
    git clone --depth=1 https://huggingface.co/AeroX2/stable-diffusion-xl-turbo-1.0-onnxstream /weights

FROM ubuntu:22.04 as app
RUN apt update && \
    apt install -y curl libxml2-utils coreutils

COPY --from=build /onnx/OnnxStream/src/build/sd /app/OnnxStream/src/build/sd
COPY --from=build /onnx/OnnxStream/src/build/sd /app/bin/sd
# COPY --from=weights /weights /app/weights

ENTRYPOINT [ "/app/scripts/run.sh" ]
