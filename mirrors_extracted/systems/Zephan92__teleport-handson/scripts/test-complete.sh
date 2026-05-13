#!/bin/bash
set -e

COVERAGE_DIR=${1:-coverage}

echo "1. Cleaning previous coverage..."
rm -rf "$COVERAGE_DIR"
mkdir -p "$COVERAGE_DIR"

echo "2. Running tests for all packages..."
# Filter out api/v1 as per Makefile logic
PACKAGES=$(go list -buildvcs=false ./... | grep -v /api/v1)

# Define coverage scope: All packages except generated api/v1
COVER_PKGS=$(go list -buildvcs=false ./... | grep -v /api/v1 | paste -sd,)

echo "=== User Mode Tests ==="
for pkg in $PACKAGES; do
    # Replace slashes and dots with underscores for file names
    pkg_slug=$(echo "$pkg" | sed 's|[/.]|_|g')
    
    # Run user tests; allow failure (|| true) to continue collecting coverage
    go test -buildvcs=false -coverpkg="$COVER_PKGS" -coverprofile="$COVERAGE_DIR/coverage-user-$pkg_slug.out" "$pkg" 2>&1 | grep -v "warning: no packages being tested depend on matches" | sed 's/ of statements in .*/ of statements/' || true
done

echo ""
echo "=== Root Mode Tests ==="
for pkg in $PACKAGES; do
    # Replace slashes and dots with underscores for file names
    pkg_slug=$(echo "$pkg" | sed 's|[/.]|_|g')
    
    # Run root tests; allow failure
    go test -buildvcs=false -exec sudo -coverpkg="$COVER_PKGS" -coverprofile="$COVERAGE_DIR/coverage-root-$pkg_slug.out" "$pkg" 2>&1 | grep -v "warning: no packages being tested depend on matches" | sed 's/ of statements in .*/ of statements/' || true
done

echo "3. Merging coverage profiles..."
# If root tests ran, files are owned by root. Fix permissions so current user can read them.
if ls "$COVERAGE_DIR"/coverage-root-*.out >/dev/null 2>&1; then
    sudo chown "$(id -u):$(id -g)" "$COVERAGE_DIR"/coverage-root-*.out
fi

echo "mode: set" > "$COVERAGE_DIR/coverage.out"
# Merge profiles using awk
awk '$1 !~ /mode:/ {k=$1 " " $2; c[k]+=$3} END {for (k in c) print k, c[k]}' "$COVERAGE_DIR"/coverage-*.out | sort >> "$COVERAGE_DIR/coverage.out"

echo "4. Generating HTML report..."
go tool cover -html="$COVERAGE_DIR/coverage.out" -o "$COVERAGE_DIR/coverage.html"

echo "5. Total Coverage Summary:"
go tool cover -func="$COVERAGE_DIR/coverage.out" | grep total:

echo "   Core Logic (excluding cmd/):"
grep -v "/cmd/" "$COVERAGE_DIR/coverage.out" > "$COVERAGE_DIR/coverage-core.out"
go tool cover -func="$COVERAGE_DIR/coverage-core.out" | grep total:

echo "6. Package Level Breakdown:"
awk '
    /^mode:/ { next }
    {
        # Extract file path (remove range suffix)
        file = $1
        sub(/:.+/, "", file)
        
        # Extract package path (remove filename)
        pkg = file
        sub(/\/[^/]+$/, "", pkg)
        
        # $2 = num statements, $3 = count
        total[pkg] += $2
        if ($3 > 0) {
            covered[pkg] += $2
        }
    }
    END {
        for (p in total) {
            pct = 0
            if (total[p] > 0) pct = (covered[p] / total[p]) * 100
            printf "%-50s %5.1f%%\n", p, pct
        }
    }
' "$COVERAGE_DIR/coverage.out" | sort

echo "Done. Report available at $COVERAGE_DIR/coverage.html"