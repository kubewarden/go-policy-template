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

func easyjson55ad6c3aDecodeGithubComKubewardenK8sObjectsApiCoreV1(in *jlexer.Lexer, out *LoadBalancerStatus) {
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
		case "ingress":
			if in.IsNull() {
				in.Skip()
				out.Ingress = nil
			} else {
				in.Delim('[')
				if out.Ingress == nil {
					if !in.IsDelim(']') {
						out.Ingress = make([]*LoadBalancerIngress, 0, 8)
					} else {
						out.Ingress = []*LoadBalancerIngress{}
					}
				} else {
					out.Ingress = (out.Ingress)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *LoadBalancerIngress
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(LoadBalancerIngress)
						}
						(*v1).UnmarshalEasyJSON(in)
					}
					out.Ingress = append(out.Ingress, v1)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjson55ad6c3aEncodeGithubComKubewardenK8sObjectsApiCoreV1(out *jwriter.Writer, in LoadBalancerStatus) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ingress\":"
		out.RawString(prefix[1:])
		if in.Ingress == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Ingress {
				if v2 > 0 {
					out.RawByte(',')
				}
				if v3 == nil {
					out.RawString("null")
				} else {
					(*v3).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LoadBalancerStatus) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson55ad6c3aEncodeGithubComKubewardenK8sObjectsApiCoreV1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LoadBalancerStatus) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson55ad6c3aEncodeGithubComKubewardenK8sObjectsApiCoreV1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LoadBalancerStatus) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson55ad6c3aDecodeGithubComKubewardenK8sObjectsApiCoreV1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LoadBalancerStatus) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson55ad6c3aDecodeGithubComKubewardenK8sObjectsApiCoreV1(l, v)
}