package readline

import "testing"

func Test_insertRune(t *testing.T) {
	type pair struct {
		cl        cmdLine
		add       rune
		expectStr string
		expectPos int
	}

	tests := []pair{
		{cmdLine{}, '测', "测", 1},                  // empty
		{cmdLine{[]rune("试了"), 0}, '测', "测试了", 1}, // middle
		{cmdLine{[]rune("测了"), 1}, '试', "测试了", 2}, // middle
		{cmdLine{[]rune("测试"), 2}, '了', "测试了", 3}, // end
	}

	for i, test := range tests {
		old := test.cl
		cl := &test.cl
		cl.insertRune(test.add)
		if cl.String() != test.expectStr || cl.curPos != test.expectPos {
			t.Errorf("[%d]: %s + [%s] = %s, expect: (%s, %d)\n", i,
				old.toString(), string(test.add), cl.toString(),
				test.expectStr, test.expectPos)
		}
	}
}

func Test_deleteRune(t *testing.T) {
	type pair struct {
		cl        cmdLine
		expectStr string
		expectPos int
	}

	tests := []pair{
		{cmdLine{}, "", 0},                    // empty
		{cmdLine{[]rune("测试了"), 0}, "试了", 0},  // head
		{cmdLine{[]rune("测试了"), 1}, "测了", 1},  // middle
		{cmdLine{[]rune("测试了"), 3}, "测试了", 3}, // end
	}

	for i, test := range tests {
		old := test.cl
		cl := &test.cl
		cl.deleteRune()
		if cl.String() != test.expectStr || cl.curPos != test.expectPos {
			t.Errorf("[%d]: %s delete = %s, expect: (%s, %d)\n", i,
				old.toString(), cl.toString(), test.expectStr, test.expectPos)
		}
	}
}

func Test_backspaceRune(t *testing.T) {
	type pair struct {
		cl        cmdLine
		expectStr string
		expectPos int
	}

	tests := []pair{
		{cmdLine{}, "", 0},                    // empty
		{cmdLine{[]rune("测试了"), 0}, "测试了", 0}, // head
		{cmdLine{[]rune("测试了"), 1}, "试了", 0},  // middle
		{cmdLine{[]rune("测试了"), 3}, "测试", 2},  // end
	}

	for i, test := range tests {
		old := test.cl
		cl := &test.cl
		cl.backspaceRune()
		if cl.String() != test.expectStr || cl.curPos != test.expectPos {
			t.Errorf("[%d]: %s delete = %s, expect: (%s, %d)\n", i,
				old.toString(), cl.toString(), test.expectStr, test.expectPos)
		}
	}
}

func Test_forwardCursor(t *testing.T) {
	type pair struct {
		cl        cmdLine
		expectStr string
		expectPos int
	}

	tests := []pair{
		{cmdLine{}, "", 0},                    // empty
		{cmdLine{[]rune("测试了"), 0}, "测试了", 1}, // head
		{cmdLine{[]rune("测试了"), 1}, "测试了", 2}, // middle
		{cmdLine{[]rune("测试了"), 3}, "测试了", 3}, // end
	}

	for i, test := range tests {
		old := test.cl
		cl := &test.cl
		cl.forwardCursor()
		if cl.String() != test.expectStr || cl.curPos != test.expectPos {
			t.Errorf("[%d]: %s delete = %s, expect: (%s, %d)\n", i,
				old.toString(), cl.toString(), test.expectStr, test.expectPos)
		}
	}
}

func Test_backwardCursor(t *testing.T) {
	type pair struct {
		cl        cmdLine
		expectStr string
		expectPos int
	}

	tests := []pair{
		{cmdLine{}, "", 0},                    // empty
		{cmdLine{[]rune("测试了"), 0}, "测试了", 0}, // head
		{cmdLine{[]rune("测试了"), 1}, "测试了", 0}, // middle
		{cmdLine{[]rune("测试了"), 3}, "测试了", 2}, // end
	}

	for i, test := range tests {
		old := test.cl
		cl := &test.cl
		cl.backwardCursor()
		if cl.String() != test.expectStr || cl.curPos != test.expectPos {
			t.Errorf("[%d]: %s delete = %s, expect: (%s, %d)\n", i,
				old.toString(), cl.toString(), test.expectStr, test.expectPos)
		}
	}
}
