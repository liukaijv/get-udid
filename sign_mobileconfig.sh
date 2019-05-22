#!/usr/bin/env bash

#/usr/bin/security find-identity -p codesigning -v

/usr/bin/security cms -S -N "your identity" -i ./udid_unsigned.mobileconfig -o ./udid.mobileconfig