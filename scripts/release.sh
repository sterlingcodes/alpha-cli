#!/bin/bash
set -e

# Alpha CLI Release Builder
# Builds binaries for all platforms and creates release artifacts

VERSION=${1:-$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")}
BINARY_NAME="alpha"
BUILD_DIR="./dist"

# Platforms to build
PLATFORMS=(
    "darwin_amd64"
    "darwin_arm64"
    "linux_amd64"
    "linux_arm64"
    "linux_arm"
    "windows_amd64"
)

echo "Building Alpha CLI $VERSION"
echo "================================"

# Clean previous builds
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# Build for each platform
for PLATFORM in "${PLATFORMS[@]}"; do
    OS="${PLATFORM%_*}"
    ARCH="${PLATFORM#*_}"

    OUTPUT_NAME="$BINARY_NAME"
    if [ "$OS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi

    echo "Building for $OS/$ARCH..."

    GOOS="$OS" GOARCH="$ARCH" go build \
        -ldflags "-s -w -X main.Version=$VERSION" \
        -o "$BUILD_DIR/$OUTPUT_NAME" \
        ./cmd/alpha

    # Create tarball (or zip for windows)
    cd "$BUILD_DIR"
    if [ "$OS" = "windows" ]; then
        zip "${BINARY_NAME}_${PLATFORM}.zip" "$OUTPUT_NAME"
        rm "$OUTPUT_NAME"
    else
        tar -czf "${BINARY_NAME}_${PLATFORM}.tar.gz" "$OUTPUT_NAME"
        rm "$OUTPUT_NAME"
    fi
    cd ..

    echo "  ${BINARY_NAME}_${PLATFORM}"
done

# Generate checksums
cd "$BUILD_DIR"
shasum -a 256 *.tar.gz *.zip 2>/dev/null > checksums.txt || shasum -a 256 *.tar.gz > checksums.txt
cd ..

echo ""
echo "================================"
echo "Release artifacts in $BUILD_DIR:"
ls -la "$BUILD_DIR"
echo ""
echo "To create a GitHub release:"
echo "  git tag -a $VERSION -m 'Release $VERSION'"
echo "  git push origin $VERSION"
echo "  gh release create $VERSION $BUILD_DIR/* --title 'Alpha CLI $VERSION' --notes 'Release notes here'"
