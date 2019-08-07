// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"sync"
	"time"
)

// An object implementing Store interface can be registered with SetCustomStore
// function to handle storage and retrieval of captcha ids and solutions for
// them, replacing the default memory store.
//
// It is the responsibility of an object to delete expired and used captchas
// when necessary (for example, the default memory store collects them in Set
// method after the certain amount of captchas has been stored.)
type Store interface {
	// Set sets the digits for the captcha id.
	Set(id string, captcha string, afterExpire time.Duration)

	// Get returns stored digits for the captcha id. Clear indicates
	// whether the captcha must be deleted from the store.
	Get(id string, clear bool) string
}

type captchaInfo struct {
	captcha    string
	expireTime time.Time //过期时间
}

// memoryStore is an internal store for captcha ids and their values.
type memoryStore struct {
	captchaMap sync.Map
	ticker     *time.Ticker
}

// NewMemoryStore returns a new standard memory store for captchas with the
// given collection threshold and expiration time (duration). The returned
// store must be registered with SetCustomStore to replace the default one.
func NewMemoryStore(interval time.Duration) Store {
	s := new(memoryStore)

	//默认每5分钟gc一次
	if interval == 0 {
		interval = time.Minute * 5
	}

	s.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case _, ok := <-s.ticker.C:
				if ok {
					s.collect()
				} else {
					return
				}
			}
		}
	}()

	return s
}

func (s *memoryStore) Set(id string, captcha string, afterExpire time.Duration) {

	s.captchaMap.Store(id, &captchaInfo{
		captcha:    captcha,
		expireTime: time.Now().Add(afterExpire),
	})

}

//释放
func (s *memoryStore) Stop() {
	s.ticker.Stop()
}

func (s *memoryStore) Get(id string, clear bool) string {

	info, ok := s.captchaMap.Load(id)

	if ok {
		info := info.(*captchaInfo)

		if clear {
			s.captchaMap.Delete(id)
		}
		if info.expireTime.After(time.Now()) {
			return info.captcha
		}

	}
	return ""
}

func (s *memoryStore) collect() {
	now := time.Now()

	s.captchaMap.Range(func(k, v interface{}) bool {
		info := v.(*captchaInfo)
		if now.After(info.expireTime) {
			s.captchaMap.Delete(k)
		}
		return true
	})
}
