package internal

import "github.com/bennycio/bundle/api"

type UserService interface {
	Get(req *api.GetUserRequest) (*api.User, error)
	Update(req *api.UpdateUserRequest) error
	Insert(user *api.User) error
}

type PluginService interface {
	Get(req *api.GetPluginRequest) (*api.Plugin, error)
	Insert(plugin *api.Plugin) error
	Update(req *api.UpdatePluginRequest) error
	Paginate(req *api.PaginatePluginsRequest) (*api.PaginatePluginsResponse, error)
}
