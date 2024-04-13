/*
Package metadata provides procedures to detect, extract and parse metadata
headers from arbitrary byte buffers.

Example document:

	date = "2024-03-02"
	site = "example.com"
	total = 76373

	The empty line above ends the metadata block.
	These last two lines are not metadata, they are the document's content.

Unlike other front-matter processors, the procedures provided in this package
do not use delimiters/fences to separate metadata from content.

The metadata block starts on the first byte of the buffer, and extends up to
the first occurrence of two consecutive newline characters (ignoring carriage
return characters for compatibility).

If the first byte of the buffer is not a valid character for a key in the
language that the metadata block is written, the entire buffer is treated as
normal content and the metadata block is presumed to not exist/be empty.

The metadata block may pass the above detection heuristic but fail to parse
correctly. The programmer must then decide whether to treat this error as an
error, or to ignore the error and treat the metadata block as absent (in other
words, to treat the block as document content).

The idiomatic way to prevent attempts to parse first line paragraphs as
metadata is to start the buffer with an empty line or white space character.

Empty, non-nil maps are returned when there is no metadata block to prevent
nil dereference errors.

Limitations: at present, only TOML is recognized as metadata.
*/
package metadata

import "cdop.pt/go/free/platepipe/metadata/toml"

// IsPresent heuristically checks if metadata in present in the buffer.
//
// If the buffer seems to have metadata, IsPresent will return true and the
// position of the first content byte in the buffer. This is useful for slicing
// the buffer for further processing, without imposing any memory allocation
// penalties.
func IsPresent(buf []byte) (bool, int) {
	bufSz := len(buf)

	if bufSz == 0 {
		return false, 0
	}

	if !isStartOfKey(buf[0]) {
		return false, 0
	}

	last := buf[0]
	for i := 1; i < bufSz; i++ {
		if buf[i] == '\r' {
			continue
		}

		if buf[i] == '\n' && last == '\n' {
			return true, i + 1
		}

		last = buf[i]
	}

	return false, 0
}

func isStartOfKey(b byte) bool {
	return b == '\'' || b == '"' ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9')
}

// FromTomlBuffer converts a TOML encoded buffer to a key/value map.
//
// FromTomlBuffer is meant to receive a slice of the original buffer, as in
// the following example:
//
//	present, pos := IsPresent(buf)
//	if present {
//	    data, _ := FromTomlBuffer(buf[:pos])
//	    ProcessContent(buf[pos:])
//	}
func FromTomlBuffer(buf []byte) (map[string]any, error) {
	ret := map[string]any{}

	err := toml.Parse(buf, &ret)
	if err != nil {
		return map[string]any{}, err
	}

	return ret, nil
}
