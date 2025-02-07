// Code generated by github.com/spacemeshos/go-scale/scalegen. DO NOT EDIT.

package types

import (
	"github.com/spacemeshos/go-scale"
)

func (t *Nonce) EncodeScale(enc *scale.Encoder) (total int, err error) {
	if n, err := scale.EncodeCompact64(enc, t.Counter); err != nil {
		return total, err
	} else {
		total += n
	}
	if n, err := scale.EncodeCompact8(enc, t.Bitfield); err != nil {
		return total, err
	} else {
		total += n
	}
	return total, nil
}

func (t *Nonce) DecodeScale(dec *scale.Decoder) (total int, err error) {
	if field, n, err := scale.DecodeCompact64(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Counter = field
	}
	if field, n, err := scale.DecodeCompact8(dec); err != nil {
		return total, err
	} else {
		total += n
		t.Bitfield = field
	}
	return total, nil
}
