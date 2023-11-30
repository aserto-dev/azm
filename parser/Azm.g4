grammar Azm;

relation
    :   union EOF
    ;
    
permission
    :   union EOF
    |   intersection EOF
    |   exclusion EOF
    ;

union
    :   rel ('|' rel)*
    ;

intersection
    :   rel ('&' rel)*
    ;

exclusion
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

ID: [a-z][a-z0-9._]*[a-z0-9] ;

WS: [ \t\n\r\f]+ -> skip ;
