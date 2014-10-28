// Copyright 2014 Darren Schnare. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package dtoo

import (
  "bytes"
  "io"
  "html"
  "errors"
  "github.com/PuerkitoBio/goquery"
)

// EMPTYSTRING is a convenient constant for an empty string.
const (
  EMPTYSTRING = ""
)

// Model is a complex model that is a map of string:interface{}.
type Model map[string]interface{}

/*
RetrieverModel is an expressive object that can perform subselection scraping. 
This data model type behaves in a similar fashion as the retriever object interface accepted by artoo.

Examples:

Retrieves a slice of post titles.

    dtoo.ScrapeFromUrl(".post", dtoo.RetrieverModel{Sel: ".post-title", Method: "text", DefaultValue: "Unknown Post"}, url) 

Retrieves a slice of post published dates via the datetime attribute.

    dtoo.ScrapeFromUrl(".post", dtoo.RetrieverModel{Sel: ".publisehed-date", Attr: "datetime"}, url) 

Retrieve a slice of slice of comment authors.
   
    dtoo.ScrapeFromUrl(".post", dtoo.RetrieverModel{Scrape: dtoo.ScrapeObject{Iterator: ".comment-author", Data: "text"}}, url)
*/
type RetrieverModel struct {
  // The name of the attribute to extract.
  Attr string
  // The CSS selector to extract from.
  Sel string
  // The method to use for the extraction. Can be set to "text", "html" or a func(*goquery.Selection)(interface{}, error). Required if Sel is set.
  Method interface{}
  // If set to a dtoo.ScrapeObject then a recursive scrape will be executed.
  Scrape ScrapeObject
  // If set and the retrieved value is nil or the empty string then returns this value.
  DefaultValue interface{}
}

// ScrapeObject is an object that specifies settings for recursive scraping.
type ScrapeObject struct {
  // The selector that will act as the root iterator for a recursive scrape.
  Iterator string
  // The data model for the recursive scrape. Can be any data model type accepted by the Scrape functions
  Data interface{}
}

// ScrapeFromStringWithLimit scrapes content from an HTML string according to the data model specified up to a limit.
// Takes a selector as its root iterator and then takes the data model you intend to extract at each iteration.
// The data model can be a string, (s *goquery.Selection) (interface{}, error), Model or RetrieverModel.
// Will iterate up to limit number of iterations. If limit is 0 then no limit is applied.
//
// Returns a value based on the data model specified. See package examples for more info.
//
// Example:
//
//    dtoo.ScrapeFromStringWithLimit("li", dtoo.Model{id: 'id', content: 'text'}, html, 0)
func ScrapeFromStringWithLimit(iterator string, model interface{}, html string, limit uint) ([]interface{}, error) {
  return ScrapeFromReaderWithLimit(iterator, model, bytes.NewBufferString(html), limit)
}

// ScrapeFromReaderWithLimit scrapes content from a file according to the data model specified up to a limit.
// Takes a selector as its root iterator and then takes the data model you intend to extract at each iteration.
// The data model can be a string, (s *goquery.Selection) (interface{}, error), Model or RetrieverModel.
// Will iterate up to limit number of iterations. If limit is 0 then no limit is applied.
//
// Returns a value based on the data model specified. See package examples for more info.
//
// Example:
// 
//    dtoo.ScrapeFromReaderWithLimit("li", dtoo.Model{id: 'id', content: 'text'}, reader, 0)
func ScrapeFromReaderWithLimit(iterator string, model interface{}, r io.Reader, limit uint) ([]interface{}, error) {
  doc, err := goquery.NewDocumentFromReader(r)

  if err == nil {
    return Scrape(iterator, model, doc.Selection, limit)
  } else {
    return nil, err
  }
}

// ScrapeFromUrlWithLimit scrapes content from a URL according to the data model specified up to a limit.
// Takes a selector as its root iterator and then takes the data model you intend to extract at each iteration.
// The data model can be a string, (s *goquery.Selection) (interface{}, error), Model or RetrieverModel.
// Will iterate up to limit number of iterations. If limit is 0 then no limit is applied.
//
// Returns a value based on the data model specified. See package examples for more info.
//
// Example:
//
//    dtoo.ScrapeFromUrlWithLimit("li", dtoo.Model{id: 'id', content: 'text'}, url, 0)
func ScrapeFromUrlWithLimit(iterator string, model interface{}, url string, limit uint) ([]interface{}, error) {
  doc, err := goquery.NewDocument(url)

  if err == nil {
    return Scrape(iterator, model, doc.Selection, limit)
  } else {
    return nil, err
  }
}

// ScrapeFromString scrapes content from an HTML string according to the data model specified.
// Takes a selector as its root iterator and then takes the data model you intend to extract at each iteration.
// The data model can be a string, (s *goquery.Selection) (interface{}, error), Model or RetrieverModel.
//
// Returns a value based on the data model specified. See package examples for more info.
//
// Example:
//
//    dtoo.ScrapeFromString("li", dtoo.Model{id: 'id', content: 'text'}, html)
func ScrapeFromString(iterator string, model interface{}, html string) ([]interface{}, error) {
  return ScrapeFromReaderWithLimit(iterator, model, bytes.NewBufferString(html), 0)
}

// ScrapeFromReader scrapes content from a file according to the data model specified.
// Takes a selector as its root iterator and then takes the data model you intend to extract at each iteration.
// The data model can be a string, (s *goquery.Selection) (interface{}, error), Model or RetrieverModel.
//
// Returns a value based on the data model specified. See package examples for more info.
//
// Example:
//
//    dtoo.ScrapeFromReader("li", dtoo.Model{id: 'id', content: 'text'}, reader)
func ScrapeFromReader(iterator string, model interface{}, r io.Reader) ([]interface{}, error) {
  return ScrapeFromReaderWithLimit(iterator, model, r, 0)
}

// ScrapeFromUrl scrapes content from a URL according to the data model specified.
// Takes a selector as its root iterator and then takes the data model you intend to extract at each iteration.
// The data model can be a string, (s *goquery.Selection) (interface{}, error), Model or RetrieverModel.
//
// Returns a value based on the data model specified. See package examples for more info.
//
// Example:
//
//    dtoo.ScrapeFromUrl("li", dtoo.Model{id: 'id', content: 'text'}, url)
func ScrapeFromUrl(iterator string, model interface{}, url string) ([]interface{}, error) {
  return ScrapeFromUrlWithLimit(iterator, model, url, 0)
}

// Scrape scrapes content from a goquery.Selection object according to the data model specified.
// Takes a selector as its root iterator and then takes the data model you intend to extract at each iteration.
// Will iterate up to limit number of iterations. If limit is 0 then no limit is applied.
// The data model can be a string, (s *goquery.Selection) (interface{}, error), Model or RetrieverModel.
//
// Returns a value based on the data model specified. See package examples for more info.
//
// Example:
//
//    doc, err := goquery.NewDocument(url)
//    if err == nil {
//      Scrape("li", dtoo.Model{id: 'id', content: 'text'}, doc.Selection, 0)
//    }
func Scrape(iterator string, model interface{}, s *goquery.Selection, limit uint) ([]interface{}, error) {
  var theError error = nil
  c := uint(0)
  result := make([]interface{}, 1)

  s.Find(iterator).EachWithBreak(func (i int, s *goquery.Selection) bool {
    data, err := extract(model, s)

    if err == nil {
      if c >= uint(len(result)) {
        result = append(result, data)
      } else {
        result[c] = data
      }
    } else {
      theError = err
    }

    c++
    
    // If we return false then the loop will be broken.
    return !(c == limit || theError != nil)
  })

  return result, theError
}

func extract(model interface{}, s *goquery.Selection) (interface{}, error) {
  switch modelValue := model.(type) {
    case string:
      return extractString(modelValue, s)
    case RetrieverModel:
      return extractRetrieverModel(modelValue, s)
    case Model:
      return extractDataModel(modelValue, s)
    case func (s *goquery.Selection) (interface{}, error):
      return modelValue(s)
    default:
      return nil, errors.New("Unsupported retriever type")
  }
}

func extractString(model string, s *goquery.Selection) (string, error) {
  switch model {
    // the text of the selection
    case "text":
      return html.UnescapeString(s.Text()), nil
    case "html":
      return s.Html()
    // attribute name
    default:
      if attrValue,hasAttr := s.Attr(model); hasAttr {
        return attrValue, nil
      } else {
        return EMPTYSTRING, nil
      }
  }
}

func extractRetrieverModel(rm RetrieverModel, s *goquery.Selection) (interface{}, error) {
  if rm.Sel != EMPTYSTRING {
    s = s.Find(rm.Sel)
  }

  if rm.Attr != EMPTYSTRING {
    if attrValue,hasAttr := s.Attr(rm.Attr); hasAttr {
      return attrValue, nil
    } else {
      return rm.DefaultValue, nil
    }
  } else if rm.Method != nil {
    switch method := rm.Method.(type) {
      case string:
        switch method {
          case "text":
            return html.UnescapeString(s.Text()), nil
          case "html":
            return s.Html()
        }
      case func (s *goquery.Selection) (interface{}, error):
        return method(s)
      default:
        return nil, errors.New("RetrieverModel: unrecognized 'method' type " + rm.Sel)
    }
  } else if rm.Scrape.Iterator != EMPTYSTRING && rm.Scrape.Data != EMPTYSTRING {
    return Scrape(rm.Scrape.Iterator, rm.Scrape.Data, s, 0)
  }

  return nil, errors.New("Empty RetrieverModel encountered")
}

func extractDataModel(model Model, s *goquery.Selection) (data Model, err error) {
  data = Model{}

  for key,retriever := range(model) {
    if data[key],err = extract(retriever, s); err != nil {
      break
    }
  }

  return
}