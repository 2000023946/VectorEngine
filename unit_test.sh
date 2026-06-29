#!/bin/bash

set -e

echo "🧪 Running Unit Tests with Coverage (NO benchmarks, NO accuracy)..."

# Run ONLY unit tests and generate coverage file
go test ./... -run "^TestUnit" -coverprofile=coverage.out

echo ""
echo "📊 Generating HTML coverage report..."

go tool cover -html=coverage.out -o coverage.html

echo ""
echo "📦 Coverage Summary:"
go tool cover -func=coverage.out

echo ""
echo "✅ Coverage report generated: coverage.html"
echo "📂 Opening coverage.html file ..."
open coverage.html