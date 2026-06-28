#!/bin/bash

set -e

echo "📊 Running VectorEngine Accuracy Test (Recall@K)..."

# Run ONLY the accuracy test
go test ./... \
  -v \
  -run=TestGraphRecallAtK

echo ""
echo "🎯 Accuracy test complete"