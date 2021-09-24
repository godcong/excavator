#!/bin/bash

curl https://www.unicode.org/Public/14.0.0/ucd/Unihan.zip -o Unihan.zip

unzip Unihan.zip

grep kIRG_GSource Unihan_IRGSources.txt|grep -v -E '#'|grep -E 'G0|G1|G7|G8|GE|GK|GRM|GXC|GXH|GZH|GCH|GHZR|GFC|GOCD|GXHZ|GJZ|GGFZ'|awk '{print $1}'|sort|uniq > GB.txt

rm Unihan*.*

