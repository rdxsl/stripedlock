package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestLocksField(size uint32) []sync.Mutex {
	locks := make([]sync.Mutex, size)
	return locks
}

func Test_stripedLock_idToIndex(t *testing.T) {
	type fields struct {
		size  uint32
		locks []sync.Mutex
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		{"id is dog", fields{256, getTestLocksField(256)}, args{"dog"}, 9},
		{"id is cat", fields{256, getTestLocksField(256)}, args{"cat123123ss"}, 237},
		{"id is globalqa.las2.test1.test.app", fields{256, getTestLocksField(256)}, args{"globalqa.las2.test1.test.app"}, 184},
		{"id is globalqa.las2.test2.test.app", fields{256, getTestLocksField(256)}, args{"globalqa.las2.test2.test.app"}, 167},
		{"id is empty string", fields{256, getTestLocksField(256)}, args{""}, 197},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := &stripedLock{
				size:  tt.fields.size,
				locks: tt.fields.locks,
			}
			if got := sl.idToIndex(tt.args.id); got != tt.want {
				t.Errorf("stripedLock.idToIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_StripedLockInterface(t *testing.T) {
	sl, internals := getTestStripedLock(256)
	lock := sl.Get("dog")
	assert.Equal(t, lock, &internals.locks[8])

	// TODO: think about how to unit test functions with no return val
	sl.Lock("dog")
	sl.Unlock("dog")

	ids := []string{"dog", "cat", "fish"}
	sl.BatchLock(ids)
	for _, id := range ids {
		sl.Unlock(id)
	}

	// we made it to the end i guess?
	assert.Equal(t, true, true)
}

func getTestStripedLock(size uint32) (StripedLock, *stripedLock) {
	sl := &stripedLock{
		size:  size,
		locks: make([]sync.Mutex, size),
	}
	return sl, sl
}

func TestNewStripedLock(t *testing.T) {
	sl := NewStripedLock(256)
	assert.NotNil(t, sl)
}
