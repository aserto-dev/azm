// Code generated from Azm.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type AzmLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var AzmLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func azmlexerLexerInit() {
	staticData := &AzmLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "'|'", "'&'", "'-'", "", "'#'", "':'", "'*'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "ARROW", "HASH", "COLON", "ASTERISK", "ID", "WS",
	}
	staticData.RuleNames = []string{
		"T__0", "T__1", "T__2", "ARROW", "HASH", "COLON", "ASTERISK", "ID",
		"WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 9, 50, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 1, 0, 1, 0, 1,
		1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 6, 1,
		6, 1, 7, 1, 7, 5, 7, 37, 8, 7, 10, 7, 12, 7, 40, 9, 7, 1, 7, 1, 7, 1, 8,
		4, 8, 45, 8, 8, 11, 8, 12, 8, 46, 1, 8, 1, 8, 0, 0, 9, 1, 1, 3, 2, 5, 3,
		7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 1, 0, 4, 1, 0, 97, 122, 4, 0, 45,
		46, 48, 57, 95, 95, 97, 122, 2, 0, 48, 57, 97, 122, 3, 0, 9, 10, 12, 13,
		32, 32, 51, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7,
		1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0,
		15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 1, 19, 1, 0, 0, 0, 3, 21, 1, 0, 0, 0,
		5, 23, 1, 0, 0, 0, 7, 25, 1, 0, 0, 0, 9, 28, 1, 0, 0, 0, 11, 30, 1, 0,
		0, 0, 13, 32, 1, 0, 0, 0, 15, 34, 1, 0, 0, 0, 17, 44, 1, 0, 0, 0, 19, 20,
		5, 124, 0, 0, 20, 2, 1, 0, 0, 0, 21, 22, 5, 38, 0, 0, 22, 4, 1, 0, 0, 0,
		23, 24, 5, 45, 0, 0, 24, 6, 1, 0, 0, 0, 25, 26, 5, 45, 0, 0, 26, 27, 5,
		62, 0, 0, 27, 8, 1, 0, 0, 0, 28, 29, 5, 35, 0, 0, 29, 10, 1, 0, 0, 0, 30,
		31, 5, 58, 0, 0, 31, 12, 1, 0, 0, 0, 32, 33, 5, 42, 0, 0, 33, 14, 1, 0,
		0, 0, 34, 38, 7, 0, 0, 0, 35, 37, 7, 1, 0, 0, 36, 35, 1, 0, 0, 0, 37, 40,
		1, 0, 0, 0, 38, 36, 1, 0, 0, 0, 38, 39, 1, 0, 0, 0, 39, 41, 1, 0, 0, 0,
		40, 38, 1, 0, 0, 0, 41, 42, 7, 2, 0, 0, 42, 16, 1, 0, 0, 0, 43, 45, 7,
		3, 0, 0, 44, 43, 1, 0, 0, 0, 45, 46, 1, 0, 0, 0, 46, 44, 1, 0, 0, 0, 46,
		47, 1, 0, 0, 0, 47, 48, 1, 0, 0, 0, 48, 49, 6, 8, 0, 0, 49, 18, 1, 0, 0,
		0, 3, 0, 38, 46, 1, 6, 0, 0,
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

// AzmLexerInit initializes any static state used to implement AzmLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewAzmLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func AzmLexerInit() {
	staticData := &AzmLexerLexerStaticData
	staticData.once.Do(azmlexerLexerInit)
}

// NewAzmLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewAzmLexer(input antlr.CharStream) *AzmLexer {
	AzmLexerInit()
	l := new(AzmLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &AzmLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "Azm.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// AzmLexer tokens.
const (
	AzmLexerT__0     = 1
	AzmLexerT__1     = 2
	AzmLexerT__2     = 3
	AzmLexerARROW    = 4
	AzmLexerHASH     = 5
	AzmLexerCOLON    = 6
	AzmLexerASTERISK = 7
	AzmLexerID       = 8
	AzmLexerWS       = 9
)
