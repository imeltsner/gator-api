package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/imeltsner/gator-api/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}

	client := http.Client{}
	req.Header.Add("User-Agent", "gator-api")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get response: %v", err)
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(content, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal xml %v", err)
	}
	rssFeed.unescape()

	return &rssFeed, nil
}

func (feed *RSSFeed) unescape() {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}
}

func (s *state) scrapeFeeds() error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get next feed: %v", err)
	}

	feedFetchedParams := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:     time.Now().UTC(),
		ID:            nextFeed.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), feedFetchedParams)
	if err != nil {
		return fmt.Errorf("unable to mark feed fetched: %v", err)
	}

	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("unable to fetch feed from %v: %v", nextFeed.Url, err)
	}
	fmt.Printf("Feed at url %v fetched successfully\n", rssFeed.Channel.Link)

	err = s.saveFeed(*rssFeed, nextFeed)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) saveFeed(feed RSSFeed, dbFeed database.Feed) error {
	for _, item := range feed.Channel.Item {
		postParams := generatePostParams(item, dbFeed)
		post, err := s.db.CreatePost(context.Background(), postParams)
		if err != nil && strings.Contains(err.Error(), "duplicate key value") {
			return nil
		} else if err != nil {
			return fmt.Errorf("unable to create post %v: %v", item.Title, err)
		}
		fmt.Printf("Successfully created post %v\n", post.Title)
	}

	return nil
}

func generatePostParams(item RSSItem, feed database.Feed) database.CreatePostParams {
	var description sql.NullString
	if item.Description == "" {
		description = sql.NullString{}
	} else {
		description = sql.NullString{String: item.Description, Valid: true}
	}

	var pubDate sql.NullTime
	if item.PubDate != "" {
		parsedPubdate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			pubDate = sql.NullTime{}
		}
		pubDate = sql.NullTime{Time: parsedPubdate, Valid: true}
	} else {
		pubDate = sql.NullTime{}
	}

	return database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Title:       item.Title,
		Url:         item.Link,
		Description: description,
		PublishedAt: pubDate,
		FeedID:      feed.ID,
	}
}
