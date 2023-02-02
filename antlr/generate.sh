antlr -Dlanguage=Go ContactQL.g4 -o ../contactql/gen -package gen -visitor
antlr -Dlanguage=Go Excellent1.g4 -o ../flows/definition/legacy/gen -package gen -visitor
antlr -Dlanguage=Go Excellent3.g4 -o ../excellent/gen -package gen -visitor