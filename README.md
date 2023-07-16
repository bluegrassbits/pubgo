## Introducing PubGo: A Dynamic Content Publishing Framework in Go

PubGo is a lightweight and customizable content publishing framework written in
Go. Inspired by the simplicity of Hugo, PubGo aims to simplify the process of
generating static sites while offering the convenience of a live server approach.
With PubGo, you can easily create, manage, and customize your website through a
simple configuration file.

### Key Features

-   **Customizability:** Update the styling, layout, and pages of your website
    effortlessly using a configuration file. Customize your site's appearance with
    custom CSS and Go templates.
-   **Embedded Webserver:** Unlike Hugo's primarily testing-focused webserver,
    PubGo places a stronger emphasis on providing a robust hypermedia system. As the
    framework matures, PubGo plans to leverage [htmx](https://htmx.org) to deliver dynamic user experiences.
-   **Simplified Content Delivery:** Adding new content to your website is as easy
    as dropping a Markdown file into a directory. PubGo handles the conversion of
    Markdown to HTML on the fly, making content management a breeze.

### Quick Start

To get started with PubGo, follow these steps:

1. Clone the PubGo repository:

    ```bash
    git clone https://github.com/bluegrassbits/pubgo.git
    cd pubgo
    ```

2. Build the PubGo executable:

    ```bash
    go build
    ```

3. Run the PubGo server:

    ```bash
    ./pubgo
    ```

By default, the PubGo server host this site locally at http://localhost:8889.
You can customize the server configuration in the PubGo configuration file.
See the [docs](https://pubgo.org/docs) section for details on bootstrapping a new project.

PubGo simplifies content publishing, allowing you to focus on creating engaging
content while providing a seamless and hassle-free publishing experience. Start
using PubGo today to build and customize your dynamic websites with ease!
