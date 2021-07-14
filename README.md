# BundleMC

# NOTE: project is a work in progress and not currently available for download

## The ultimate solution to Spigot plugin management

BundleMC is a modern solution to plugin management allowing users to spend less days housekeeping and more days enjoying Minecraft with friends. BundleMC enables server owners to download, update, and manage a server's spigot plugins and enables plugin developers to effortlessly upload their products to the cloud for anyone to access. BundleMC has a variety of interfaces and is flexible to adapt to the needs of any server owner or developer

## Features

- Update all plugins within seconds
- etc.

## Installation

All downloads of the various products associated with BundleMC can found in the releases tab ->

## Official Website

https://bundlemc.io/

Search through available plugins, manage your uploaded plugins, find tutorials and downloads, and more.

# Apps

Below are various applications associated with BundleMC that make it easy for you to find a solution best tailored for your server or your development process.

## Command-Line Interface

The command-line interface is an easy to use interface with commands in your terminal just like you would use in Minecraft. Use this interface to download plugins, update plugins, boostrap a server, upload plugins to the cloud, and more.

### For Server Owners

_Note: If you do not currently have a server setup on your machine and would like to bootstrap the creation of a server, open a terminal where you would like to create a server folder and type_ `bundle boostrap`

To get started, you are going to want to first initialize a `bundle.yml` file. This file will store information on what plugins you have and the versions they are running in a format you are probably familiar with. If not, check out this tutorial on YML files [here]. You can create this file yourself at the root of your server folder (the top level directory where you server jar is and your plugins folder is) and use the format as described [here] or you can open a new terminal in the root of your server folder and type the following command:

```
bundle init
```

Once the file is initialized, you can open the file, which will look something like this:

```yml
Plugins:
  PluginName: "1.0.0"
```

Here is where you are going to specify which plugins you would like on your server and the version of that plugin that should be downloaded.

_Note: Use the keyword_ `latest` _instead of a version number to specify that you would always like to have the latest version of the plugin downloaded._

_Note: The version must always be surrounded in quotations_

Lets say you would like to have the plugins, EssentialsX, WorldEdit, and Vault on your server, you might make your `bundle.yml` file look like this:

```yml
Plugins:
  EssentialsX: "latest"
  WorldEdit: "latest"
  Vault: "latest"
```

_Note: Plugin names are case-insensitive meaning you do not need to use correct capitilization_

You may then type the following command in your terminal to install each of your plugins (Make sure your terminal working directory is in the root folder of your server with the bundle file):

```
bundle install
```

Once the command completes, you will see in your plugins folder your newly installed plugins.

To get a list of commands you can use with the command-line interface, simply type `bundle` into your terminal. Some common commands you might use are listed below:

- `bundle update`
- `bundle uninstall`

### For Developers

Developing with Bundle is easy! The only prerequisite to utilizing the command-line interface is to have an account on [our official website](https://bundlemc.io/). Once you have made an account and you have built your plugin into a jar with a valid plugin.yml file, you can upload your plugin to the official Bundle Repository by typing the following command:

```
bundle upload [path to plugin jar]
```

Easy as that! You will even get a link to your plugin's new web page! If you would like to add a description to your plugin, you can use the same command to upload a README file (must be in the .md format, similar to a GitHub README file). You may also manage a plugin's description and more on the web page generated for your plugin.

## Licensed Under the MIT License
