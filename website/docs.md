---
title: PubGo Documentation
include_toc: true
---

<br/>
# Documentation

Welcome to the PubGo documentation. This section provides comprehensive
information on how to install, configure, and use PubGo for your content
publishing needs.

## Getting Started

### Basics

How to install PubGo on your local machine or server.

Assuming you have `git` and `go` installed:

```bash
git clone https://github.com/bluegrassbits/pubgo.git
cd pubgo
go build
./pubgo
```

You can also copy this binary someplace thats in your `$PATH` on your system.

`config.yaml` is the default configuration filename. There is a `config.yaml`
at the root of the `pubgo` repository. This will serve up this website locally.

If you'd like to start from scratch:

```bash
mkdir new_project
cp pubgo new_project/
# or
go build -o new_project/
cd new_project
./pubgo
```

This will bootstrap the current directory with a boilerplate config file,
static content directory, and a `home.md`.

### Docker

You can use the following to build a docker container to host your content.

```bash
mkdir new_project

docker build -t pubgo \
    --build-arg UID=$UID \
    --build-arg GID=$UID .

docker run -it -p 8080:8080  -v ./new_project:/opt/pubgo pubgo
```

### Static site compiling

To generate a static site, use the `build` mode:

```bash
./pubgo -mode build -content_dir ./website -out ./out
```

## Configuration:

Understanding the configuration file structure and how to customize PubGo for your
specific requirements.

Commandline flags:

```bash
$ ./pubgo -h
Usage of ./pubgo:
  -config string
        Path to config file (default "config.yaml")
  -content_dir string
        Content directory (default "./website")
  -mode string
        Run mode: <serve> or <build> static site (default "serve")
  -out string
        Output directory for static site (default "./out")
```

The from scratch instructions above should result in a config file that looks
something like this:

### Config example

```yaml
# config.yaml
content_dir: ./website
base_url: ""
port: 8080
site:
    name: My Site
    logo: ""
    logo_width: ""
    logo_height: 32px
    pages:
        "0":
            name: home
            path: /
            hide_from_nav: true
            collection: false
            hero:
                background: "#dcdcdc"
                content: Welcome to my site
                sub_content: This is a simple site
                image: https://picsum.photos/200
                background_image: https://picsum.photos/400
        "1":
            name: blog
            path: /blog
            hide_from_nav: false
            collection: true
            hero:
                background: ""
                content: ""
                sub_content: ""
                image: ""
                background_image: ""
    theme:
        background_color: "#eeeeee"
        text_color: "#020202"
        main_color: "#020202"
        accent_color: "#020202"
        font_family: monospace
    title: ""
    footer_content: Powered by PubGo (https://pubgo.org) - CopyRight © 2019 My Site
    favicon: ""
    stylesheet: ""
```

**NOTE:** Any static files used in the configuration should be placed in
`<content_dir>/static`. If this directory doesn't exist, pubgo will attempt to
create it on startup.

### Main config breakdown

| Option                      | type       | Description                                                                                                    | Default |
| --------------------------- | ---------- | -------------------------------------------------------------------------------------------------------------- | ------- |
| **config_dir**              | string     | path to user supplied data                                                                                     | "."     |
| **port**                    | string/int | listening port for server                                                                                      | 8080    |
| **site**                    | string     | site specific nested config                                                                                    | ~N/A~   |
| **site**>**name**           | string     | site name, used for header and title                                                                           | PubGo   |
| **site**>**logo**           | string     | path to logo image. used for header if present. path is relative to `config_dir`. should start with `/static/` |         |
| **site**>**logo_width**     | string     | width of _logo_ image in header                                                                                | "32px"  |
| **site**>**logo_height**    | string     | height of _logo_ image in header                                                                               |         |
| **site**>**title**          | string     | tagline appended to site name in main page title                                                               |         |
| **site**>**footer_content** | string     | text of footer                                                                                                 |         |
| **site**>**favicon**        | string     | path to favicon image. used if present. path is relative to `config_dir`. should start with `/static/`         |         |
| **site**>**stylesheet**     | string     | path to favicon image. used if present. path is relative to `config_dir`. should start with `/static/`         |         |

### Site config breakdown

In addition to the options listed above, there are a couple of nested options.
These are the **site**>**theme** and the **site**>**pages**.

For **site**>**theme**, please reference the example config above. Most of it
should be self-explanatory. This page provides a way to do some basic
customization.

### Site > Pages config breakdown

This is were the pages of your site can be setup. **sites**>**pages** is
a _key/value_ data structure. The _key_ here should be a numeric value. This
value is used for sorting, specifically the sort order of pages that are
included in the header navigation.

The _value_ portion of this is the configuration for a single page.

### Page config

```yaml
pages:
    name: home
    path: /
    hide_from_nav: true
    collection: false
    hero:
        background: ""
        content: ""
        sub_content: ""
        image: ""
        background_image: ""
```

#### Page config details

An important piece of page configuration is the **collection** bit. If _true_,
then the page is expected to display multiple entries of content.

-   Collections pages will look for content in **\<content_dir\>/\<page_name\>/**
-   Non-Collections pages expect **\<content_dir\>/\<page_name\>.md**

<br/>

#### Other page options

| Option            | Type   | Description                                                           |
| ----------------- | ------ | --------------------------------------------------------------------- |
| **name**          | string | the page name. for not "/" routes this is used for the page \<title\> |
| **path**          | string | this is used in route handling configuration. for collections         |
| **hide_from_nav** | bool   | whether or not to hide page from navbar                               |
| **collection**    | bool   | whether or not this page contains multiple entries                    |
| **hero**          | object | page hero configuration, see example above for options                |

<br/>
## Usage Guide

### Creating Content:

How to create different types of content, such as blog posts, articles, or portfolio pages.

This Section is WIP

### Managing pages:

How to modify and organize pages within your website through the configuration file.

This Section is WIP

### Customization:

Customizing the design, layout, and themes of your PubGo website.

This Section is WIP

### Deployment:

This Section is WIP

## Todo

-   [x] figure out how to overwrite templates with user defined templates
-   [x] intial intro and docs pages
-   [x] change default header text when no img present
-   [x] add to nav should be default. setting should toggle if we dont want it
-   [x] if collection page is configured but folder doesn't exist in content dir create it
-   [x] if home.md doesn't exist in content dir create it with some default content
-   [x] ~~make make home page called index globally~~
-   [x] ~~auto-detect new content~~ new router will lookup entries on request
-   [x] Test clean install
-   [x] Create single page table of content style pages
-   [x] allow enable table of content on individual collection entries
-   [ ] improve table of content presentation
-   [ ] evaluate duplication between markdown parser and Page type, refactor
-   [ ] entries sorting and listing options, possibly pagination
-   [x] figure out how to remove the need for `br` tags in markdown
-   [x] create docker deployment
-   [x] move example to website, update config. create examples dir with content
-   [x] setup 404 page
-   [ ] work on image based content publishing and presentation
-   [x] meta/seo/robot/sitemap stuff(added headMeta template with example of override. will setup better configuration options in the future.
-   [ ] improve default site theme
-   [ ] debug comments not showing for posts loaded dynamically with htmx
-   [ ] improve server logging
-   [ ] explore building flatfile commenting into pubgo
