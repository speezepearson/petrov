#!/bin/bash

set -e
set -x

cat World_map_blank_without_borders.svg | sed 's/#bcbcbc/#0088dd/g' > test.svg
inkscape     --export-png=export.png --export-dpi=100     --export-background-opacity=0 --without-gui test.svg
convert export.png -background '#000022' -alpha remove big.png
convert big.png -resize 200x200 smol.png
convert smol.png -scale 1000% map.png
optipng map.png
