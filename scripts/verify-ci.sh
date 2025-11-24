#!/bin/bash
set -e

echo "🚀 Starting Local CI Verification..."

echo "--------------------------------------------------"
echo "📦 Backend Verification"
echo "--------------------------------------------------"

echo "🔍 Verifying Go modules..."
go mod verify

echo "🧪 Running Go tests..."
for dir in services/* pkg/*; do
  if [ -d "$dir" ] && [ -f "$dir/go.mod" ]; then
    echo "   Testing $(basename $dir)..."
    (cd "$dir" && go test ./...)
  fi
done

echo "🏗️  Building Services..."
for dir in services/*; do
  if [ -d "$dir" ]; then
    echo "   Building $(basename $dir)..."
    (cd "$dir" && go build -o /dev/null ./...)
  fi
done

echo "✅ Backend Verification Passed!"

echo "--------------------------------------------------"
echo "🎨 Frontend Verification"
echo "--------------------------------------------------"

cd frontend

echo "📦 Installing dependencies..."
npm install

echo "🧹 Linting..."
npm run lint

echo "🏗️  Building Frontend..."
npm run build

echo "✅ Frontend Verification Passed!"

echo "--------------------------------------------------"
echo "🎉 All CI checks passed successfully!"
echo "--------------------------------------------------"
