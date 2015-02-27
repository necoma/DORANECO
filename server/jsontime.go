package main

import (
	"time"
	)

// from http://qiita.com/taizo/items/2c3a338f1aeea86ce9e2
type jsonTime struct {
	time.Time
}

// formatを設定
func (j jsonTime) format() string {
  return j.Time.Format("2006-01-02 15:04:05 MST")
}

// MarshalJSON() の実装
func (j jsonTime) MarshalJSON() ([]byte, error) {
  return []byte(`"` + j.format() + `"`), nil
}



