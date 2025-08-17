#!/bin/bash
# build.sh - Railway build script

set -e

echo "Installing templ..."
go install github.com/a-h/templ/cmd/templ@latest

echo "Cleaning old generated files..."
find . -name '*_templ.go' -delete

echo "Generating templ files..."
templ generate

echo "Tidying Go modules..."
go mod tidy

echo "Installing npm dependencies..."
npm ci

echo "Building frontend assets..."
npm run build

echo "Building Go application..."
go build -o orbit2x .

echo "Build complete!"