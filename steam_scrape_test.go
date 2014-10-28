// Copyright 2014 Darren Schnare. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package dtoo

import (

  "os"
  "io"
  "net/http"
  "testing"
  "bufio"

  "strings"
  "regexp"
  "strconv"
  "github.com/PuerkitoBio/goquery"
)

type Game struct {
  Id string
  Name string
  DetailsUrl string
  LogoSmall string
  Logo string
  Metascore int
  ReleaseDate string
  Genres []string
}

func downloadSteamFixture() {
  var err error
  var file *os.File
  var resp *http.Response

  defer func () {
    if file != nil {
      file.Close()
    }
  }()

  if file,err = os.Open("./fixtures/steam.html"); err != nil {
    if file,err = os.Create("./fixtures/steam.html"); err == nil {
      if resp,err = http.Get("http://store.steampowered.com/search/"); err == nil {
        var line []byte
        buf := bufio.NewReader(resp.Body)
        for err == nil {
          if line,err = buf.ReadBytes('\n'); err == nil {
            _,err = file.Write(line)
          }
        }
      }
    }
  }

  if err != nil && err != io.EOF {
    panic(err)
  }
}

func TestSteamScrape(t *testing.T) {
  downloadSteamFixture()

  genreRegexp := regexp.MustCompile(`\s*-\s+Released.+`)
  commaRegexp := regexp.MustCompile(`\s*,\s*`)
  idRegexp := regexp.MustCompile(`https?:\/\/store\.steampowered\.com\/app\/([^\/]+?)\/`)
  logoRegexp := regexp.MustCompile(`sm_\d+`)

  if file,err := os.Open("./fixtures/steam.html"); err == nil {
    if games,err := ScrapeFromReader(".search_result_row", Model{
      "Id": func (s *goquery.Selection) (interface{}, error) {
        if href,exists := s.Attr("href"); exists {
          matches := idRegexp.FindAllStringSubmatch(href, -1)
          if len(matches) == 1 {
            return matches[0][1], nil
          }
        }

        return "", nil
      },
      "Name": RetrieverModel{Sel: ".search_name h4", Method: "text"},
      "DetailsUrl": "href",
      "LogoSmall": RetrieverModel{Sel: ".search_capsule img", Attr: "src"},
      "Metascore": RetrieverModel{Sel: ".search_metascore", Method: "text"},
      "ReleaseDate": RetrieverModel{Sel: ".search_released", Method: "text"},
      "Genres": RetrieverModel{
        Sel: ".search_name p",
        Method: func (s *goquery.Selection) (interface{}, error) {
          text := s.Text()
          text = strings.TrimSpace(text)
          text = genreRegexp.ReplaceAllLiteralString(text, "")
          genres := commaRegexp.Split(text, -1)
          return genres, nil
        },
      },
      "Logo": func (s *goquery.Selection) (interface{}, error) {
        if src,exists := s.Find(".search_capsule img").Attr("src"); exists {
          return logoRegexp.ReplaceAllLiteralString(src, "184x69"), nil
        }

        return "", nil
      },
    }, file); err == nil {
      if len(games) != 25 {
        t.Fatalf("game count invalid: expected %v got %v", 25, len(games))
      }

      for i,data := range(games) {
        if obj,ok := data.(Model); ok {
          game := toGame(obj)

          switch i {
            case 0:
              testGame(t, game, Game{
                Id: "570",
                Name: "Dota 2",
                DetailsUrl: "http://store.steampowered.com/app/570/?snr=1_7_7_230_150_1",
                LogoSmall: "http://cdn.akamai.steamstatic.com/steam/apps/570/capsule_sm_120.jpg?t=1404424435",
                Logo: "http://cdn.akamai.steamstatic.com/steam/apps/570/capsule_184x69.jpg?t=1404424435",
                Metascore: 90,
                ReleaseDate: "9 Jul 2013",
                Genres: []string{"Action", "Free to Play", "Strategy"},
              })
            case 4:
              testGame(t, game, Game{
                Id: "48700",
                Name: "Mount & Blade: Warband",
                DetailsUrl: "http://store.steampowered.com/app/48700/?snr=1_7_7_230_150_1",
                LogoSmall: "http://cdn.akamai.steamstatic.com/steam/apps/48700/capsule_sm_120.jpg?t=1405012491",
                Logo: "http://cdn.akamai.steamstatic.com/steam/apps/48700/capsule_184x69.jpg?t=1405012491",
                Metascore: 78,
                ReleaseDate: "31 Mar 2010",
                Genres: []string{"Action", "RPG"},
              })
            case 24:
              testGame(t, game, Game{
                Id: "49520",
                Name: "Borderlands 2",
                DetailsUrl: "http://store.steampowered.com/app/49520/?snr=1_7_7_230_150_1",
                LogoSmall: "http://cdn.akamai.steamstatic.com/steam/apps/49520/capsule_sm_120.jpg?t=1398979144",
                Logo: "http://cdn.akamai.steamstatic.com/steam/apps/49520/capsule_184x69.jpg?t=1398979144",
                Metascore: 89,
                ReleaseDate: "17 Sep 2012",
                Genres: []string{"Action", "RPG"},
              })
          }
        } else {
          t.Fatalf("invalid game encountered: %v", data)
        }
      }
    } else {
      panic(err)
    }

    if file != nil {
      file.Close()
    }
  } else {
    panic(err)
  }
}

func toGame(data Model) Game {
  return Game{
    Id: getStr(data, "Id", ""),
    Name: getStr(data, "Name", ""),
    DetailsUrl: getStr(data, "DetailsUrl", ""),
    LogoSmall: getStr(data, "LogoSmall", ""),
    Logo: getStr(data, "Logo", ""),
    Metascore: toInt(getStr(data, "Metascore", ""), 0),
    ReleaseDate: getStr(data, "ReleaseDate", ""),
    Genres: getStrSlice(data, "Genres"),
  }
}

func getStrSlice(data Model, key string) []string {
  if intrface,exists := data[key]; exists {
    if value,ok := intrface.([]string); ok {
      return value
    }
  }

  return make([]string, 0)
}

func toInt(s string, defaultValue int) int {
  if i,err := strconv.ParseInt(s, 10, 32); err == nil {
    return int(i)
  }
  return defaultValue
}

func testGame(t *testing.T, game Game, expected Game) {
  if game.Id != expected.Id {
    t.Fatalf("game Id invalid: expected %v got %v", expected.Id, game.Id)
  }

  if game.Name != expected.Name {
    t.Fatalf("game Name invalid: expected %v got %v", expected.Name, game.Name)
  }

  if game.DetailsUrl != expected.DetailsUrl {
    t.Fatalf("game DetailsUrl invalid: expected %v got %v", expected.DetailsUrl, game.DetailsUrl)
  }

  if game.LogoSmall != expected.LogoSmall {
    t.Fatalf("game LogoSmall invalid: expected %v got %v", expected.LogoSmall, game.LogoSmall)
  }

  if game.Logo != expected.Logo {
    t.Fatalf("game Logo invalid: expected %v got %v", expected.Logo, game.Logo)
  }

  if game.Metascore != expected.Metascore {
    t.Fatalf("game Metascore invalid: expected %v got %v", expected.Metascore, game.Metascore)
  }

  if game.ReleaseDate != expected.ReleaseDate {
    t.Fatalf("game ReleaseDate invalid: expected %v got %v", expected.ReleaseDate, game.ReleaseDate)
  }

  if strings.Join(game.Genres, ",") != strings.Join(expected.Genres, ",") {
    t.Fatalf("game Genres invalid: expected %v got %v", expected.Genres, game.Genres)
  }
}