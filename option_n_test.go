package safe

import (
	gc "gopkg.in/check.v1"
)

type MultiOptionSuite struct{}

var _ = gc.Suite(&MultiOptionSuite{})

type testMultiOptions struct {
	S string
	N Option2[string, int8]
	O Option2[string, int8]
	P Option2[*string, int32]
	Q Option2[*string, int32]
}

func newTestMultiStruct(s1, s2 string, s3 *string, i4 int32) testMultiOptions {
	t := testMultiOptions{}
	t.S = s1
	t.N = None2[string, int8]()
	t.O = Some2[string, string, int8](s2)
	t.P = Some2[*string, *string, int32](s3)
	t.Q = Some2[int32, *string, int32](i4)
	return t
}

func (s *OptionSuite) TestMultiOptions(c *gc.C) {
	res2 := Some[*testMultiOptions](nil)
	c.Check(res2.IsNone(), gc.Equals, true)
	switch res2.Unwrap().(type) {
	case *testMultiOptions:
		c.Log("impossible type")
		c.Fail()
	}

	test := newTestMultiStruct("test", "", nil, 42)
	res3 := Some(&test)
	c.Check(res3.IsSome(), gc.Equals, true)

	var test2 testMultiOptions
	switch match := res3.Unwrap().(type) {
	case *testMultiOptions:
		test2 = *match
		c.Check(test2, gc.DeepEquals, test)
	default:
		c.Log("impossible type")
		c.Fail()
	}

	c.Check(test.Q.IsSome(), gc.Equals, true)
	switch match := test.Q.Unwrap().(type) {
	case int32:
		c.Check(match, gc.DeepEquals, int32(42))
	default:
		c.Log("impossible type")
		c.Fail()
	}
}
