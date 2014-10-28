// Copyright 2014 Darren Schnare. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package dtoo

import (
  "os"
  "testing"
)

type Post struct {
  Title string
  Date string
  Summary string
}

func TestScrapeText(t *testing.T) {
  if file,err := os.Open("./fixtures/index.html"); err == nil {
    if posts,err := ScrapeFromReader("time", "text", file); err == nil {
      if len(posts) != 6 {
        t.Fatalf("post count invalid: expected %v got %v", 6, len(posts))
      }

      for _,data := range(posts) {
        if time,ok := data.(string); ok {
          if len(time) == 0 {
            t.Fatalf("invalid time retrieved: expected a non-empty string")
          }
        } else {
          t.Fatalf("Invalid time retrieved: expected a string")
        }
      }
    }

    if file != nil {
      file.Close()
    }
  }
}

func TestScrapeModel(t *testing.T) {
  if file,err := os.Open("./fixtures/index.html"); err == nil {
    if posts,err := ScrapeFromReader(".row .col", Model{
      "Date": RetrieverModel{Sel: "time", Method: "text"},
      "Title": RetrieverModel{Sel: ".post-title", Method: "text"},
      "Summary": RetrieverModel{Sel: ".post-summary", Method: "text"},
    }, file); err == nil {
      if len(posts) != 6 {
        t.Fatalf("post count invalid: expected %v got %v", 6, len(posts))
      }

      for _,data := range(posts) {
        if obj,ok := data.(Model); ok {
          post := toPost(obj)

          if len(post.Title) == 0 {
            t.Fatalf("invalid title encountered: expected a non-empty string")
          }
          if len(post.Date) == 0 {
            t.Fatalf("invalid date encountered: expected a non-empty string")
          }
          if len(post.Summary) == 0 {
            t.Fatalf("invalid summary encountered: expected a non-empty string")
          }
        } else {
          t.Fatalf("invalid post encountered: %v", data)
        }
      }
    }

    if file != nil {
      file.Close()
    }
  }
}

func TestScrapeModelWithScrapeRetriever(t *testing.T) {
  if file,err := os.Open("./fixtures/index.html"); err == nil {
    if rows,err := ScrapeFromReader(".row", RetrieverModel{
      Scrape: ScrapeObject{
        Iterator: ".post",
        Data: Model{
          "Date": RetrieverModel{Sel: "time", Method: "text"},
          "Title": RetrieverModel{Sel: ".post-title", Method: "text"},
          "Summary": RetrieverModel{Sel: ".post-summary", Method: "text"},
        },
      },
    }, file); err == nil {
      if len(rows) != 3 {
        t.Fatalf("row count invalid: expected %v got %v", 3, len(rows))
      }

      for _,data := range(rows) {
        if row,ok := data.([]interface{}); ok {

          if len(row) != 2 {
            t.Fatalf("post count invalid: expected %v got %v", 2, len(row))
          }

          for _,intrface := range(row) {
            if obj,ok := intrface.(Model); ok {
              post := toPost(obj)

              if len(post.Title) == 0 {
                t.Fatalf("invalid title encountered: expected a non-empty string")
              }
              if len(post.Date) == 0 {
                t.Fatalf("invalid date encountered: expected a non-empty string")
              }
              if len(post.Summary) == 0 {
                t.Fatalf("invalid summary encountered: expected a non-empty string")
              }
            } else {
              t.Fatalf("invalid post encountered: %v", obj)
            }
          }
        } else {
          t.Fatalf("invalid row encountered: %v", row)
        }
      }
    }

    if file != nil {
      file.Close()
    }
  }
}

func toPost(data Model) Post {
  return Post{
    Title: getStr(data, "Title", ""),
    Date: getStr(data, "Date", ""),
    Summary: getStr(data, "Summary", ""),
  }
}

func getStr(data Model, key string, defaultValue string) string {
  if intrface,exists := data[key]; exists {
    if value,ok := intrface.(string); ok {
      return value
    }
  }

  return defaultValue
}