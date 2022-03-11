#!/bin/bash

##https://www.unicode.org/reports/tr38/tr38-32.html
curl https://www.unicode.org/Public/15.0.0/ucd/Unihan-15.0.0d2.zip -o Unihan.zip

unzip Unihan.zip

grep kIRG_GSource Unihan_IRGSources.txt|grep -v -E '#'|grep -E 'G0|G1|G7|G8|GE|GK|GRM|GXC|GXH|GZH|GCH|GHZR|GFC|GOCD|GXHZ|GJZ|GGFZ'|awk '{print $1}'|sort|uniq > GB.txt

rm Unihan*.*

