#!/usr/bin/env bash
set -euo pipefail

TPC=Third_Party_Code

IGNORE=()
IGNORE+=(--ignore)
IGNORE+=(github.com/Juniper) # don't bother with Juniper licenses

go run github.com/google/go-licenses/v2 save   ${IGNORE[@]} --save_path "${TPC}" --force ./...
go run github.com/google/go-licenses/v2 report ${IGNORE[@]} --template .notices.tpl ./... > "${TPC}/NOTICES.md"

# The `save` command above collects only license and notice files from packages with licenses identified as
# `RestrictionsShareLicense` and collects the entire source tree when the license is identified as
# `RestrictionsShareCode`.
#
# It's true that some licenses require us to "make available" the upstream source code, but I'm not sure
# that doing so as *part of this repository* is appropriate.
# 1. The go package system makes it perfectly clear what we're using and where we got it.
# 2. If somebody wants to really push the issue, we'll find a way to deliver the source independent of this repository.
#
# The line below deletes "saved" files other than those beginning with "LICENSE" and "NOTICE"
find "$TPC" -type f ! -name 'LICENSE*' ! -name 'NOTICE*' -print0 | xargs -0 rm --

# We now likely have some empty directories. Get rid of 'em.
find "$TPC" -depth -type d -empty -exec rmdir -- "{}" \;
