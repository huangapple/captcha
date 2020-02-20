// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	if c == "" {
		t.Errorf("expected id, got empty string")
	}
}

func TestVerify(t *testing.T) {
	id := New()
	if Verify(id, "00") {
		t.Errorf("verified wrong captcha")
	}
	id = New()
	d := globalStore.Get(id, false) // cheating
	if !Verify(id, d) {
		t.Errorf("proper captcha not verified")
	}
}

func TestReload(t *testing.T) {
	id := New()
	d1 := globalStore.Get(id, false) // cheating
	Reload(id, Expiration, DefaultLeftTimes)
	d2 := globalStore.Get(id, false) // cheating again
	if d1 == d2 {
		t.Errorf("reload didn't work: %v = %v", d1, d2)
	}
}

func TestRandomDigits(t *testing.T) {
	d1 := RandomDigits(10)
	for _, v := range d1 {
		if v > 9 {
			t.Errorf("digits not in range 0-9: %v", d1)
		}
	}
	d2 := RandomDigits(10)
	if bytes.Equal(d1, d2) {
		t.Errorf("digits seem to be not random")
	}
}

func TestLeftTimes(t *testing.T) {

	captcha := RandomDigitsString(6)
	id := NewID()

	SetID(id, captcha, Expiration, DefaultLeftTimes)

	for i := 0; i < DefaultLeftTimes; i++ {
		if Verify(id, captcha+"1111") {
			t.Errorf("TestLeftTimes error")
		}
	}
	if Verify(id, captcha) {
		t.Errorf("TestLeftTimes error")
	}
}

func TestLeftTimes2(t *testing.T) {

	captcha := RandomDigitsString(6)
	id := NewID()

	SetID(id, captcha, Expiration, DefaultLeftTimes)

	for i := 0; i < DefaultLeftTimes-1; i++ {
		if Verify(id, captcha+"1111") {
			t.Errorf("TestLeftTimes error")
		}
	}
	if !Verify(id, captcha) {
		t.Errorf("TestLeftTimes error")
	}
	if Verify(id, captcha) {
		t.Errorf("TestLeftTimes error")
	}
}

func TestLeftTimes3(t *testing.T) {

	captcha := RandomDigitsString(6)
	id := NewID()

	SetID(id, captcha, Expiration, DefaultLeftTimes)

	for i := 0; i < DefaultLeftTimes-2; i++ {
		if Verify(id, captcha+"1111") {
			t.Errorf("TestLeftTimes error")
		}
	}
	if !Verify(id, captcha) {
		t.Errorf("TestLeftTimes error")
	}
	if Verify(id, captcha) {
		t.Errorf("TestLeftTimes error")
	}
}
