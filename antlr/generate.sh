#!/bin/sh

# NOTE that Excellent2.g4 isn't used and is only there for historical records

antlr -Dlanguage=Go ContactQL.g4 -o gen/contactql -package gen -visitor
antlr -Dlanguage=Go Excellent1.g4 -o gen/excellent1 -package gen -visitor
antlr -Dlanguage=Go Excellent3.g4 -o gen/excellent3 -package gen -visitor