// Inspired from tidb testleak
// https://github.com/pingcap/tidb/tree/master/util/testleak

package testleak

import (
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"
)

// defaultTestLeakWhiteList with default values
var defaultTestLeakWhiteList = []string{"testing.Main(", "runtime.goexit", "testing.(*T).Run"}

// testLeakWhiteList initialize with defaultTestLeakWhiteList
var testLeakWhiteList = defaultTestLeakWhiteList

// SetTestLeakWhiteList will set variable testLeakWhiteList
func SetTestLeakWhiteList(strSlice []string) {
	testLeakWhiteList = strSlice
}

// AppendTestLeakWhiteList will append strSlice to existing testLeakWhiteList
func AppendTestLeakWhiteList(strSlice ...string) {
	testLeakWhiteList = append(testLeakWhiteList, strSlice...)
}

// RestoreDefaultTestLeakWhiteList will set testLeakWhiteList with defaultTestLeakWhiteList
func RestoreDefaultTestLeakWhiteList() {
	testLeakWhiteList = defaultTestLeakWhiteList
}

// containsAnyInStringSlice will return whether s contains any substr in substrSlice
func containsAnyInStringSlice(s string, substrSlice []string) bool {
	for _, substr := range substrSlice {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

func interestingGoroutines() (gs []string) {
	buf := make([]byte, 2<<20)
	buf = buf[:runtime.Stack(buf, true)]
	for _, g := range strings.Split(string(buf), "\n\n") {
		sl := strings.SplitN(g, "\n", 2)
		if len(sl) != 2 {
			continue
		}
		stack := strings.TrimSpace(sl[1])
		if stack == "" || containsAnyInStringSlice(stack, testLeakWhiteList) {
			continue
		}
		gs = append(gs, stack)
	}
	sort.Strings(gs)
	return
}

var beforeTestGorountines = map[string]bool{}

// beforeTest gets the current goroutines.
func beforeTest() {
	for _, g := range interestingGoroutines() {
		beforeTestGorountines[g] = true
	}
}

// TestLeak gets the current goroutines and runs the returned function to
// get the goroutines at that time to contrast whether any goroutines leaked.
// Usage: defer testleak.TestLeak(t)()
func TestLeak(t *testing.T) func() {
	if len(beforeTestGorountines) == 0 {
		beforeTest()
	}

	return func() {
		defer func() {
			beforeTestGorountines = map[string]bool{}
		}()

		var leaked []string
		for i := 0; i < 50; i++ {
			leaked = leaked[:0]
			for _, g := range interestingGoroutines() {
				if !beforeTestGorountines[g] {
					leaked = append(leaked, g)
				}
			}
			// Bad stuff found, but goroutines might just still be
			// shutting down, so give it some time.
			if len(leaked) != 0 {
				time.Sleep(50 * time.Millisecond)
				continue
			}

			return
		}
		for _, g := range leaked {
			t.Errorf("Test appears to have leaked: %v", g)
		}
	}
}
