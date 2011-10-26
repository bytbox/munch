#!/bin/sh

echo "package main"
echo -n "var $1 = \`"
cat
echo "\`"

