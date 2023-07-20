package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Theme struct {
	BackgroundColor string `yaml:"background"`
	Fg              string `yaml:"text"`
	Bg              string `yaml:"main"`
	Accent          string `yaml:"accent"`
	MutedAccent     string `yaml:"muted_accent"`
	MainFont        string `yaml:"font_family"`
	SyntaxHighlight bool   `yaml:"syntax_highlight"`
	SyntaxTheme     string `yaml:"syntax_theme"`
}

type Site struct {
	Name          string `yaml:"name"`
	Logo          string `yaml:"logo"`
	LogoText      string `yaml:"logo_text"`
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
	Collection  bool   `yaml:"collection"`
	Hero        Hero   `yaml:"hero"`
}

type Pages map[string]Page

type Config struct {
	ContentDir string `yaml:"content_dir"`
	BaseURL    string `yaml:"base_url"`
	OutputDir  string `yaml:"-"`
	AdminUser  string `yaml:"admin_user"`
	AdminPass  string `yaml:"admin_pass"`
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
			FooterContent: "CopyRight © 2019 My Site",
			LogoHeight:    "32px",
			Pages: Pages{
				"0": {
					Name:        "home",
					Path:        "/",
					HideFromNav: true,
					Collection:  false,
					Hero: Hero{
						Content:    "Welcome to my site",
						SubContent: "This is a simple site",
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
				Bg:              "#f8fafb",
				Fg:              "#020202",
				Accent:          "#020202",
				MutedAccent:     "#020202",
				SyntaxHighlight: false,
				SyntaxTheme:     "dracula",
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
