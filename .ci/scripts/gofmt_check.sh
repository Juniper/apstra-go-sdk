#!/usr/bin/env bash

# Copyright (c) Juniper Networks, Inc., 2024-2024.
# All rights reserved.
# SPDX-License-Identifier: Apache-2.0

echo -n "Check formatting... "
require_formatting=$(gofmt -l .)

if [[ -n "${require_formatting}" ]]; then
  echo "FAILED"
  echo "${require_formatting}"
  exit 1
else
  echo "OK"
  exit 0
fi
