#!/usr/bin/python

import sys

# check for filename on cmd line and use that,
# else use stdin for parsing
fn = ""
if len(sys.argv) > 1:
    fn = sys.argv[1]
if fn:
    inf = open(fn)
else:
    inf = sys.stdin

# print header for json data
print "{"
print "\t\"dishes\": ["

# loop each line
for line in inf:
    a = line.split(":")
    name = a[0].strip()
    print "{"
    print "\"name\": " + "\"%s\"," %name
    print "\"ingredients\": " + "[",
    ingredients = a[1].split(",")

    for item in ingredients:
        print "\"" + item.strip() + "\",",
    print "\"\" ]"

    print "},"

# print footer for json data
print "{}\n\t]\n}"
