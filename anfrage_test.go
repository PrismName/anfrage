package anfrage

import (
	"testing"
)

func Test_Get(t *testing.T) {
	_, err := Get("http://www.baidu.com")
	if err != nil {
		t.Fatalf("error %+v", err)
	}
}
