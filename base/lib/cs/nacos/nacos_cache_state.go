package nacos

type cacheState struct {
	state bool
}

func (c *cacheState) active() bool {
	return c.state
}

func (c *cacheState) enable() {
	c.state = true
}

func (c *cacheState) disable() {
	c.state = false
}
