/* Copyright 2019 The Bazel Authors. All rights reserved.

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

type LeftParenthesis struct{}

func (s *LeftParenthesis) String() string { return "<LeftParenthesis>" }
func (_ *LeftParenthesis) Type() string   { return "LeftParenthesis" }

var _ Token = &LeftParenthesis{}

type RightParenthesis struct{}

func (s *RightParenthesis) String() string { return "<RightParenthesis>" }
func (_ *RightParenthesis) Type() string   { return "RightParenthesis" }

var _ Token = &RightParenthesis{}

type PlusSign struct{}

func (s *PlusSign) String() string { return "<PlusSign>" }
func (_ *PlusSign) Type() string   { return "PlusSign" }

var _ Token = &PlusSign{}

type Comma struct{}

func (s *Comma) String() string { return "<Comma>" }
func (_ *Comma) Type() string   { return "Comma" }

var _ Token = &Comma{}

type GreaterThanSign struct{}

func (s *GreaterThanSign) String() string { return "<GreaterThanSign>" }
func (_ *GreaterThanSign) Type() string   { return "GreaterThanSign" }

var _ Token = &GreaterThanSign{}

type CDC struct{}

func (s *CDC) String() string { return "<CDC>" }
func (_ *CDC) Type() string   { return "CDC" }

var _ Token = &CDC{}

type Colon struct{}

func (s *Colon) String() string { return "<Colon>" }
func (_ *Colon) Type() string   { return "Colon" }

var _ Token = &Colon{}

type Semicolon struct{}

func (s *Semicolon) String() string { return "<Semicolon>" }
func (_ *Semicolon) Type() string   { return "Semicolon" }

var _ Token = &Semicolon{}

type LessThanSign struct{}

func (s *LessThanSign) String() string { return "<LessThanSign>" }
func (_ *LessThanSign) Type() string   { return "LessThanSign" }

var _ Token = &LessThanSign{}

type CDO struct{}

func (s *CDO) String() string { return "<CDO>" }
func (_ *CDO) Type() string   { return "CDO" }

var _ Token = &CDO{}

type LeftSquareBracket struct{}

func (s *LeftSquareBracket) String() string { return "<LeftSquareBracket>" }
func (_ *LeftSquareBracket) Type() string   { return "LeftSquareBracket" }

var _ Token = &LeftSquareBracket{}

type RightSquareBracket struct{}

func (s *RightSquareBracket) String() string { return "<RightSquareBracket>" }
func (_ *RightSquareBracket) Type() string   { return "RightSquareBracket" }

var _ Token = &RightSquareBracket{}

type LeftCurlyBracket struct{}

func (s *LeftCurlyBracket) String() string { return "<LeftCurlyBracket>" }
func (_ *LeftCurlyBracket) Type() string   { return "LeftCurlyBracket" }

var _ Token = &LeftCurlyBracket{}

type RightCurlyBracket struct{}

func (s *RightCurlyBracket) String() string { return "<RightCurlyBracket>" }
func (_ *RightCurlyBracket) Type() string   { return "RightCurlyBracket" }

var _ Token = &RightCurlyBracket{}
