package safe

import (
	"errors"

	gc "gopkg.in/check.v1"
)

type ResultSuite struct{}

var _ = gc.Suite(&ResultSuite{})

type testResults struct {
	S string
	N Result[Option[string]]
	O Result[Option[string]]
	P Result[*string]
	Q Result[*string]
}

func newTestResults(s1, s2 string, s3 *string) testResults {
	t := testResults{}
	t.S = s1
	t.N = Result[Option[string]]{}
	t.O = OK(Some(s2))
	t.P = OK(s3)
	t.Q = Err[*string](errors.New(""))
	return t
}

func (s *ResultSuite) TestResults(c *gc.C) {
	ptr := "adieu"
	doc := newTestResults("so long", "", &ptr)
	res := OK(doc)
	c.Check(res.IsOK(), gc.Equals, true)
	switch match := res.Unwrap().(type) {
	case testResults:
		c.Check(match, gc.DeepEquals, doc)
	default:
		c.Log("impossible type")
		c.Fail()
	}

	res2 := Err[*testResults](nil)
	c.Check(res2.IsErr(), gc.Equals, true)
	switch match := res2.Unwrap().(type) {
	case *testResults:
		c.Log("impossible type")
		c.Fail()
	case error:
		c.Check(match, gc.DeepEquals, errors.New(""))
	}

	res3 := OK(&testResults{S: "test"})
	switch match := res3.Unwrap().(type) {
	case *testResults:
		c.Check(match, gc.DeepEquals, &testResults{S: "test"})
	default:
		c.Log("impossible type")
		c.Fail()
	}
}
