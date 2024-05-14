grammar ContactQL;

import LexUnicode;

// Lexer rules
fragment HAS: [Hh][Aa][Ss];
fragment IS: [Ii][Ss];

LPAREN: '(';
RPAREN: ')';
AND: [Aa][Nn][Dd];
OR: [Oo][Rr];
COMPARATOR: (
		'='
		| '!='
		| '~'
		| '>='
		| '<='
		| '>'
		| '<'
		| HAS
		| IS
	);
STRING: '"' (~["] | '\\"')* '"';
NAME: (UnicodeLetter | UnicodeDigit | '_' | '.')+;  // e.g. fields.num_goats or urns.tel or name
TEXT: (
		UnicodeLetter
		| UnicodeDigit
		| '_'
		| '.'
		| '-'
		| '+'
		| '/'
		| '\''
		| '@'
		| ':'
	)+;

WS: [ \t\n\r]+ -> skip; // ignore whitespace

ERROR: .;

// Parser rules
parse: expression EOF;

expression:
	expression AND expression	# combinationAnd
	| expression expression		# combinationImpicitAnd
	| expression OR expression	# combinationOr
	| LPAREN expression RPAREN	# expressionGrouping
	| NAME COMPARATOR literal	# condition
	| literal					# implicitCondition;

literal: 
	NAME # textLiteral 
	| TEXT # textLiteral 
	| STRING # stringLiteral;