package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/blevesearch/bleve/mapping"
)

// NewSearch creates a new instance of bleve.Index without the boilerplate.
func NewSearch(name string, mapping *mapping.IndexMappingImpl) (bleve.Index, error) {
	path := filepath.Join(".", "_data", name+".bleve")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		_ = os.RemoveAll(path)
	}

	index, err := bleve.NewUsing(path, mapping, scorch.Name, scorch.Name, nil)
	if err != nil {
		return nil, err
	}

	return index, nil
}

// UpsertEntry inserts a new entry into a specific search index.
func UpsertEntry(index bleve.Index, id string, data interface{}) error {
	return index.Index(id, data)
}

// DeleteEntry removes an existing entry from the search index.
func DeleteEntry(index bleve.Index, id string) error {
	return index.Delete(id)
}

func injectSearchQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value("searchRequest").(*bleve.SearchRequest)
		if !ok {
			searchReq := bleve.NewSearchRequest(bleve.NewConjunctionQuery())
			ctx := context.WithValue(r.Context(), "searchRequest", searchReq)
			next.ServeHTTP(rw, r.WithContext(ctx))
		} else {
			next.ServeHTTP(rw, r)
		}
	})
}
