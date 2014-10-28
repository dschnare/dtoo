PACKAGE DOCUMENTATION

package dtoo
    import "github.com/dschnare/dtoo"

    Package dtoo exposes an HTML scraper API inspired by the artoo.js Scrape
    API https://medialab.github.io/artoo/scrape/.

    The dtoo scrape API closely follows artoo's scrape API, but is slightly
    modified to suit Go. The biggest changes are that "scrapeTable" is not
    implemented and "scrapeOne" is implemented via the
    ScrapeFromXxxWithLimit functions.

    The artoo example

	artoo.scrape('li', {id: 'id', content: 'text'});

    written for dtoo would look like this.

	dtoo.ScrapeFromUrl("li", dtoo.Model{"id": "id", "content": "text"}, url)

    The above example will return a slice of dtoo.Model objects each with
    the following keys: {id, content}.

    The dtoo data model passed to the ScrapeXxx and ScrapeXxxWithLimit
    functions can be a string, func (s *goquery.Selection) (interface{},
    error), dtoo.Model or dtoo.RetrieverModel.

    Retrieves a slice of post id attributes using a string data model.

	dtoo.ScrapeFromUrl(".post", "id", url)

    Retrieves a slice of post comments using a function data model. When
    using a function you are exposed to the low level goquery API for
    scraping out content from the DOM.

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

    Retrieves a slice of dtoo.Model objects each with the following keys
    {Id, Title, PublishDate}.

	dtoo.ScrapeFromUrl(".post", dtoo.Model{
	  "Id": "id",
	  "Title": "title",
	  "PublishDate": "datetime"
	}, url)

    Using dtoo.Model you can nest dtoo.RetrieverModel objects and functions
    as keys in your model to retrieve infinitely complex models. This
    example retrieves a slice of dtoo.Model objects with the following
    properties {Id, Title, PublishedDate, Comments: [{Id, Author,
    Content}]}.

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

    The above example could be alternatively achieved by using the recursive
    Scrape setting of RetrieverModel.

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

CONSTANTS

const (
    EMPTYSTRING = ""
)
    EMPTYSTRING is a convenient constant for an empty string.

FUNCTIONS

func Scrape(iterator string, model interface{}, s *goquery.Selection, limit uint) ([]interface{}, error)
    Scrape scrapes content from a goquery.Selection object according to the
    data model specified. Takes a selector as its root iterator and then
    takes the data model you intend to extract at each iteration. Will
    iterate up to limit number of iterations. If limit is 0 then no limit is
    applied. The data model can be a string, (s *goquery.Selection)
    (interface{}, error), Model or RetrieverModel.

    Returns a value based on the data model specified. See package examples
    for more info.

    Example:

	doc, err := goquery.NewDocument(url)
	if err == nil {
	  Scrape("li", dtoo.Model{id: 'id', content: 'text'}, doc.Selection, 0)
	}

func ScrapeFromReader(iterator string, model interface{}, r io.Reader) ([]interface{}, error)
    ScrapeFromReader scrapes content from a file according to the data model
    specified. Takes a selector as its root iterator and then takes the data
    model you intend to extract at each iteration. The data model can be a
    string, (s *goquery.Selection) (interface{}, error), Model or
    RetrieverModel.

    Returns a value based on the data model specified. See package examples
    for more info.

    Example:

	dtoo.ScrapeFromReader("li", dtoo.Model{id: 'id', content: 'text'}, reader)

func ScrapeFromReaderWithLimit(iterator string, model interface{}, r io.Reader, limit uint) ([]interface{}, error)
    ScrapeFromReaderWithLimit scrapes content from a file according to the
    data model specified up to a limit. Takes a selector as its root
    iterator and then takes the data model you intend to extract at each
    iteration. The data model can be a string, (s *goquery.Selection)
    (interface{}, error), Model or RetrieverModel. Will iterate up to limit
    number of iterations. If limit is 0 then no limit is applied.

    Returns a value based on the data model specified. See package examples
    for more info.

    Example:

	dtoo.ScrapeFromReaderWithLimit("li", dtoo.Model{id: 'id', content: 'text'}, reader, 0)

func ScrapeFromString(iterator string, model interface{}, html string) ([]interface{}, error)
    ScrapeFromString scrapes content from an HTML string according to the
    data model specified. Takes a selector as its root iterator and then
    takes the data model you intend to extract at each iteration. The data
    model can be a string, (s *goquery.Selection) (interface{}, error),
    Model or RetrieverModel.

    Returns a value based on the data model specified. See package examples
    for more info.

    Example:

	dtoo.ScrapeFromString("li", dtoo.Model{id: 'id', content: 'text'}, html)

func ScrapeFromStringWithLimit(iterator string, model interface{}, html string, limit uint) ([]interface{}, error)
    ScrapeFromStringWithLimit scrapes content from an HTML string according
    to the data model specified up to a limit. Takes a selector as its root
    iterator and then takes the data model you intend to extract at each
    iteration. The data model can be a string, (s *goquery.Selection)
    (interface{}, error), Model or RetrieverModel. Will iterate up to limit
    number of iterations. If limit is 0 then no limit is applied.

    Returns a value based on the data model specified. See package examples
    for more info.

    Example:

	dtoo.ScrapeFromStringWithLimit("li", dtoo.Model{id: 'id', content: 'text'}, html, 0)

func ScrapeFromUrl(iterator string, model interface{}, url string) ([]interface{}, error)
    ScrapeFromUrl scrapes content from a URL according to the data model
    specified. Takes a selector as its root iterator and then takes the data
    model you intend to extract at each iteration. The data model can be a
    string, (s *goquery.Selection) (interface{}, error), Model or
    RetrieverModel.

    Returns a value based on the data model specified. See package examples
    for more info.

    Example:

	dtoo.ScrapeFromUrl("li", dtoo.Model{id: 'id', content: 'text'}, url)

func ScrapeFromUrlWithLimit(iterator string, model interface{}, url string, limit uint) ([]interface{}, error)
    ScrapeFromUrlWithLimit scrapes content from a URL according to the data
    model specified up to a limit. Takes a selector as its root iterator and
    then takes the data model you intend to extract at each iteration. The
    data model can be a string, (s *goquery.Selection) (interface{}, error),
    Model or RetrieverModel. Will iterate up to limit number of iterations.
    If limit is 0 then no limit is applied.

    Returns a value based on the data model specified. See package examples
    for more info.

    Example:

	dtoo.ScrapeFromUrlWithLimit("li", dtoo.Model{id: 'id', content: 'text'}, url, 0)

TYPES

type Model map[string]interface{}
    Model is a complex model that is a map of string:interface{}.

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
    RetrieverModel is an expressive object that can perform subselection
    scraping. This data model type behaves in a similar fashion as the
    retriever object interface accepted by artoo.

    Examples:

    Retrieves a slice of post titles.

	dtoo.ScrapeFromUrl(".post", dtoo.RetrieverModel{Sel: ".post-title", Method: "text", DefaultValue: "Unknown Post"}, url)

    Retrieves a slice of post published dates via the datetime attribute.

	dtoo.ScrapeFromUrl(".post", dtoo.RetrieverModel{Sel: ".publisehed-date", Attr: "datetime"}, url)

    Retrieve a slice of slice of comment authors.

	dtoo.ScrapeFromUrl(".post", dtoo.RetrieverModel{Scrape: dtoo.ScrapeObject{Iterator: ".comment-author", Data: "text"}}, url)

type ScrapeObject struct {
    // The selector that will act as the root iterator for a recursive scrape.
    Iterator string
    // The data model for the recursive scrape. Can be any data model type accepted by the Scrape functions
    Data interface{}
}
    ScrapeObject is an object that specifies settings for recursive
    scraping.

SUBDIRECTORIES

	fixtures

