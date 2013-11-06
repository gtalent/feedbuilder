/*
   Copyright 2013 gtalent2@gmail.com

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package parser

import (
	"../../lex"
	"fmt"
)

type Feed struct {
	Name       string
	FamilySize int
}

/*
  Parses a set of feed specs into a list of feed models.
*/
func Parse(input string) ([]Feed, error) {
	symbols := []string{"[", "]"}
	keywords := []string{}
	stringTypes := []lex.Pair{}
	commentTypes := []lex.Pair{{"#", "\n"}}
	l := lex.NewAnalyzer(symbols, keywords, stringTypes, commentTypes, true)

	tokens := l.TokenList(input)

	line := 1
	var feeds []Feed
	var err error
	for t := tokens.Next(); tokens.HasNext(); t = tokens.Next() {
		switch t.Type {
		case lex.Comment:
			// ignore comment tokens
		case lex.Whitespace:
			if t.String() == "\n" {
				line++
			}
		case lex.Identifier:
			var feed Feed
			feed.Name = t.String()
			feed.FamilySize = 1

			//parse size
			if tokens.Peak().Type != lex.Whitespace {
				if tokens.Peak().String() == "[" && tokens.PeakTo(1).Type == lex.IntLiteral && tokens.PeakTo(2).String() == "]" {
					tokens.Next()
					feed.FamilySize = tokens.Next().Int()
					tokens.Next()
				} else {
					return feeds, fmt.Errorf("Error on line %d:\n\tError: unexepected token: \"%s\" (token type: %s)", line, t.String(), lex.TokenType(t.Type))
				}
			}

			feeds = append(feeds, feed)
		default:
			err = fmt.Errorf("Error on line %d:\n\tError: unexpected token: \"%s\" (token type: %s)", line+1, t.String(), lex.TokenType(t.Type))
			return feeds, err
		}
	}
	return feeds, err
}
