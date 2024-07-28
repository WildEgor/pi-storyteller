#!/bin/bash

apt-get install -y cmake build-essential curl libxml2-utils coreutils git git-lfs && git lfs install

git clone https://github.com/google/XNNPACK.git /app/XNNPACK
cd /app/XNNPACK
git checkout 579de32260742a24166ecd13213d2e60af862675

mkdir /app/XNNPACK/build
cd /app/XNNPACK/build
cmake -DXNNPACK_BUILD_TESTS=OFF -DXNNPACK_BUILD_BENCHMARKS=OFF ..
cmake --build . --config Release

git clone https://github.com/vitoplantamura/OnnxStream.git /app/OnnxStream
mkdir /app/OnnxStream/src/build
cd /app/OnnxStream/src/build
cmake -DMAX_SPEED=ON -DXNNPACK_DIR=/app/XNNPACK ..
cmake --build . -j --config Release

mkdir -p /app/weights 
cd /app/weights
git clone --depth=1 https://huggingface.co/AeroX2/stable-diffusion-xl-turbo-1.0-onnxstream .

cd /
mkdir -p /app/bin
mkdir -p /app/out

chmod +x /app/OnnxStream/src/build/sd