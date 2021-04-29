package web

import (
	"net/http"

	"github.com/bennycio/bundle/api"
)

func getProfileFromCookie(r *http.Request) (Profile, error) {

	c, err := r.Cookie("access_token")
	if err != nil {
		return Profile{}, err
	}
	user, err := getUserFromToken(c.Value)
	if err != nil {
		return Profile{}, err
	}
	return userToProfile(user), nil

}

func userToProfile(user *api.User) Profile {
	return Profile{
		Username: user.Username,
		Email:    user.Email,
		Tag:      user.Tag,
		Scopes:   user.Scopes,
	}
}

func pluginToInfo(plugin *api.Plugin) PluginInfo {
	return PluginInfo{
		Name:        plugin.Name,
		Version:     plugin.Version,
		Author:      plugin.Author,
		Description: plugin.Description,
	}
}

func pluginsToInfos(plugins []*api.Plugin) []PluginInfo {
	result := []PluginInfo{}
	for _, v := range plugins {
		result = append(result, pluginToInfo(v))
	}
	return result
}
