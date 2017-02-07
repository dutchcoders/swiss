package cmd

var (
	plugins = Plugins{}
)

type Plugins map[string]Plugin

type Plugin interface {
	Do(*Environment) error
}

func Register(name string, plugin Plugin) Plugin {
	plugins[name] = plugin
	return plugin
}
