// Code generated from Azm.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Azm
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type AzmParser struct {
	*antlr.BaseParser
}

var AzmParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func azmParserInit() {
	staticData := &AzmParserStaticData
	staticData.LiteralNames = []string{
		"", "'|'", "'&'", "'-'", "", "'#'", "':'", "'*'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "ARROW", "HASH", "COLON", "ASTERISK", "ID", "WS",
	}
	staticData.RuleNames = []string{
		"relation", "permission", "union", "intersection", "exclusion", "rel",
		"singleRel", "subjectRel", "wildcardRel", "arrowRel",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 9, 75, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 4,
		2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 1, 0, 1, 0,
		1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 33, 8,
		1, 1, 2, 1, 2, 1, 2, 5, 2, 38, 8, 2, 10, 2, 12, 2, 41, 9, 2, 1, 3, 1, 3,
		1, 3, 5, 3, 46, 8, 3, 10, 3, 12, 3, 49, 9, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1,
		5, 1, 5, 1, 5, 1, 5, 3, 5, 59, 8, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7,
		1, 8, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 0, 0, 10, 0, 2, 4,
		6, 8, 10, 12, 14, 16, 18, 0, 0, 71, 0, 20, 1, 0, 0, 0, 2, 32, 1, 0, 0,
		0, 4, 34, 1, 0, 0, 0, 6, 42, 1, 0, 0, 0, 8, 50, 1, 0, 0, 0, 10, 58, 1,
		0, 0, 0, 12, 60, 1, 0, 0, 0, 14, 62, 1, 0, 0, 0, 16, 66, 1, 0, 0, 0, 18,
		70, 1, 0, 0, 0, 20, 21, 3, 4, 2, 0, 21, 22, 5, 0, 0, 1, 22, 1, 1, 0, 0,
		0, 23, 24, 3, 4, 2, 0, 24, 25, 5, 0, 0, 1, 25, 33, 1, 0, 0, 0, 26, 27,
		3, 6, 3, 0, 27, 28, 5, 0, 0, 1, 28, 33, 1, 0, 0, 0, 29, 30, 3, 8, 4, 0,
		30, 31, 5, 0, 0, 1, 31, 33, 1, 0, 0, 0, 32, 23, 1, 0, 0, 0, 32, 26, 1,
		0, 0, 0, 32, 29, 1, 0, 0, 0, 33, 3, 1, 0, 0, 0, 34, 39, 3, 10, 5, 0, 35,
		36, 5, 1, 0, 0, 36, 38, 3, 10, 5, 0, 37, 35, 1, 0, 0, 0, 38, 41, 1, 0,
		0, 0, 39, 37, 1, 0, 0, 0, 39, 40, 1, 0, 0, 0, 40, 5, 1, 0, 0, 0, 41, 39,
		1, 0, 0, 0, 42, 47, 3, 10, 5, 0, 43, 44, 5, 2, 0, 0, 44, 46, 3, 10, 5,
		0, 45, 43, 1, 0, 0, 0, 46, 49, 1, 0, 0, 0, 47, 45, 1, 0, 0, 0, 47, 48,
		1, 0, 0, 0, 48, 7, 1, 0, 0, 0, 49, 47, 1, 0, 0, 0, 50, 51, 3, 10, 5, 0,
		51, 52, 5, 3, 0, 0, 52, 53, 3, 10, 5, 0, 53, 9, 1, 0, 0, 0, 54, 59, 3,
		12, 6, 0, 55, 59, 3, 16, 8, 0, 56, 59, 3, 14, 7, 0, 57, 59, 3, 18, 9, 0,
		58, 54, 1, 0, 0, 0, 58, 55, 1, 0, 0, 0, 58, 56, 1, 0, 0, 0, 58, 57, 1,
		0, 0, 0, 59, 11, 1, 0, 0, 0, 60, 61, 5, 8, 0, 0, 61, 13, 1, 0, 0, 0, 62,
		63, 5, 8, 0, 0, 63, 64, 5, 5, 0, 0, 64, 65, 5, 8, 0, 0, 65, 15, 1, 0, 0,
		0, 66, 67, 5, 8, 0, 0, 67, 68, 5, 6, 0, 0, 68, 69, 5, 7, 0, 0, 69, 17,
		1, 0, 0, 0, 70, 71, 5, 8, 0, 0, 71, 72, 5, 4, 0, 0, 72, 73, 5, 8, 0, 0,
		73, 19, 1, 0, 0, 0, 4, 32, 39, 47, 58,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// AzmParserInit initializes any static state used to implement AzmParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewAzmParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func AzmParserInit() {
	staticData := &AzmParserStaticData
	staticData.once.Do(azmParserInit)
}

// NewAzmParser produces a new parser instance for the optional input antlr.TokenStream.
func NewAzmParser(input antlr.TokenStream) *AzmParser {
	AzmParserInit()
	this := new(AzmParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &AzmParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "Azm.g4"

	return this
}

// AzmParser tokens.
const (
	AzmParserEOF      = antlr.TokenEOF
	AzmParserT__0     = 1
	AzmParserT__1     = 2
	AzmParserT__2     = 3
	AzmParserARROW    = 4
	AzmParserHASH     = 5
	AzmParserCOLON    = 6
	AzmParserASTERISK = 7
	AzmParserID       = 8
	AzmParserWS       = 9
)

// AzmParser rules.
const (
	AzmParserRULE_relation     = 0
	AzmParserRULE_permission   = 1
	AzmParserRULE_union        = 2
	AzmParserRULE_intersection = 3
	AzmParserRULE_exclusion    = 4
	AzmParserRULE_rel          = 5
	AzmParserRULE_singleRel    = 6
	AzmParserRULE_subjectRel   = 7
	AzmParserRULE_wildcardRel  = 8
	AzmParserRULE_arrowRel     = 9
)

// IRelationContext is an interface to support dynamic dispatch.
type IRelationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Union() IUnionContext
	EOF() antlr.TerminalNode

	// IsRelationContext differentiates from other interfaces.
	IsRelationContext()
}

type RelationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRelationContext() *RelationContext {
	var p = new(RelationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_relation
	return p
}

func InitEmptyRelationContext(p *RelationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_relation
}

func (*RelationContext) IsRelationContext() {}

func NewRelationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationContext {
	var p = new(RelationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_relation

	return p
}

func (s *RelationContext) GetParser() antlr.Parser { return s.parser }

func (s *RelationContext) Union() IUnionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionContext)
}

func (s *RelationContext) EOF() antlr.TerminalNode {
	return s.GetToken(AzmParserEOF, 0)
}

func (s *RelationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RelationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterRelation(s)
	}
}

func (s *RelationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitRelation(s)
	}
}

func (p *AzmParser) Relation() (localctx IRelationContext) {
	localctx = NewRelationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, AzmParserRULE_relation)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(20)
		p.Union()
	}
	{
		p.SetState(21)
		p.Match(AzmParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPermissionContext is an interface to support dynamic dispatch.
type IPermissionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Union() IUnionContext
	EOF() antlr.TerminalNode
	Intersection() IIntersectionContext
	Exclusion() IExclusionContext

	// IsPermissionContext differentiates from other interfaces.
	IsPermissionContext()
}

type PermissionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPermissionContext() *PermissionContext {
	var p = new(PermissionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_permission
	return p
}

func InitEmptyPermissionContext(p *PermissionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_permission
}

func (*PermissionContext) IsPermissionContext() {}

func NewPermissionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PermissionContext {
	var p = new(PermissionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_permission

	return p
}

func (s *PermissionContext) GetParser() antlr.Parser { return s.parser }

func (s *PermissionContext) Union() IUnionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionContext)
}

func (s *PermissionContext) EOF() antlr.TerminalNode {
	return s.GetToken(AzmParserEOF, 0)
}

func (s *PermissionContext) Intersection() IIntersectionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntersectionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntersectionContext)
}

func (s *PermissionContext) Exclusion() IExclusionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExclusionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExclusionContext)
}

func (s *PermissionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PermissionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PermissionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterPermission(s)
	}
}

func (s *PermissionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitPermission(s)
	}
}

func (p *AzmParser) Permission() (localctx IPermissionContext) {
	localctx = NewPermissionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, AzmParserRULE_permission)
	p.SetState(32)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(23)
			p.Union()
		}
		{
			p.SetState(24)
			p.Match(AzmParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(26)
			p.Intersection()
		}
		{
			p.SetState(27)
			p.Match(AzmParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(29)
			p.Exclusion()
		}
		{
			p.SetState(30)
			p.Match(AzmParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnionContext is an interface to support dynamic dispatch.
type IUnionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRel() []IRelContext
	Rel(i int) IRelContext

	// IsUnionContext differentiates from other interfaces.
	IsUnionContext()
}

type UnionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnionContext() *UnionContext {
	var p = new(UnionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_union
	return p
}

func InitEmptyUnionContext(p *UnionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_union
}

func (*UnionContext) IsUnionContext() {}

func NewUnionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionContext {
	var p = new(UnionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_union

	return p
}

func (s *UnionContext) GetParser() antlr.Parser { return s.parser }

func (s *UnionContext) AllRel() []IRelContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRelContext); ok {
			len++
		}
	}

	tst := make([]IRelContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRelContext); ok {
			tst[i] = t.(IRelContext)
			i++
		}
	}

	return tst
}

func (s *UnionContext) Rel(i int) IRelContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelContext)
}

func (s *UnionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterUnion(s)
	}
}

func (s *UnionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitUnion(s)
	}
}

func (p *AzmParser) Union() (localctx IUnionContext) {
	localctx = NewUnionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, AzmParserRULE_union)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(34)
		p.Rel()
	}
	p.SetState(39)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == AzmParserT__0 {
		{
			p.SetState(35)
			p.Match(AzmParserT__0)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(36)
			p.Rel()
		}

		p.SetState(41)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIntersectionContext is an interface to support dynamic dispatch.
type IIntersectionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRel() []IRelContext
	Rel(i int) IRelContext

	// IsIntersectionContext differentiates from other interfaces.
	IsIntersectionContext()
}

type IntersectionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntersectionContext() *IntersectionContext {
	var p = new(IntersectionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_intersection
	return p
}

func InitEmptyIntersectionContext(p *IntersectionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_intersection
}

func (*IntersectionContext) IsIntersectionContext() {}

func NewIntersectionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntersectionContext {
	var p = new(IntersectionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_intersection

	return p
}

func (s *IntersectionContext) GetParser() antlr.Parser { return s.parser }

func (s *IntersectionContext) AllRel() []IRelContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRelContext); ok {
			len++
		}
	}

	tst := make([]IRelContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRelContext); ok {
			tst[i] = t.(IRelContext)
			i++
		}
	}

	return tst
}

func (s *IntersectionContext) Rel(i int) IRelContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelContext)
}

func (s *IntersectionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntersectionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntersectionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterIntersection(s)
	}
}

func (s *IntersectionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitIntersection(s)
	}
}

func (p *AzmParser) Intersection() (localctx IIntersectionContext) {
	localctx = NewIntersectionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, AzmParserRULE_intersection)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(42)
		p.Rel()
	}
	p.SetState(47)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == AzmParserT__1 {
		{
			p.SetState(43)
			p.Match(AzmParserT__1)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(44)
			p.Rel()
		}

		p.SetState(49)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExclusionContext is an interface to support dynamic dispatch.
type IExclusionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllRel() []IRelContext
	Rel(i int) IRelContext

	// IsExclusionContext differentiates from other interfaces.
	IsExclusionContext()
}

type ExclusionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExclusionContext() *ExclusionContext {
	var p = new(ExclusionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_exclusion
	return p
}

func InitEmptyExclusionContext(p *ExclusionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_exclusion
}

func (*ExclusionContext) IsExclusionContext() {}

func NewExclusionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExclusionContext {
	var p = new(ExclusionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_exclusion

	return p
}

func (s *ExclusionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExclusionContext) AllRel() []IRelContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRelContext); ok {
			len++
		}
	}

	tst := make([]IRelContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRelContext); ok {
			tst[i] = t.(IRelContext)
			i++
		}
	}

	return tst
}

func (s *ExclusionContext) Rel(i int) IRelContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelContext)
}

func (s *ExclusionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExclusionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExclusionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterExclusion(s)
	}
}

func (s *ExclusionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitExclusion(s)
	}
}

func (p *AzmParser) Exclusion() (localctx IExclusionContext) {
	localctx = NewExclusionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, AzmParserRULE_exclusion)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(50)
		p.Rel()
	}
	{
		p.SetState(51)
		p.Match(AzmParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(52)
		p.Rel()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRelContext is an interface to support dynamic dispatch.
type IRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SingleRel() ISingleRelContext
	WildcardRel() IWildcardRelContext
	SubjectRel() ISubjectRelContext
	ArrowRel() IArrowRelContext

	// IsRelContext differentiates from other interfaces.
	IsRelContext()
}

type RelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRelContext() *RelContext {
	var p = new(RelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_rel
	return p
}

func InitEmptyRelContext(p *RelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_rel
}

func (*RelContext) IsRelContext() {}

func NewRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelContext {
	var p = new(RelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_rel

	return p
}

func (s *RelContext) GetParser() antlr.Parser { return s.parser }

func (s *RelContext) SingleRel() ISingleRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISingleRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISingleRelContext)
}

func (s *RelContext) WildcardRel() IWildcardRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWildcardRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWildcardRelContext)
}

func (s *RelContext) SubjectRel() ISubjectRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISubjectRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISubjectRelContext)
}

func (s *RelContext) ArrowRel() IArrowRelContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrowRelContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrowRelContext)
}

func (s *RelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterRel(s)
	}
}

func (s *RelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitRel(s)
	}
}

func (p *AzmParser) Rel() (localctx IRelContext) {
	localctx = NewRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, AzmParserRULE_rel)
	p.SetState(58)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(54)
			p.SingleRel()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(55)
			p.WildcardRel()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(56)
			p.SubjectRel()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(57)
			p.ArrowRel()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISingleRelContext is an interface to support dynamic dispatch.
type ISingleRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode

	// IsSingleRelContext differentiates from other interfaces.
	IsSingleRelContext()
}

type SingleRelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySingleRelContext() *SingleRelContext {
	var p = new(SingleRelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_singleRel
	return p
}

func InitEmptySingleRelContext(p *SingleRelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_singleRel
}

func (*SingleRelContext) IsSingleRelContext() {}

func NewSingleRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SingleRelContext {
	var p = new(SingleRelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_singleRel

	return p
}

func (s *SingleRelContext) GetParser() antlr.Parser { return s.parser }

func (s *SingleRelContext) ID() antlr.TerminalNode {
	return s.GetToken(AzmParserID, 0)
}

func (s *SingleRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SingleRelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SingleRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterSingleRel(s)
	}
}

func (s *SingleRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitSingleRel(s)
	}
}

func (p *AzmParser) SingleRel() (localctx ISingleRelContext) {
	localctx = NewSingleRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, AzmParserRULE_singleRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(60)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISubjectRelContext is an interface to support dynamic dispatch.
type ISubjectRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	HASH() antlr.TerminalNode

	// IsSubjectRelContext differentiates from other interfaces.
	IsSubjectRelContext()
}

type SubjectRelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySubjectRelContext() *SubjectRelContext {
	var p = new(SubjectRelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_subjectRel
	return p
}

func InitEmptySubjectRelContext(p *SubjectRelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_subjectRel
}

func (*SubjectRelContext) IsSubjectRelContext() {}

func NewSubjectRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubjectRelContext {
	var p = new(SubjectRelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_subjectRel

	return p
}

func (s *SubjectRelContext) GetParser() antlr.Parser { return s.parser }

func (s *SubjectRelContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(AzmParserID)
}

func (s *SubjectRelContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(AzmParserID, i)
}

func (s *SubjectRelContext) HASH() antlr.TerminalNode {
	return s.GetToken(AzmParserHASH, 0)
}

func (s *SubjectRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubjectRelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SubjectRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterSubjectRel(s)
	}
}

func (s *SubjectRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitSubjectRel(s)
	}
}

func (p *AzmParser) SubjectRel() (localctx ISubjectRelContext) {
	localctx = NewSubjectRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, AzmParserRULE_subjectRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(62)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(63)
		p.Match(AzmParserHASH)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(64)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IWildcardRelContext is an interface to support dynamic dispatch.
type IWildcardRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	COLON() antlr.TerminalNode
	ASTERISK() antlr.TerminalNode

	// IsWildcardRelContext differentiates from other interfaces.
	IsWildcardRelContext()
}

type WildcardRelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWildcardRelContext() *WildcardRelContext {
	var p = new(WildcardRelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_wildcardRel
	return p
}

func InitEmptyWildcardRelContext(p *WildcardRelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_wildcardRel
}

func (*WildcardRelContext) IsWildcardRelContext() {}

func NewWildcardRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WildcardRelContext {
	var p = new(WildcardRelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_wildcardRel

	return p
}

func (s *WildcardRelContext) GetParser() antlr.Parser { return s.parser }

func (s *WildcardRelContext) ID() antlr.TerminalNode {
	return s.GetToken(AzmParserID, 0)
}

func (s *WildcardRelContext) COLON() antlr.TerminalNode {
	return s.GetToken(AzmParserCOLON, 0)
}

func (s *WildcardRelContext) ASTERISK() antlr.TerminalNode {
	return s.GetToken(AzmParserASTERISK, 0)
}

func (s *WildcardRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WildcardRelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WildcardRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterWildcardRel(s)
	}
}

func (s *WildcardRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitWildcardRel(s)
	}
}

func (p *AzmParser) WildcardRel() (localctx IWildcardRelContext) {
	localctx = NewWildcardRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, AzmParserRULE_wildcardRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(66)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(67)
		p.Match(AzmParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(68)
		p.Match(AzmParserASTERISK)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArrowRelContext is an interface to support dynamic dispatch.
type IArrowRelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllID() []antlr.TerminalNode
	ID(i int) antlr.TerminalNode
	ARROW() antlr.TerminalNode

	// IsArrowRelContext differentiates from other interfaces.
	IsArrowRelContext()
}

type ArrowRelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrowRelContext() *ArrowRelContext {
	var p = new(ArrowRelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_arrowRel
	return p
}

func InitEmptyArrowRelContext(p *ArrowRelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = AzmParserRULE_arrowRel
}

func (*ArrowRelContext) IsArrowRelContext() {}

func NewArrowRelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrowRelContext {
	var p = new(ArrowRelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = AzmParserRULE_arrowRel

	return p
}

func (s *ArrowRelContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrowRelContext) AllID() []antlr.TerminalNode {
	return s.GetTokens(AzmParserID)
}

func (s *ArrowRelContext) ID(i int) antlr.TerminalNode {
	return s.GetToken(AzmParserID, i)
}

func (s *ArrowRelContext) ARROW() antlr.TerminalNode {
	return s.GetToken(AzmParserARROW, 0)
}

func (s *ArrowRelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrowRelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArrowRelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.EnterArrowRel(s)
	}
}

func (s *ArrowRelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(AzmListener); ok {
		listenerT.ExitArrowRel(s)
	}
}

func (p *AzmParser) ArrowRel() (localctx IArrowRelContext) {
	localctx = NewArrowRelContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, AzmParserRULE_arrowRel)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(70)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(71)
		p.Match(AzmParserARROW)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(72)
		p.Match(AzmParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
