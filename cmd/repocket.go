package main

import (
	"log"
	"os"

	"github.com/srvaroa/repocket/pkg/repocket"
	"github.com/srvaroa/repocket/pkg/util"
)

func usageAndExit() {
	log.Fatal("Usage: %s [favs|delete|archive|unread|sync]\n", os.Args[0])
}

func main() {

	log.SetFlags(0)
	log.SetOutput(new(util.LogWriter))

	if len(os.Args) != 2 {
		usageAndExit()
	}

	cmd := os.Args[1]

	r := repocket.Repocket{}
	err := r.Load()
	if err != nil {
		log.Fatalf("Unable to load configuration!", err)
	}
	r.Authenticate()

	switch cmd {
	case "delete":
		r.SyncDeletions()
		break
	case "archive":
		r.SyncArchived()
		break
	case "favs":
		r.SyncFavs()
		break
	case "unread":
		r.SyncUnread()
		break
	case "sync":
		log.Printf("Full sync")
		log.Printf("First push deletions..")
		r.SyncDeletions()
		log.Printf("Then push archived..")
		r.SyncArchived()
		log.Printf("Then push favs..")
		r.SyncFavs()
		log.Printf("Then pull unreads..")
		r.SyncUnread()
		log.Printf("All done!")
		break
	default:
		usageAndExit()
	}

}
