#!/usr/bin/env bash

# Copyright (c) Juniper Networks, Inc., 2024-2026.
# All rights reserved.
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

# files we're not interested in
skip_regexes=()
skip_regexes+=("^LICENSE$")
skip_regexes+=("^NOTICE$")
skip_regexes+=("^README.md$")
skip_regexes+=("^Third_Party_Code/.*$")
skip_regexes+=("^\.gitignore$")
skip_regexes+=("^\.notices.tpl$")
skip_regexes+=("^go.mod$")
skip_regexes+=("^go.sum$")
skip_regexes+=("^kb/.*md$")

copyright_template="Copyright \\(c\\) Juniper Networks, Inc\\., %s-%s\\."
arr_line="All rights reserved\."
spdx_line="SPDX-License-Identifier: Apache-2.0"
leading_comment="[ #/]{0,3}"

problem_files=()
problem_headers=()

# loop over files changed relative to "main" branch
for file in $(git diff --name-only origin/main)
do
  # skip over non-files and files which don't exist in this branch
  [ ! -f "$file" ] && continue

  # skip over files which don't require the license header
  skip=""
  for re in "${skip_regexes[@]}"
  do
    grep -Eq "$re" <<< "$file" && skip="yes" && break
  done
  # shellcheck disable=SC2059
  [ -n "$skip" ] && printf "skipping %s\n" "$file" && continue

  # shellcheck disable=SC2059
  printf "checking %s...  " "$file"



  # Determine the year the file was introduced and the year it was most recently modified.
  #
  # Note that Git 2.56 changed how `git log --follow` tracks renames across merges. The
  # same repository history can therefore yield a different introduction year than earlier
  # Git versions, even with the default rename similarity threshold unchanged. We accept
  # the lineage inferred by the Git version running this check. Files headers previously
  # validated with an earlier Git release may need their headers updated to reflect the
  # Git's current understanding of the file's history.
  years=$(git log --follow --format='%ad' --date=format:'%Y' -- "$file" | sort -n)
  if [ -n "$years" ]; then
    first_year=$(sed -n '1p' <<< "$years")
    recent_year=$(sed -n '$p' <<< "$years")
  else
    first_year=$(date '+%Y')
    recent_year="$first_year"
  fi

  # shellcheck disable=SC2059
  copyright_line=$(printf "${copyright_template}" "$first_year" "$recent_year")

  # collect the first few lines of the file
  head=$(head "$file")

  # check the file head
  failed=0
  grep -Eq "${leading_comment}${copyright_line}" <<< "$head" || ((failed+=1))
  grep -Eq "${leading_comment}${arr_line}"       <<< "$head" || ((failed+=2))
  grep -Eq "${leading_comment}${spdx_line}"      <<< "$head" || ((failed+=4))

  if [ "$failed" -eq 0 ]
  then
    echo "ok"
  else
    echo "failure reason: $failed"
    problem_files+=("$file")
    problem_headers+=("${copyright_line//\\/}\n${arr_line//\\/}\n${spdx_line//\\/}")
  fi
done

# print a report if necessary
if [ "${#problem_files[@]}" -gt 0 ]
then
  printf "\nThe following files need to have their license header updated:\n\n"

  for i in $(seq 0 $(( ${#problem_files[@]} - 1 )) )
  do
    printf '%s\n' "${problem_files[$i]}"
    printf -- '-%.0s' {1..40}
    # shellcheck disable=SC2059
    printf "\n${problem_headers[$i]}\n\n"
  done

  exit 1
fi
