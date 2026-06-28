#!/bin/bash

set -e

echo "🧪 Running Unit Tests with Coverage (NO benchmarks, NO accuracy)..."

# Run ONLY unit tests (exclude benchmarks + specific accuracy test)
go test ./... \
  -coverprofile=coverage.out \
  -covermode=atomic \
  -run='^(Test|Example)'

echo ""
echo "📊 Generating HTML coverage report..."

go tool cover -html=coverage.out -o coverage.html

echo ""
echo "✅ Coverage report generated: coverage.html"
echo "📂 Open it in your browser to view coverage"