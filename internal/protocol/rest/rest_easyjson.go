// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package rest

import (
	json "encoding/json"
	_vault "github.com/k1nky/gophkeeper/internal/entity/vault"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonDfb66c62DecodeGithubComK1nkyGophkeeperInternalProtocolRest(in *jlexer.Lexer, out *RegisterUserRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "login":
			out.Login = string(in.String())
		case "password":
			out.Password = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonDfb66c62EncodeGithubComK1nkyGophkeeperInternalProtocolRest(out *jwriter.Writer, in RegisterUserRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"login\":"
		out.RawString(prefix[1:])
		out.String(string(in.Login))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RegisterUserRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDfb66c62EncodeGithubComK1nkyGophkeeperInternalProtocolRest(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegisterUserRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDfb66c62EncodeGithubComK1nkyGophkeeperInternalProtocolRest(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegisterUserRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDfb66c62DecodeGithubComK1nkyGophkeeperInternalProtocolRest(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegisterUserRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDfb66c62DecodeGithubComK1nkyGophkeeperInternalProtocolRest(l, v)
}
func easyjsonDfb66c62DecodeGithubComK1nkyGophkeeperInternalProtocolRest1(in *jlexer.Lexer, out *NewSecretRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = _vault.MetaID(in.String())
		case "extra":
			out.Extra = string(in.String())
		case "secret":
			out.Secret = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonDfb66c62EncodeGithubComK1nkyGophkeeperInternalProtocolRest1(out *jwriter.Writer, in NewSecretRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"extra\":"
		out.RawString(prefix)
		out.String(string(in.Extra))
	}
	{
		const prefix string = ",\"secret\":"
		out.RawString(prefix)
		out.String(string(in.Secret))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NewSecretRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDfb66c62EncodeGithubComK1nkyGophkeeperInternalProtocolRest1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NewSecretRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDfb66c62EncodeGithubComK1nkyGophkeeperInternalProtocolRest1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NewSecretRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDfb66c62DecodeGithubComK1nkyGophkeeperInternalProtocolRest1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NewSecretRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDfb66c62DecodeGithubComK1nkyGophkeeperInternalProtocolRest1(l, v)
}
func easyjsonDfb66c62DecodeGithubComK1nkyGophkeeperInternalProtocolRest2(in *jlexer.Lexer, out *LoginUserRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "login":
			out.Login = string(in.String())
		case "password":
			out.Password = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonDfb66c62EncodeGithubComK1nkyGophkeeperInternalProtocolRest2(out *jwriter.Writer, in LoginUserRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"login\":"
		out.RawString(prefix[1:])
		out.String(string(in.Login))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LoginUserRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDfb66c62EncodeGithubComK1nkyGophkeeperInternalProtocolRest2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LoginUserRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonDfb66c62EncodeGithubComK1nkyGophkeeperInternalProtocolRest2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LoginUserRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDfb66c62DecodeGithubComK1nkyGophkeeperInternalProtocolRest2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LoginUserRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDfb66c62DecodeGithubComK1nkyGophkeeperInternalProtocolRest2(l, v)
}
