package app

import "fmt"

type Commands struct {
	CmdHandlerMap map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CmdHandlerMap[name] = f
}
func (c *Commands) Run(s *State, cmd Command) error {
	if handler, exists := c.CmdHandlerMap[cmd.Name]; exists {
		err := handler(s, cmd)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("command does not exist")
	}
	return nil
}
