package safe

import (
	"encoding/json"

	gc "gopkg.in/check.v1"
)

type OptionSerdeSuite struct{}

var _ = gc.Suite(&OptionSerdeSuite{})

func (s *OptionSerdeSuite) TestJSON(c *gc.C) {
	ptr := "adieu"
	doc := newTestStruct("so long", "", &ptr)
	newDoc := testOptions{} // do not use zero literals in production

	buf := []byte(`{"S": "so long", "N": {}, "O": "", "P": "adieu", "Q": {}}`)
	err := json.Unmarshal(buf, &newDoc)
	c.Check(err, gc.IsNil)
	c.Check(newDoc, gc.DeepEquals, doc)

	buf, err = json.Marshal(doc)
	c.Check(err, gc.IsNil)
	if err != nil {
		c.Log(err)
	}
	c.Log(string(buf))
	err = json.Unmarshal(buf, &newDoc)
	c.Assert(err, gc.IsNil)
	c.Check(newDoc, gc.DeepEquals, doc)
}
