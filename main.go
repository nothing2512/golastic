package golastic

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
)

var _client *elastic.Client

func Connect(uri string) error {
	client, err := elastic.NewClient(elastic.SetURL(uri), elastic.SetSniff(false))
	if err != nil {
		return err
	}
	_client = client
	return nil
}

func Delete(table, key string, value interface{}) error {
	_, err := _client.DeleteByQuery(table).
		Query(elastic.NewTermQuery(key, value)).
		Do(context.Background())
	return err
}

func Search(obj interface{}, table, value string, keys ...string) {
	sr, _ := _client.Search().
		Index(table).
		Query(term(value, keys...)).
		Do(context.TODO())
	responses := []interface{}{}
	for _, hit := range sr.Hits.Hits {
		var src map[string]interface{}
		json.Unmarshal(hit.Source, &src)
		responses = append(responses, src)
	}
	b, _ := json.Marshal(responses)
	_ = json.Unmarshal(b, obj)
}

func term(value string, keys ...string) *elastic.BoolQuery {
	var queries []elastic.Query
	for _, v := range keys {
		q := elastic.NewFuzzyQuery(v, value).
			Fuzziness(2).
			Transpositions(true)
		queries = append(queries, q)
	}
	return elastic.NewBoolQuery().MinimumShouldMatch("1").Should(queries...)
}

func Save(table string, data interface{}) error {
	if _client == nil {
		return errors.New("Disconnected")
	}
	_, err := _client.Index().
		Index(table).
		Id(uuid()).
		BodyJson(data).
		Do(context.Background())
	return err
}

func uuid() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return ""
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	now := time.Now().UnixNano()
	uuid[0] = byte(now >> 54)
	uuid[1] = byte(now >> 48)
	uuid[2] = byte(now >> 42)
	uuid[3] = byte(now >> 36)
	uuid[4] = byte(now >> 30)
	uuid[5] = byte(now >> 24)
	uuid[6] = byte(now >> 18)
	uuid[7] = byte(now >> 12)
	uuid[8] = byte(now >> 6)

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
