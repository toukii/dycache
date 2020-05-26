package gcache

import (
	// "fmt"
	"testing"
)

func TestBitmapExist(t *testing.T) {
	var ex bool
	bitmap := new(Bitmap)
	for i := 0; i < 100; i++ {
		bitmap.Set(i)
	}
	// fmt.Println(bitmap.bits)
	ex = bitmap.Exist(99)
	t.Logf("%+v", ex)
	for i := 0; i < 90; i++ {
		bitmap.Purge(i)
	}

	for i := 0; i < 100; i++ {
		ex := bitmap.Exist(i)
		t.Logf("%d %+v", i, ex)
	}

}
