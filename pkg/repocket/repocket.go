package repocket

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/srvaroa/repocket/pkg/pocket"
	"github.com/srvaroa/repocket/pkg/util"
)

// GetArticleIds finds all the ids of articles cached in a given dir
func GetArticleIds(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to list files in %s %s", dir, err)
	}

	var ids []string
	for _, f := range files {
		if !f.IsDir() {
			pieces := strings.Split(f.Name(), "_")
			ids = append(ids, pieces[0])
		}
	}

	return ids
}

func DumpArticle(outputDir string, a *pocket.Article) {

	title := a.ResolvedTitle
	if len(a.ResolvedTitle) <= 0 && len(a.GivenTitle) > 0 {
		title = a.GivenTitle
	}

	// Clean up path
	re := regexp.MustCompile(`[\.|/\\]+`)
	path := outputDir + "/" +
		a.ItemId +
		"_" +
		string(re.ReplaceAll([]byte(title), []byte("-")))

	// If the article is there, leave it alone
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		log.Printf("Skipping (already downloaded): %s", title)
		return
	}

	log.Printf("Downloading: `%s` to `%s`", title, path)

	// We add the metadata from Pocket to the local copy
	metaBytes, err := yaml.Marshal(a)
	if err != nil {
		log.Printf("Failed to serialize article meta: %s %s", title, err)
	}

	// Prepare the contents of the article
	txtBytes, err := util.DumpUrl(a.ResolvedUrl)
	if err != nil {
		log.Print("Failed to download %s, %s", a.ResolvedUrl, err)
		return
	}

	// Open and write the file
	file, err := os.Create(path)
	if err != nil {
		log.Printf("Failed to create file for %s: %s", title, err)
		return
	}
	defer file.Close()

	_, err = io.WriteString(file, string(metaBytes)+"\n\n---\n\n"+string(txtBytes))
	if err != nil {
		log.Printf("Failed to write %s: %s", title, err)
	}
}

// SyncDeleted reads the deleted directory and marks the articles as
// deleted upstream, then removes all the files in the deleted
// directory
func (r *Repocket) SyncDeletions() {
	ids := GetArticleIds(r.DeletedDir)
	log.Printf("Will delete file: %s", ids)
	pocket.Delete(r.AccessToken, r.ConsumerKey, ids)
	util.EmptyDir(r.DeletedDir)
}

// SyncArchived reads the archived directory and marks the articles as
// archived upstream, then removes all the files in the archived
// directory
func (r *Repocket) SyncArchived() {
	ids := GetArticleIds(r.ArchivedDir)
	log.Printf("Will delete file: %s", ids)
	pocket.Archive(r.AccessToken, r.ConsumerKey, ids)
	util.EmptyDir(r.ArchivedDir)
}

// SyncFavs does ths following steps:
// * Keep a local copy of all articles fav'd upstream
// * Ensure that all articles in the local fav dir are fav'd and
//   archived upstream
func (r *Repocket) SyncFavs() {

	if len(r.FavsDir) == 0 {
		log.Fatalf("No output directory provided")
	}
	util.EnsureDir(r.FavsDir)

	knownIds := map[string]bool{}
	favs := pocket.QueryFavourites(r.AccessToken, r.ConsumerKey)
	for _, item := range favs {
		DumpArticle(r.FavsDir, &item)
		knownIds[item.ItemId] = true
	}

	// Find articles that were NOT in the favourited list
	var missingIds []string
	for _, id := range GetArticleIds(r.FavsDir) {
		if !knownIds[id] {
			missingIds = append(missingIds, id)
		}
	}

	log.Printf("Uploading %d new favs %s", len(missingIds), missingIds)
	pocket.Fav(r.AccessToken, r.ConsumerKey, missingIds)
	pocket.Archive(r.AccessToken, r.ConsumerKey, missingIds)
}

// SyncUnread reads all articles not archived and stores them in the
// corresponding directory.  It does no upstream sync.
func (r *Repocket) SyncUnread() {
	favs := pocket.QueryUnread(r.AccessToken, r.ConsumerKey)
	if len(r.UnreadDir) == 0 {
		log.Fatalf("No output directory provided")
	}
	util.EnsureDir(r.UnreadDir)
	for _, item := range favs {
		DumpArticle(r.UnreadDir, &item)
	}
}
