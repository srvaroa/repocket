package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/kelseyhightower/envconfig"
	"github.com/srvaroa/repocket/pkg/pocket"
)

type config struct {
	ConsumerKey string `required:"true" split_words:"true"`
	AccessToken string
	OutputDir   string `required:"true" split_words:"true"`
}

const RepocketConfigFile = ".repocket"

func ensureDir(path string) {
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err != nil {
			log.Fatalf("ERROR: expecting directory %s to exist", path)
		}
	}
	if !f.IsDir() {
		log.Fatalf("ERROR: expecting path %s to be a directory", path)
	}
}

func dumpArticle(outputDir string, a *pocket.Article) {

	// Clean up path
	re := regexp.MustCompile(`[\.|/\\]+`)
	path := outputDir + "/" +
		a.ItemId +
		"_" +
		string(re.ReplaceAll([]byte(a.ResolvedTitle), []byte("-")))

	// If the article is there, leave it alone
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		log.Printf("Skipping (already downloaded): %s", a.ResolvedTitle)
		return
	}

	log.Printf("Downloading: `%s` to `%s`", a.ResolvedTitle, path)

	file, err := os.Create(path)
	if err != nil {
		log.Printf("Failed to create file for %s: %s", a.ResolvedTitle, err)
		return
	}
	defer file.Close()

	out, err := exec.Command("links2", "-dump", a.ResolvedUrl).Output()
	if err != nil {
		log.Print("Failed to download %s, %s", a.ResolvedUrl, err)
		return
	}

	_, err = io.WriteString(file, string(out))
	if err != nil {
		log.Printf("Failed to write %s: %s", a.ResolvedTitle, err)
	}
}

func main() {

	cfg := config{}

	err := envconfig.Process("REPOCKET", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	ensureDir(cfg.OutputDir)
	cfg.AccessToken, err = pocket.Authorize(cfg.ConsumerKey)
	if err != nil {
		log.Fatal("Failed to authorize against Pocket", err)
	}

	list := pocket.QueryFavourites(cfg.AccessToken, cfg.ConsumerKey)
	for _, item := range list {
		dumpArticle(cfg.OutputDir, &item)
	}

}
