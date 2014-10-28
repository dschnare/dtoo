// Copyright 2014 Darren Schnare. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

/*
Package dtoo exposes an HTML scraper API inspired by the artoo.js Scrape API  https://medialab.github.io/artoo/scrape/.

The dtoo scrape API closely follows artoo's scrape API, but is slightly modified to suit Go.
The biggest changes are that "scrapeTable" is not implemented and "scrapeOne" is implemented via
the ScrapeFromXxxWithLimit functions.

The artoo example

    artoo.scrape('li', {id: 'id', content: 'text'});

written for dtoo would look like this.

    dtoo.ScrapeFromUrl("li", dtoo.Model{"id": "id", "content": "text"}, url)

The above example will return a slice of dtoo.Model objects each with the following keys: {id, content}.

The dtoo data model passed to the ScrapeXxx and ScrapeXxxWithLimit functions 
can be a string, func (s *goquery.Selection) (interface{}, error), dtoo.Model or dtoo.RetrieverModel.

Retrieves a slice of post id attributes using a string data model.

    dtoo.ScrapeFromUrl(".post", "id", url)

Retrieves a slice of post comments using a function data model. When using a function you
are exposed to the low level goquery API for scraping out content from the DOM.

    type Comment struct {
      Author sring
      Content string
    }

    dtoo.ScrapeFromUrl(".post", func (s *goquery.Selection) (interface{}, error) {
      var err error = nil
      comments := make([]Comment, 0)

      s.Find(".comments").EachWithBreak(func (i int, s *goquery.Selection) bool {
        comment := Comment{
          Author: s.Find('.comment-author').Text(),
        }

        if content,e := s.Find('.comment-content').Html(); e == nil {
          comment.Content = content
        } else {
          // Save error and exit
          err = e
          return false
        }

        comments = append(comments, comment)
        return true
      })

      return comments, err
    }, url)

Retrieves a slice of dtoo.Model objects each with the following keys {Id, Title, PublishDate}.

    dtoo.ScrapeFromUrl(".post", dtoo.Model{
      "Id": "id",
      "Title": "title",
      "PublishDate": "datetime"
    }, url)

Using dtoo.Model you can nest dtoo.RetrieverModel objects and functions as keys in your model to
retrieve infinitely complex models. This example retrieves a slice of dtoo.Model objects with the
following properties {Id, Title, PublishedDate, Comments: [{Id, Author, Content}]}.

    dtoo.ScrapeFromUrl(".post", dtoo.Model{
      "Id": "id",
      "Title": dtoo.RetrieverModel{Sel: ".post-title", Method: "text"},
      "PublishDate": dtoo.RetrieverModel{Sel: ".publisehed-date", Attr: "datetime"},
      "Comments": func (s *goquery.Selection) (interface{}, error) {
        // We make a recursive call to dtoo.Scrape to make things easy.
        return dtoo.Scrape(".comment", dtoo.Model{
          "Id": "id",
          "Author": dtoo.RetrieverModel{Sel: ".comment-author", Method: "text"},
          "Content": dtoo.RetrieverModel{Sel: ".comment-content", Method: "html"},
        }, s)
      },
    }, url)

The above example could be alternatively achieved by using the recursive Scrape setting of RetrieverModel.

    dtoo.ScrapeFromUrl(".post", dtoo.Model{
      "Id": "id",
      "Title": dtoo.RetrieverModel{Sel: ".post-title", Method: "text"},
      "PublishDate": dtoo.RetrieverModel{Sel: ".publisehed-date", Attr: "datetime"},
      "Comments": dtoo.RetrieverModel{
        Scrape: dtoo.ScrapeObject{
          Iterator: ".comment", 
          Data: dtoo.Model{
            "Id": "id",
            "Author": dtoo.RetrieverObject{Sel: ".comment-author", Method: "text"},
            "Content": dtoo.RetrieverObject{Sel: ".comment-content", Method: "html"},
          },
        },
      },
    }, url)
*/

package dtoo