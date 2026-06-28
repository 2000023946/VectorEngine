#!/bin/bash

set -e

echo "🚀 Running VectorEngine Benchmarks..."

go test ./... -bench=. -benchmem -run=^$ -count=1

echo ""
echo "📊 Benchmark complete"