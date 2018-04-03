package random

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	m sync.Mutex
	r *rand.Rand
)

func init() {
	src := rand.NewSource(time.Now().UnixNano())
	r = rand.New(src)
}

func New64BitID() string {
	m.Lock()
	defer m.Unlock()

	for id := r.Uint64(); id != 0; {
		return fmt.Sprintf("%016x", id)
	}
	panic("")
}
