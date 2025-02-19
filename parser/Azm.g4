grammar Azm;

relation
    :   rel ('|' rel)*  EOF
    ;

permission
    :   union EOF           #UnionPerm
    |   intersection EOF    #IntersectionPerm
    |   exclusion EOF       #ExclusionPerm
    ;

union
    :   perm ('|' perm)*
    ;

intersection
    :   perm '&' perm ('&' perm)*
    ;

exclusion
    :   perm '-' perm
    ;

rel
    :   direct      #DirectRel
    |   wildcard    #WildcardRel
    |   subject     #SubjectRel
    ;

perm
    :   direct      #DirectPerm
    |   arrow       #ArrowPerm
    ;

direct
    :   ID
    ;

subject
    :   ID HASH ID
    ;

wildcard
    :   ID COLON ASTERISK
    ;

arrow
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

ID: [a-zA-Z][a-zA-Z0-9._-]*[a-zA-Z0-9] ;

WS: [ \t\n\r\f]+ -> skip ;
