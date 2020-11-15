package ms2k

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}
