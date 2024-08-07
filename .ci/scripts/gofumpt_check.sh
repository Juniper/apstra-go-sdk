#!/usr/bin/env bash
set -euo pipefail

# init array of files which need updating by gofumpt
needs_update=()

echo "git branch"
git branch
echo "git branch"

# loop over files changed relative to "main" branch
for file in $(git diff --name-only main)
do
  OUT="" # init variable because of `set -u`

  # skip over non-Go files
  [[ $file = *.go ]] && OUT=$(go run mvdan.cc/gofumpt -l "$file")

  if [ -n "$OUT" ]
  then
    # save file which was recorded to OUT
    needs_update+=("$OUT")
  fi
done

if [ ${#needs_update[@]} -gt 0 ]
then
  echo "Formatting required:"
  echo ""

  for f in "${needs_update[@]}"
  do
     echo "  go run mvdan.cc/gofumpt -w '$file'"
  done

  echo ""
  exit 1
fi
