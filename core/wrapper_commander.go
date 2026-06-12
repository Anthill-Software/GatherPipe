package core

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type ExecuteRequest struct {
	Cmd  string
	Args []string
}

// 1. Le Client RPC (Core -> Plugin)
type CommanderRPCClient struct{ client *rpc.Client }

func (c *CommanderRPCClient) SupportedCommands() ([]CommandArg, error) {
	var reply []CommandArg
	err := c.client.Call("Plugin.SupportedCommands", new(interface{}), &reply)
	return reply, err
}

func (c *CommanderRPCClient) ExecuteCommand(cmd string, args []string) (string, error) {
	var reply string
	req := ExecuteRequest{Cmd: cmd, Args: args}
	err := c.client.Call("Plugin.ExecuteCommand", req, &reply)
	return reply, err
}

// 2. Le Serveur RPC (Plugin -> Core)
type CommanderRPCServer struct{ Impl Commander }

func (s *CommanderRPCServer) SupportedCommands(args interface{}, reply *[]CommandArg) error {
	cmds, err := s.Impl.SupportedCommands()
	if err != nil {
		return err
	}
	*reply = cmds
	return nil
}

func (s *CommanderRPCServer) ExecuteCommand(req ExecuteRequest, reply *string) error {
	res, err := s.Impl.ExecuteCommand(req.Cmd, req.Args)
	if err != nil {
		return err
	}
	*reply = res
	return nil
}

// 3. Le Plugin
type CommanderPlugin struct {
	Impl Commander
}

func (p *CommanderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &CommanderRPCServer{Impl: p.Impl}, nil
}

func (p *CommanderPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CommanderRPCClient{client: c}, nil
}
