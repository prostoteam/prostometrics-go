package prostometrics

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"
)

func newSessionID() string {
	var data [16]byte
	if _, err := rand.Read(data[:]); err == nil {
		return hex.EncodeToString(data[:])
	}
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
