package register

type PluginServer struct {
	Impl Plugin
}

func (ps *PluginServer) About(args any, resp *string) error {
	*resp = ps.Impl.About()
	return nil
}
