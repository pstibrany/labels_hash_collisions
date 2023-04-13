echo "!stringlabels"
echo

go run ./cmd/check

echo "stringlabels"
echo

go run -tags=stringlabels ./cmd/check
