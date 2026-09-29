package safe

import (
	"testing"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) { gc.TestingT(t) }

type OptionSuite struct{}

var _ = gc.Suite(&OptionSuite{})

type testOptions struct {
	S string
	N Option[string]
	O Option[string]
	P Option[*string]
	Q Option[*string]
}

func newTestStruct(s1, s2 string, s3 *string) testOptions {
	t := testOptions{}
	t.S = s1
	t.N = None[string]()
	t.O = Some(s2)
	t.P = Some(s3)
	t.Q = Some[*string](nil)
	return t
}

func (s *OptionSuite) TestOptions(c *gc.C) {
	res2 := Some[*testOptions](nil)
	c.Check(res2.IsNone(), gc.Equals, true)
	switch res2.Match().(type) {
	case MatchSome[*testOptions]:
		c.Log("impossible type")
		c.Fail()
	}

	res3 := Some(&testOptions{S: "test"})
	c.Check(res3.IsSome(), gc.Equals, true)

	switch match := res3.Match().(type) {
	case MatchSome[*testOptions]:
		c.Check(match.Some, gc.DeepEquals, &testOptions{S: "test"})
	default:
		c.Log("impossible type")
		c.Fail()
	}
}
