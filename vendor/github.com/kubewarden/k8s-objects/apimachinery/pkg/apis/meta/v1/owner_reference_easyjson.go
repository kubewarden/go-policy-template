// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package v1

import (
	json "encoding/json"
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

func easyjson2644fe6dDecodeGithubComKubewardenK8sObjectsApimachineryPkgApisMetaV1(in *jlexer.Lexer, out *OwnerReference) {
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
		case "apiVersion":
			if in.IsNull() {
				in.Skip()
				out.APIVersion = nil
			} else {
				if out.APIVersion == nil {
					out.APIVersion = new(string)
				}
				*out.APIVersion = string(in.String())
			}
		case "blockOwnerDeletion":
			out.BlockOwnerDeletion = bool(in.Bool())
		case "controller":
			out.Controller = bool(in.Bool())
		case "kind":
			if in.IsNull() {
				in.Skip()
				out.Kind = nil
			} else {
				if out.Kind == nil {
					out.Kind = new(string)
				}
				*out.Kind = string(in.String())
			}
		case "name":
			if in.IsNull() {
				in.Skip()
				out.Name = nil
			} else {
				if out.Name == nil {
					out.Name = new(string)
				}
				*out.Name = string(in.String())
			}
		case "uid":
			if in.IsNull() {
				in.Skip()
				out.UID = nil
			} else {
				if out.UID == nil {
					out.UID = new(string)
				}
				*out.UID = string(in.String())
			}
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
func easyjson2644fe6dEncodeGithubComKubewardenK8sObjectsApimachineryPkgApisMetaV1(out *jwriter.Writer, in OwnerReference) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"apiVersion\":"
		out.RawString(prefix[1:])
		if in.APIVersion == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.APIVersion))
		}
	}
	if in.BlockOwnerDeletion {
		const prefix string = ",\"blockOwnerDeletion\":"
		out.RawString(prefix)
		out.Bool(bool(in.BlockOwnerDeletion))
	}
	if in.Controller {
		const prefix string = ",\"controller\":"
		out.RawString(prefix)
		out.Bool(bool(in.Controller))
	}
	{
		const prefix string = ",\"kind\":"
		out.RawString(prefix)
		if in.Kind == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Kind))
		}
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		if in.Name == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Name))
		}
	}
	{
		const prefix string = ",\"uid\":"
		out.RawString(prefix)
		if in.UID == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.UID))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v OwnerReference) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2644fe6dEncodeGithubComKubewardenK8sObjectsApimachineryPkgApisMetaV1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OwnerReference) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2644fe6dEncodeGithubComKubewardenK8sObjectsApimachineryPkgApisMetaV1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *OwnerReference) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2644fe6dDecodeGithubComKubewardenK8sObjectsApimachineryPkgApisMetaV1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OwnerReference) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2644fe6dDecodeGithubComKubewardenK8sObjectsApimachineryPkgApisMetaV1(l, v)
}
