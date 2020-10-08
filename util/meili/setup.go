package meili

import (
	"log"

	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/viper"
)

// Client client for meili search server
var Client *meilisearch.Client

// SetupMeiliSearch setups the meili search server index
func SetupMeiliSearch() {
	Client = meilisearch.NewClient(meilisearch.Config{
		Host:   viper.GetString("meili.url"),
		APIKey: viper.GetString("meili.key"),
	})

	_, err := Client.Indexes().Get("data-portal")
	if err != nil {
		_, err = Client.Indexes().Create(meilisearch.CreateIndexRequest{
			UID:        "data-portal",
			PrimaryKey: "object_id",
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = Client.Settings("data-portal").UpdateAttributesForFaceting([]string{"kind"})
	if err != nil {
		log.Fatal(err)
	}

	// Add searchable attributes in index
	searchableAttributes := []string{"name", "slug", "description", "title", "contact_name", "contact_email", "license", "caption", "alt_text", "plan_name", "plan_info"}
	_, err = Client.Settings("data-portal").UpdateSearchableAttributes(searchableAttributes)
	if err != nil {
		log.Fatal(err)
	}

}
