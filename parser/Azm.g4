grammar Azm;

prog:   stat+ ;

stat:   unionRel NEWLINE
    |   intersectRel NEWLINE
    |   exclusionRel NEWLINE
    |   NEWLINE
    ;

unionRel
    :   rel ('|' rel)*
    ;

intersectRel
    :   rel ('&' rel)*
    ;

exclusionRel
    :   rel '-' rel
    ;

rel
    :   singleRel
    |   wildcardRel
    |   subjectRel
    |   arrowRel
    ;

singleRel
    :   ID
    ;
    
subjectRel
    :   ID HASH ID
    ;

wildcardRel
    :   ID COLON ASTERISK
    ;

arrowRel
    :   ID ARROW ID
    ;
    
ARROW:
    '-''>' ;

HASH:   
    '#' ;

COLON:  
    ':' ;

ASTERISK:
    '*';

ID: [a-z][a-z0-9._-]*[a-z0-9] ;

NEWLINE:
    '\r'? '\n' ;

WS: [ \t\n\r\f]+ -> skip ;
