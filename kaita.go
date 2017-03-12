package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	Qiita struct {
		AccessToken string `yaml:"access_token"`
	} `yaml:"qiita.com"`
}

type Article struct {
	Id        string
	URL       string `json:"url"`
	Title     string
	Body      string
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "kaita: config not found\n")
		fmt.Fprintf(os.Stderr, "Usage: kaita <config>\n")
		os.Exit(1)
	}

	filepath := os.Args[1]
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kaita: %v\n", err)
		os.Exit(1)
	}

	config := ParseConfig(data)
	articles := GetQiitaArticles(config)

	for _, article := range articles {
		fmt.Printf("- [%s](%s)\n", article.Title, article.URL)
	}
}

func ParseConfig(data []byte) Config {
	config := Config{}

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kaita: %v\n", err)
		os.Exit(1)
	}

	return config
}

func GetQiitaArticles(config Config) []Article {
	req, err := http.NewRequest("GET", "https://qiita.com/api/v2/authenticated_user/items", nil)
	req.Header.Add("Authorization", "Bearer "+config.Qiita.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kaita: %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	var articles []Article
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &articles); err != nil {
		fmt.Fprintf(os.Stderr, "kaita: %v\n", err)
		os.Exit(1)
	}

	return articles
}
