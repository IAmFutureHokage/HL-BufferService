// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package kafka_dto

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

func easyjson5cb20343DecodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto(in *jlexer.Lexer, out *WaterLevelRecords) {
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
		case "waterlevels":
			if in.IsNull() {
				in.Skip()
				out.Waterlevels = nil
			} else {
				in.Delim('[')
				if out.Waterlevels == nil {
					if !in.IsDelim(']') {
						out.Waterlevels = make([]WaterLevel, 0, 1)
					} else {
						out.Waterlevels = []WaterLevel{}
					}
				} else {
					out.Waterlevels = (out.Waterlevels)[:0]
				}
				for !in.IsDelim(']') {
					var v1 WaterLevel
					(v1).UnmarshalEasyJSON(in)
					out.Waterlevels = append(out.Waterlevels, v1)
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
func easyjson5cb20343EncodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto(out *jwriter.Writer, in WaterLevelRecords) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"waterlevels\":"
		out.RawString(prefix[1:])
		if in.Waterlevels == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Waterlevels {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WaterLevelRecords) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5cb20343EncodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WaterLevelRecords) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5cb20343EncodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WaterLevelRecords) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5cb20343DecodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WaterLevelRecords) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5cb20343DecodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto(l, v)
}
func easyjson5cb20343DecodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto1(in *jlexer.Lexer, out *WaterLevel) {
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
		case "post_code":
			out.PostCode = string(in.String())
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Date).UnmarshalJSON(data))
			}
		case "water_level":
			out.WaterLevel = int32(in.Int32())
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
func easyjson5cb20343EncodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto1(out *jwriter.Writer, in WaterLevel) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"post_code\":"
		out.RawString(prefix[1:])
		out.String(string(in.PostCode))
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.Raw((in.Date).MarshalJSON())
	}
	{
		const prefix string = ",\"water_level\":"
		out.RawString(prefix)
		out.Int32(int32(in.WaterLevel))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WaterLevel) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5cb20343EncodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WaterLevel) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5cb20343EncodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WaterLevel) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5cb20343DecodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WaterLevel) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5cb20343DecodeGithubComIAmFutureHokageHLBufferServiceInternalAppServicesKafkaDto1(l, v)
}