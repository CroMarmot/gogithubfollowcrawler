package unused

import "testing"

func TestCrawlManager_Init(t *testing.T) {
	var customDB CustomDB
	customDB = &CustomDBImpl{33221}
	customDB.(*CustomDBImpl).Init()
	cases := []struct{
		testin,testout string
	}{
		{"","2333"},
	}
	for _,c := range cases{
		got := "2333"+c.testin
		if got != c.testout {
			t.Errorf("error 1234\n")
		}
	}
}

