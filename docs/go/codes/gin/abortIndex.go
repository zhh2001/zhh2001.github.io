const abortIndex int8 = math.MaxInt8 >> 1

func (c *Context) Abort() {
	c.index = abortIndex
}