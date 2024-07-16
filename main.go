package golastic

import (
	"context"
	"encoding/json"
	"errors"

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

func Delete(table string, id int) error {
	_, err := _client.Delete().
		Index(table).
		Do(context.Background())
	return err
}

func Update(table string, id int, data interface{}) error {
	_, err := _client.Update().
		Index(table).
		Id(string(id)).
		Doc(data).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		return err
	}
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

func Save(table string, id int, data interface{}) error {
	if _client == nil {
		return errors.New("Disconnected")
	}
	_, err := _client.Index().
		Index(table).
		Id(string(id)).
		BodyJson(data).
		Do(context.Background())
	return err
}
