package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Theme struct {
	WhiteColor          string `yaml:"white_color"`
	BlackColor          string `yaml:"black_color"`
	MainBackgroundColor string `yaml:"main_background_color"`
	DividerColor        string `yaml:"divider_color"`
	MainDarkColor       string `yaml:"main_dark_color"`
	MainMidColor        string `yaml:"main_mid_color"`
	MainLightColor      string `yaml:"main_light_color"`
	AccentColor         string `yaml:"accent_color"`
	AccentLightColor    string `yaml:"accent_light_color"`
	AccentDarkColor     string `yaml:"accent_dark_color"`
	FontFamily          string `yaml:"font_family"`
}

type Site struct {
	Name          string `yaml:"name"`
	Logo          string `yaml:"logo"`
	LogoWidth     string `yaml:"logo_width"`
	LogoHeight    string `yaml:"logo_height"`
	Pages         Pages  `yaml:"pages"`
	Theme         Theme  `yaml:"theme"`
	Title         string `yaml:"title"`
	FooterContent string `yaml:"footer_content"`
	Favicon       string `yaml:"favicon"`
	Stylesheet    string `yaml:"stylesheet"`
}

type Hero struct {
	Background      string `yaml:"background"`
	Content         string `yaml:"content"`
	SubContent      string `yaml:"sub_content"`
	Image           string `yaml:"image"`
	BackgroundImage string `yaml:"background_image"`
}

type Page struct {
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	HideFromNav bool   `yaml:"hide_from_nav"`
	Breadcrumb  bool   `yaml:"breadcrumb"`
	Comments    bool   `yaml:"comments"`
	Collection  bool   `yaml:"collection"`
	Hero        Hero   `yaml:"hero"`
	Toc         bool   `yaml:"toc"`
}

type Pages map[string]Page

type Config struct {
	ContentDir string `yaml:"content_dir"`
	BaseURL    string `yaml:"base_url"`
	OutputDir  string `yaml:"-"`
	Mode       string `yaml:"-"`
	Port       int    `yaml:"port"`
	Site       Site   `yaml:"site"`
}

func NewConfig() Config {
	cfg := Config{
		BaseURL: "",
		Port:    8080,
		Site: Site{
			Name:          "My Site",
			FooterContent: "Powered by PubGo (https://pubgo.org) - " + "CopyRight © 2019 My Site",
			LogoHeight:    "32px",
			Pages: Pages{
				"0": {
					Name:        "home",
					Path:        "/",
					HideFromNav: true,
					Collection:  false,
					Hero: Hero{
						Background:      "#dcdcdc",
						Content:         "Welcome to my site",
						SubContent:      "This is a simple site",
						Image:           "https://picsum.photos/200",
						BackgroundImage: "https://picsum.photos/400",
					},
				},
				"1": {
					Name:        "blog",
					Path:        "/blog",
					HideFromNav: false,
					Collection:  true,
				},
			},
			Theme: Theme{
				WhiteColor:          "#020202",
				MainBackgroundColor: "#eeeeee",
				BlackColor:          "#eeeeee",
				DividerColor:        "#293437",
				MainDarkColor:       "#fefefe",
				MainMidColor:        "#eeeeee",
				MainLightColor:      "#111111",
				AccentColor:         "#020202",
				AccentLightColor:    "#0f0f0f",
				AccentDarkColor:     "#020202",
				FontFamily:          "monospace",
			},
		},
	}

	return cfg
}

// loadConfig
func LoadConfig(filename string, cfg *Config) {
	data, err := os.ReadFile(filename)

	// If file doesn't exist, create it with default content
	if err != nil {
		log.Printf("\nError reading config file.\n" +
			"\tPlease make sure config.yaml exists and is readable.\n" +
			"\tYou can specify a different config file with the -config flag.\n" +
			"\tContinuing with default config.\n")

		data, err = yaml.Marshal(cfg)

		if err != nil {
			log.Println("Error marshalling default config:", err)
			panic(err)
		}

		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			log.Println("Error creating config file:", err)
			panic(err)
		}
	}

	data, err = os.ReadFile(filename)
	if err != nil {
		log.Println("Error reading config file:", err)
		panic(err)
	}

	err = yaml.Unmarshal(data, cfg)

	if err != nil {
		log.Println("Error unmarshalling config:", err)
		panic(err)
	}
}
