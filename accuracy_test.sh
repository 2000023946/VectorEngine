#!/bin/bash

set -e

echo "📊 Running VectorEngine Accuracy Tests (Recall@K)..."

# -----------------------------
# OUTPUT FILES
# -----------------------------
RESULTS_TXT="accuracy_results.txt"
RESULTS_HTML="accuracy_results.html"

# -----------------------------
# RUN TESTS
# -----------------------------
go test ./... \
  -v \
  -run "^TestAccuracy" \
  | tee $RESULTS_TXT

echo ""
echo "📦 Generating HTML report..."

# -----------------------------
# SIMPLE HTML WRAPPER
# -----------------------------
cat <<EOF > $RESULTS_HTML
<html>
<head>
    <title>VectorEngine Accuracy Report</title>
    <style>
        body { font-family: monospace; background: #111; color: #0f0; padding: 20px; }
        pre { background: #000; padding: 15px; border-radius: 8px; }
    </style>
</head>
<body>
    <h2>📊 VectorEngine Accuracy Results</h2>
    <pre>
EOF

# inject test output safely
cat $RESULTS_TXT >> $RESULTS_HTML

cat <<EOF >> $RESULTS_HTML
    </pre>
</body>
</html>
EOF

echo ""
echo "🎯 Accuracy test complete"
echo "📄 Text report: $RESULTS_TXT"
echo "🌐 HTML report: $RESULTS_HTML"
echo "Openning accuracy_results.html ..."
open accuracy_results.html