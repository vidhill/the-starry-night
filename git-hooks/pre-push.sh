#!/bin/sh

unformatted=$(gofmt -l .)

[ -z "$unformatted" ] && exit 0

echo >&2 "Go files should be formatted with gofmt.\n"
echo >&2 "Unformatted files: \n"
echo >&2 "${unformatted} \n"
echo >&2 "Please run:"
echo >&2 "  \"gofmt -w .\" \n"
echo >&2 " - to automatically format go files" 
exit 1