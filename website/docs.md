---
title: PubGo Documentation
include_toc: true
---

# Documentation
---

Welcome to the PubGo documentation. This section provides comprehensive
information on how to install, configure, and use PubGo for your content
publishing needs.

<figure>
PubGo is currently in a prototype state and should be considered experimental.
Breaking changes will be frequent at this stage. With that being said, it is
somewhat functional and can be used to generate small projects and tinkerings.
</figure>

## Getting Started
### Basics

Assuming you have `git` and `go` installed, you can install **pubgo** on your
local machine or server by `cloning`, `building`, and `running` the resulting
binary.

It is also convienient to add the binary that is someplace in your **$PATH**.

If you run **pubgo** without any commandline flags, it will look for a
`config.yaml` file in the directory in which it was executed. If that file does
not exist, it will attempt to create it.

There is a `config.yaml` at the root of the `pubgo` repository so if you run
**pubgo** from within the repo, it will serve up this website locally.


```bash
# clone, build, run
git clone https://github.com/bluegrassbits/pubgo.git
cd pubgo
go build
./pubgo
```

### Bootstrapping
<figure>
If you'd prefer starting from scratch, you can specify a config and content
directory via commandline arguments. If pubgo is in your <b>$PATH</b> as
mentioned before, you can just <chip class="info">cd</chip> into your project's
directory and run pubgo with no arguments to bootstrap your new project.
</figure>

```bash
# bootstraping a new project
mkdir new_project
./pubgo -content_dir ./new_project -config ./new_project/config.yaml
```
This will result in a very plain starter. You can get started hacking on a
custom stylesheet by specifying a stylesheet path in your `config.yaml` and
creating a stylesheet in `./<content_dir>/static/`

Another option would be to download or specify a remote stylesheet. For this
platform, [missing.css](https://missing.style/) is recommended as a good
starter.

```yaml
# config.yaml
site:
    stylesheet: "/static/custom.css"
# or
site:
    stylesheet: "https://unpkg.com/missing.css@1.0.9/dist/missing.min.css"
```

## Configuration
Understanding the configuration file structure and how to customize PubGo for your
specific requirements.

~Commandline Flags~

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

~Config Example~

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
        sytnax_highlight: false
        bg: "#eeeeee"
        fg: "#020202"
        accent: "#020202"
        muted_accent: "#020202"
        main_font: monospace
    title: ""
    footer_content: CopyRight © 2019 My Site
    favicon: ""
    stylesheet: ""
```

**NOTE:** Any static files used in the configuration should be placed in
`<content_dir>/static`. If this directory doesn't exist, pubgo will attempt to
create it on startup.

### Main Config

| Option                      | type       | Description                                                                                                     | Default |
| --------------------------- | ---------- | --------------------------------------------------------------------------------------------------------------- | ------- |
| **content_dir**             | string     | path to user supplied data                                                                                      | "."     |
| **port**                    | string/int | listening port for server                                                                                       | 8080    |
| **site**                    | string     | site specific nested config                                                                                     | ~N/A~   |
| **site.name**           | string     | site name, used for header and title                                                                            | PubGo   |
| **site.logo**           | string     | path to logo image. used for header if present. path is relative to `content_dir`. should start with `/static/` |         |
| **site.logo_width**     | string     | width of _logo_ image in header                                                                                 | "32px"  |
| **site.logo_height**    | string     | height of _logo_ image in header                                                                                |         |
| **site.title**          | string     | tagline appended to site name in main page title                                                                |         |
| **site.footer_content** | string     | text of footer                                                                                                  |         |
| **site.favicon**        | string     | path to favicon image. used if present. path is relative to `content_dir`. should start with `/static/`         |         |
| **site.stylesheet**     | string     | path to stylesheet. used if present. path is relative to `content_dir`. should start with `/static/`. can also reference a remote stylesheet         |         |

### Site Config

In addition to the options listed above, there are a couple of nested options.
These are the **site.theme** and the **site.pages**.

For **site.theme**, please reference the example config above. Most of it
should be self-explanatory. This page provides a way to do some basic
customization.


| Option | type | Description | Default |
| ------ | ---- | ----------- | ------- |
| **site.theme.syntax_highlight** | bool | whether or not to apply syntax highlighting on code blocks | false |
| **site.theme.syntax_theme** | string | syntax highlighting theme | "dracula" |

`Note:` see https://github.com/alecthomas/chroma/tree/master/styles for a list of usable themes

### Pages Config

This is were the pages of your site can be setup. **sites.pages** is
a _key/value_ data structure. The _key_ here should be a numeric value. This
value is used for sorting, specifically the sort order of pages that are
included in the header navigation.

The _value_ portion of this is the configuration for a single page.

~Page Config~

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

#### Page Config Details

An important piece of page configuration is the **collection** bit. If _true_,
then the page is expected to display multiple entries of content.

-   Collections pages will look for content in **\<content_dir\>/\<page_name\>/**
-   Non-Collections pages expect **\<content_dir\>/\<page_name\>.md**


#### Other page options

| Option            | Type   | Description                                                           |
| ----------------- | ------ | --------------------------------------------------------------------- |
| **name**          | string | the page name. for not "/" routes this is used for the page \<title\> |
| **path**          | string | this is used in route handling configuration. for collections         |
| **hide_from_nav** | bool   | whether or not to hide page from navbar                               |
| **collection**    | bool   | whether or not this page contains multiple entries                    |
| **hero**          | object | page hero configuration, see example above for options                |


## Usage Guide
### Creating Content
<figure class="info">
<p>
How to create different types of content, such as blog posts, articles, or portfolio pages.
</p>
<strong>This Section is WIP</strong>
</figure>


### Managing Pages
<figure class="info">
<p>
How to modify and organize pages within your website through the configuration file.
</p>
<strong>This Section is WIP</strong>
</figure>

### Customization
<figure class="info">
<p>
Customizing the design, layout, and themes of your PubGo website.
</p>
<strong>This Section is WIP</strong>
</figure>


## Deployment

<figure>
As of now, there isn't any "official" deployment process or pattern for pubgo.
We hope that it is portable enough to provide for maximum flexiblity. That
being said, we'll explore a few potential options.
</figure>

### Docker

You can use the following to build a docker container to host your content.

```bash
mkdir new_project

docker build -t pubgo \
    --build-arg UID=$UID \
    --build-arg GID=$UID .

docker run -it -p 8080:8080  -v ./new_project:/opt/pubgo pubgo
```

### Service

Service configurations do not exist right now but examples will likely be
provided in the future. It should be pretty straight forward to use `init.d`,
`systemd`, or whatever service manager is on your server.

You could always YOLO and just run it in a `tmux` or `sreen` session on your
server. 😚

If you have ideas on how we can make **pugbo** more service manager friendly,
suggestions are welcome in the issues section on the repo's
[Github page](https://github.com/bluegrassbits/pubgo).

### Static Site Compilation

Static site generation is always an option. This is a fantastic option for many
projects, and honestly at this point, pubgo in server mode doesn't offer many
advantages over a statically generated site. It is our hope that as this
platform matures this will change.

To generate a static site, use the `build` mode:

```bash
./pubgo -mode build -content_dir ./website -out ./out
```

## Todo

-   [ ] improve server logging
-   [ ] explore building flatfile commenting into pubgo
-   [ ] improve table of content presentation
-   [ ] evaluate duplication between markdown parser and Page type, refactor
-   [ ] entries sorting and listing options, possibly pagination
-   [ ] work on image based content publishing and presentation
-   [x] debug comments not showing for posts loaded dynamically with htmx
-   [x] cleanup default site theme
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
-   [x] figure out how to remove the need for `br` tags in markdown
-   [x] create docker deployment
-   [x] move example to website, update config. create examples dir with content
-   [x] setup 404 page
-   [x] meta/seo/robot/sitemap stuff(added headMeta template with example of override. will setup better configuration options in the future.
