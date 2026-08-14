#!/bin/bash
set -e
export ANDROID_HOME=${ANDROID_HOME:-/root/Android/Sdk}
export PATH="$HOME/go/bin:$PATH"
cd "$(dirname "$0")/../go/runtime"
echo "Building Go Mobile AAR..."
gomobile bind -target=android -androidapi 24 -o /tmp/libdsh.aar ./mobile
echo "Extracting .so files..."
mkdir -p ../android/app/src/main/jniLibs
unzip -o /tmp/libdsh.aar -d /tmp/libdsh_extracted
cp -r /tmp/libdsh_extracted/jni/* ../android/app/src/main/jniLibs/
echo "Done!"
