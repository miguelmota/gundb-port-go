package dup

import (
	"math/rand"
	"time"
)

// Opt ...
type Opt struct {
	Max int64
	Age int64
}

// Dup ...
type Dup struct {
	value map[string]interface{}
	opt   *Opt
}

// NewDup ...
func NewDup() *Dup {
	value := make(map[string]interface{})
	value["s"] = make(map[string]interface{})

	return &Dup{
		value: value,
		opt: &Opt{
			Max: 1000,
			Age: 1000 * 9,
		},
	}
}

// Check ...
func (d *Dup) Check(id string) bool {
	_, ok := d.value["s"].(map[string]interface{})[id]
	if ok {
		return d.Track(id)
	}

	return false
}

// Track ...
func (d *Dup) Track(id string) bool {
	d.value["s"].(map[string]interface{})[id] = time.Now()

	_, ok := d.value["to"]
	if !ok {
		d.value["to"] = time.AfterFunc(time.Duration(d.opt.Age), func() {
			for id, v := range d.value["s"].(map[string]interface{}) {
				t := v.(time.Time)
				if d.opt.Age > time.Now().Unix()-t.Unix() {
					continue
				}

				delete(d.value["s"].(map[string]interface{}), id)
			}

			delete(d.value, "to")
		})
	}

	return true
}

// Random ...
func (d *Dup) Random() string {
	return randStringRunes(3)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
