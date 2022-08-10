#!/bin/sh

python get.py | sort | uniq | xargs -n 1 -P 8 wget -q
