package golastic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/olivere/elastic/v7"
)

type GlTabler interface {
	TableName() string
}

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

func Update(data interface{}) error {
	table, body, err := getObjValue(data)
	if err != nil {
		return err
	}
	_, err = _client.Update().
		Index(table).
		Id(fmt.Sprintf("%v", body["id"])).
		Doc(body).
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

func Save(data interface{}) error {
	if _client == nil {
		return errors.New("Disconnected")
	}
	table, body, err := getObjValue(data)
	if err != nil {
		return err
	}
	_, err = _client.Index().
		Index(table).
		Id(fmt.Sprintf("%v", body["id"])).
		BodyJson(body).
		Do(context.Background())
	return err
}

func getObjValue(ptr interface{}) (string, map[string]interface{}, error) {

	tabler, ok := ptr.(GlTabler)
	if !ok {
		return "", nil, errors.New("must pointer to golastic.GlTabler")
	}

	v := reflect.ValueOf(ptr)

	res := map[string]interface{}{}

	val := v.Elem()

	if val.Kind() != reflect.Struct {
		return "", nil, errors.New("invalid object")
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("gl")
		if tag != "" && tag != "-" {
			res[tag] = val.FieldByName(typ.Field(i).Name).Interface()
		}
	}

	if res["id"] == nil || res["id"] == "" {
		return "", nil, errors.New("column id not found")
	}

	return (tabler).TableName(), res, nil
}
