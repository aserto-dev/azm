package parser

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aserto-dev/azm/model"
	"github.com/pkg/errors"
)

func ParseRelation(input string) ([]*model.RelationRef, error) {
	p := newParser(input)
	rTree, err := p.Relation(), p.Error()
	if err != nil {
		return nil, err
	}

	var v RelationVisitor
	return v.Visit(rTree).([]*model.RelationRef), nil
}

func ParsePermission(input string) (*model.Permission, error) {
	p := newParser(input)
	pTree, err := p.Permission(), p.Error()
	if err != nil {
		return nil, err
	}

	var v PermissionVisitor
	return v.Visit(pTree).(*model.Permission), nil
}

func newParser(input string) *parser {
	lexer := NewAzmLexer(antlr.NewInputStream(input))
	stream := antlr.NewCommonTokenStream(lexer, 0)
	listener := newErrorListener(input)

	p := NewAzmParser(stream)
	p.AddErrorListener(listener)
	if os.Getenv("AZM_DIAGNOSTICS") == "1" {
		p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	}
	return &parser{AzmParser: p, listener: listener}
}

type parser struct {
	*AzmParser
	listener *errorListener
}

func (p *parser) Error() error {
	return p.listener.err
}

type errorListener struct {
	*antlr.DefaultErrorListener

	input string
	err   error
}

func newErrorListener(input string) *errorListener {
	return &errorListener{antlr.NewDefaultErrorListener(), input, nil}
}

func (l *errorListener) SyntaxError(
	_ antlr.Recognizer, _ any,
	line, column int,
	msg string,
	e antlr.RecognitionException,
) {
	l.err = errors.Wrap(model.ErrInvalidIdentifier, fmt.Sprintf("%s in '%s'", msg, l.input))
}
