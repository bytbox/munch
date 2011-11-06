#!/bin/sh

echo "package main"
echo
echo -n "var $1 = \`"
cat
echo "\`"

