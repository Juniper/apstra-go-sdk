#!/usr/bin/env bash

# Copyright (c) Juniper Networks, Inc., 2024-2025.
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

  # determine the year the file was introduced and the year it was most recently modified
  first_year=""
  recent_year=""
  while read -r line
  do
    if [ -z "$recent_year" ]
    then
      recent_year="$line"
    else
      first_year="$line"
    fi
  done <<< "$( (git log --follow --pretty=format:"%ad" --date=format:'%Y' "$file"; echo) | sed -n '1p;$p' )"

  # assume current year for both values if git log didn't find anything (new file)
  [ -z "$first_year" ] && first_year=$(date '+%Y')
  [ -z "$recent_year" ] && recent_year=$(date '+%Y')

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
