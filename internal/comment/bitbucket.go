package comment

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Bitbucket doesn't support markdown comments
var bitbucketDefaultTag = "generated by Infracost"

// Bitbucket Cloud URL
var bitbucketDefaultServerURL = "https://bitbucket.org"

// bitbucketComment represents a comment on an Bitbucket pull request. It
// implements the Comment interface.
type bitbucketComment struct {
	id        int64
	body      string
	createdAt string
	url       string
}

// Body returns the body of the comment
func (c *bitbucketComment) Body() string {
	return c.body
}

// Ref returns the reference to the comment. For Bitbucket this is an API URL
// of the comment.
func (c *bitbucketComment) Ref() string {
	return c.url
}

// Less compares the comment to another comment and returns true if this
// comment should be sorted before the other comment.
func (c *bitbucketComment) Less(other Comment) bool {
	j := other.(*bitbucketComment)

	if c.createdAt != j.createdAt {
		return c.createdAt < j.createdAt
	}

	return c.id < j.id
}

// IsHidden always returns false for Bitbucket since Bitbucket doesn't have a
// feature for hiding comments.
func (c *bitbucketComment) IsHidden() bool {
	return false
}

// BitbucketExtra contains any extra inputs that can be passed to the Bitbucket
// comment handlers.
type BitbucketExtra struct {
	// ServerURL is the URL of the Bitbucket server. This can be set to a custom URL if
	// using Bitbucket Server/Data Center. If not set, the default Bitbucket server URL will be used.
	ServerURL string
	// OmitDetails is used to specify a format that excludes details output.
	OmitDetails bool
	// Tag is used to identify the Infracost comment.
	Tag string
	// Token is the Bitbucket access token.
	Token string
}

// bitbucketAPIComment represents API response structure of Azure Repos comment.
type bitbucketAPIComment struct {
	ID        int64  `json:"id"`
	CreatedOn string `json:"created_on"`
	IsDeleted bool   `json:"deleted"`
	Content   struct {
		Raw string `json:"raw"`
	} `json:"content"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
}

// newBitbucketAPIClient creates a HTTP client.
func newBitbucketAPIClient(ctx context.Context, token string) (*http.Client, error) {
	accessToken, tokenType := token, "Bearer"

	if strings.Contains(token, ":") {
		accessToken = base64.StdEncoding.EncodeToString([]byte(accessToken))
		tokenType = "Basic"
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken,
			TokenType:   tokenType,
		},
	)
	httpClient := oauth2.NewClient(ctx, ts)

	return httpClient, nil
}

// bitbucketPRHandler is a PlatformHandler for Bitbucket pull requests. It
// implements the PlatformHandler interface and contains the functions
// for finding, creating, updating, deleting comments on Bitbucket pull requests.
type bitbucketPRHandler struct {
	httpClient *http.Client
	apiURL     string
	prNumber   int
}

// NewBitbucketPRHandler creates a new PlatformHandler for Bitbucket pull requests.
func NewBitbucketPRHandler(ctx context.Context, repo string, targetRef string, extra BitbucketExtra) (*CommentHandler, error) {
	prNumber, err := strconv.Atoi(targetRef)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing targetRef as pull request number")
	}

	httpClient, err := newBitbucketAPIClient(ctx, extra.Token)
	if err != nil {
		return nil, err
	}

	var apiURL string
	var h PlatformHandler

	if strings.EqualFold(extra.ServerURL, bitbucketDefaultServerURL) {
		apiURL = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/", repo)

		h = &bitbucketPRHandler{
			httpClient: httpClient,
			apiURL:     apiURL,
			prNumber:   prNumber,
		}
	} else {
		serverRepo := strings.Split(repo, "/")

		if !strings.HasSuffix(extra.ServerURL, "/") {
			extra.ServerURL += "/"
		}
		apiURL = fmt.Sprintf("%srest/api/1.0/projects/%s/repos/%s/", extra.ServerURL, serverRepo[0], serverRepo[1])

		h = &bitbucketServerPRHandler{
			httpClient: httpClient,
			apiURL:     apiURL,
			prNumber:   prNumber,
		}
	}

	tag := bitbucketDefaultTag
	if extra.Tag != "" {
		tag = fmt.Sprintf("%s, tag: %s", tag, extra.Tag)
	}

	return NewCommentHandler(ctx, h, tag), nil
}

// CallFindMatchingComments calls the Bitbucket API to find the pull request
// comments that match the given tag, which has been embedded in the comment.
func (h *bitbucketPRHandler) CallFindMatchingComments(ctx context.Context, tag string) ([]Comment, error) {
	url := fmt.Sprintf("%spullrequests/%d/comments", h.apiURL, h.prNumber)

	matchingComments := []Comment{}

	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error getting comments")
		}

		res, err := h.httpClient.Do(req)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error getting comments")
		}

		if res.StatusCode != http.StatusOK {
			return []Comment{}, errors.Errorf("Error getting comments: %s", res.Status)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error reading response body")
		}

		var resData = struct {
			Comments    []bitbucketAPIComment `json:"values"`
			NextPageURL string                `json:"next"`
		}{}

		err = json.Unmarshal(resBody, &resData)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error unmarshaling response body")
		}

		for _, c := range resData.Comments {
			if c.IsDeleted || !strings.Contains(c.Content.Raw, tag) {
				continue
			}

			matchingComments = append(matchingComments, &bitbucketComment{
				id:        c.ID,
				body:      c.Content.Raw,
				createdAt: c.CreatedOn,
				url:       c.Links.Self.Href,
			})
		}

		if resData.NextPageURL != "" {
			url = resData.NextPageURL
		} else {
			break
		}
	}

	return matchingComments, nil
}

// CallCreateComment calls the Bitbucket API to create a new comment on the pull request.
func (h *bitbucketPRHandler) CallCreateComment(ctx context.Context, body string) (Comment, error) {
	reqData, err := json.Marshal(map[string]interface{}{
		"content": map[string]string{
			"raw": body,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error marshaling comment body")
	}

	url := fmt.Sprintf("%spullrequests/%d/comments", h.apiURL, h.prNumber)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, errors.Wrap(err, "Error creating request")
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating comment")
	}

	if res.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("Error creating comment: %s", res.Status)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	resData := bitbucketAPIComment{}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling response body")
	}

	return &bitbucketComment{
		id:        resData.ID,
		body:      resData.Content.Raw,
		createdAt: resData.CreatedOn,
		url:       resData.Links.HTML.Href,
	}, nil
}

// CallUpdateComment calls the Bitbucket API to update the body of a comment on the pull request.
func (h *bitbucketPRHandler) CallUpdateComment(ctx context.Context, comment Comment, body string) error {
	reqData, err := json.Marshal(map[string]interface{}{
		"content": map[string]string{
			"raw": body,
		},
	})
	if err != nil {
		return errors.Wrap(err, "Error marshaling comment body")
	}

	req, err := http.NewRequest("PUT", comment.Ref(), bytes.NewBuffer(reqData))
	if err != nil {
		return errors.Wrap(err, "Error creating request")
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := h.httpClient.Do(req)

	if res.Body != nil {
		defer res.Body.Close()
	}

	return err
}

// CallDeleteComment calls the Bitbucket API to delete the pull request comment.
func (h *bitbucketPRHandler) CallDeleteComment(ctx context.Context, comment Comment) error {
	req, err := http.NewRequest("DELETE", comment.Ref(), nil)
	if err != nil {
		return errors.Wrap(err, "Error creating request")
	}

	res, err := h.httpClient.Do(req)

	if res.Body != nil {
		defer res.Body.Close()
	}

	return err
}

// CallHideComment calls the Bitbucket API to minimize the pull request comment.
func (h *bitbucketPRHandler) CallHideComment(ctx context.Context, comment Comment) error {
	return errors.New("Not implemented")
}

// AddMarkdownTag appends a tag to the end of the given string. Bitbucket
// doesn't support markdown comments.
func (h *bitbucketPRHandler) AddMarkdownTag(s string, tag string) string {
	return fmt.Sprintf("%s%s", s, bitbucketMarkdownTag(tag))
}

// bitbucketCommitHandler is a PlatformHandler for Bitbucket commits. It
// implements the PlatformHandler interface and contains the functions
// for finding, creating, updating, deleting comments on Bitbucket commits.
type bitbucketCommitHandler struct {
	httpClient *http.Client
	apiURL     string
	commitSHA  string
}

// NewBitbucketCommitHandler creates a new PlatformHandler for Bitbucket commits.
func NewBitbucketCommitHandler(ctx context.Context, repo string, targetRef string, extra BitbucketExtra) (*CommentHandler, error) {
	httpClient, err := newBitbucketAPIClient(ctx, extra.Token)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(extra.ServerURL, bitbucketDefaultServerURL) {
		return nil, errors.New("Posting comments on commits is not available for Bitbucket Server")
	}

	apiURL := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/", repo)

	h := &bitbucketCommitHandler{
		httpClient: httpClient,
		apiURL:     apiURL,
		commitSHA:  targetRef,
	}

	tag := bitbucketDefaultTag
	if extra.Tag != "" {
		tag = fmt.Sprintf("%s, tag: %s", tag, extra.Tag)
	}

	return NewCommentHandler(ctx, h, tag), nil
}

// CallFindMatchingComments calls the Bitbucket API to find the pull request
// comments that match the given tag, which has been embedded in the comment.
func (h *bitbucketCommitHandler) CallFindMatchingComments(ctx context.Context, tag string) ([]Comment, error) {
	url := fmt.Sprintf("%scommit/%s/comments", h.apiURL, h.commitSHA)

	matchingComments := []Comment{}

	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error getting comments")
		}

		res, err := h.httpClient.Do(req)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error getting comments")
		}

		if res.StatusCode != http.StatusOK {
			return []Comment{}, errors.Errorf("Error getting comments: %s", res.Status)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error reading response body")
		}

		var resData = struct {
			Comments    []bitbucketAPIComment `json:"values"`
			NextPageURL string                `json:"next"`
		}{}

		err = json.Unmarshal(resBody, &resData)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error unmarshaling response body")
		}

		for _, c := range resData.Comments {
			if c.IsDeleted || !strings.Contains(c.Content.Raw, bitbucketMarkdownTag(tag)) {
				continue
			}

			matchingComments = append(matchingComments, &bitbucketComment{
				id:        c.ID,
				body:      c.Content.Raw,
				createdAt: c.CreatedOn,
				url:       c.Links.Self.Href,
			})
		}

		if resData.NextPageURL != "" {
			url = resData.NextPageURL
		} else {
			break
		}
	}

	return matchingComments, nil
}

// CallCreateComment calls the Bitbucket API to create a new comment on the pull request.
func (h *bitbucketCommitHandler) CallCreateComment(ctx context.Context, body string) (Comment, error) {
	reqData, err := json.Marshal(map[string]interface{}{
		"content": map[string]string{
			"raw": body,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error marshaling comment body")
	}

	url := fmt.Sprintf("%scommit/%s/comments", h.apiURL, h.commitSHA)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, errors.Wrap(err, "Error creating request")
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating comment")
	}

	if res.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("Error creating comment: %s", res.Status)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	resData := bitbucketAPIComment{}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling response body")
	}

	return &bitbucketComment{
		id:        resData.ID,
		body:      resData.Content.Raw,
		createdAt: resData.CreatedOn,
		url:       resData.Links.HTML.Href,
	}, nil
}

// CallUpdateComment calls the Bitbucket API to update the body of a comment on the pull request.
func (h *bitbucketCommitHandler) CallUpdateComment(ctx context.Context, comment Comment, body string) error {
	reqData, err := json.Marshal(map[string]interface{}{
		"content": map[string]string{
			"raw": body,
		},
	})
	if err != nil {
		return errors.Wrap(err, "Error marshaling comment body")
	}

	req, err := http.NewRequest("PUT", comment.Ref(), bytes.NewBuffer(reqData))
	if err != nil {
		return errors.Wrap(err, "Error creating request")
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := h.httpClient.Do(req)

	if res.Body != nil {
		defer res.Body.Close()
	}

	return err
}

// CallDeleteComment calls the Bitbucket API to delete the pull request comment.
func (h *bitbucketCommitHandler) CallDeleteComment(ctx context.Context, comment Comment) error {
	req, err := http.NewRequest("DELETE", comment.Ref(), nil)
	if err != nil {
		return errors.Wrap(err, "Error creating request")
	}

	res, err := h.httpClient.Do(req)

	if res.Body != nil {
		defer res.Body.Close()
	}

	return err
}

// CallHideComment calls the Bitbucket API to minimize the pull request comment.
func (h *bitbucketCommitHandler) CallHideComment(ctx context.Context, comment Comment) error {
	return errors.New("Not implemented")
}

// AddMarkdownTag appends a tag to the end of the given string. Bitbucket
// doesn't support markdown comments.
func (h *bitbucketCommitHandler) AddMarkdownTag(s string, tag string) string {
	return fmt.Sprintf("%s%s", s, bitbucketMarkdownTag(tag))
}

// bitbucketServerAPIComment represents Bitbucket Server API comment.
type bitbucketServerAPIComment struct {
	ID          int64  `json:"id"`
	CreatedDate int64  `json:"createdDate"`
	Text        string `json:"text"`
	Version     int64  `json:"version"`
}

// URL returns comment's Bitbucket Server API URL.
func (c *bitbucketServerAPIComment) URL(apiURL string, prNumber int) string {
	return fmt.Sprintf("%spull-requests/%d/comments/%d", apiURL, prNumber, c.ID)
}

// bitbucketServerAPIActivity represents Bitbucket Server API activity on
// a pull request.
type bitbucketServerAPIActivity struct {
	Action        string                    `json:"action"`
	CommentAction string                    `json:"commentAction"`
	Comment       bitbucketServerAPIComment `json:"comment"`
	CommentAnchor *struct{}                 `json:"commentAnchor"`
}

// Match checks if activity's comment matches the provided tag.
func (a *bitbucketServerAPIActivity) Match(tag string) bool {
	return a.Action == "COMMENTED" &&
		a.CommentAction == "ADDED" &&
		a.CommentAnchor == nil &&
		strings.Contains(a.Comment.Text, bitbucketMarkdownTag(tag))
}

// bitbucketServerPRHandler is a PlatformHandler for Bitbucket Server pull requests.
// It implements the PlatformHandler interface and contains the functions
// for finding, creating, updating, deleting comments on Bitbucket pull requests.
type bitbucketServerPRHandler struct {
	httpClient *http.Client
	apiURL     string
	prNumber   int
}

// CallFindMatchingComments calls the Bitbucket Server API to find the pull request
// comments that match the given tag, which has been embedded in the comment.
func (h *bitbucketServerPRHandler) CallFindMatchingComments(ctx context.Context, tag string) ([]Comment, error) {
	matchingComments := []Comment{}

	start := 0
	limit := 100

	for {
		url := fmt.Sprintf("%spull-requests/%d/activities?start=%d&limit=%d", h.apiURL, h.prNumber, start, limit)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error getting comments")
		}

		res, err := h.httpClient.Do(req)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error getting comments")
		}

		if res.StatusCode != http.StatusOK {
			return []Comment{}, errors.Errorf("Error getting comments: %s", res.Status)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error reading response body")
		}

		var resData = struct {
			Activities    []bitbucketServerAPIActivity `json:"values"`
			IsLastPage    bool                         `json:"isLastPage"`
			NextPageStart int                          `json:"nextPageStart"`
		}{}

		err = json.Unmarshal(resBody, &resData)
		if err != nil {
			return []Comment{}, errors.Wrap(err, "Error unmarshaling response body")
		}

		for _, activity := range resData.Activities {
			if !activity.Match(tag) {
				continue
			}

			comment := activity.Comment

			matchingComments = append(matchingComments, &bitbucketComment{
				id:        comment.ID,
				body:      comment.Text,
				createdAt: fmt.Sprint(comment.CreatedDate),
				url:       comment.URL(h.apiURL, h.prNumber),
			})
		}

		if resData.IsLastPage {
			break
		}

		start = resData.NextPageStart
	}

	return matchingComments, nil
}

// CallCreateComment calls the Bitbucket Server API to create a new comment on the pull request.
func (h *bitbucketServerPRHandler) CallCreateComment(ctx context.Context, body string) (Comment, error) {
	reqData, err := json.Marshal(map[string]string{
		"text": body,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error marshaling comment body")
	}

	url := fmt.Sprintf("%spull-requests/%d/comments", h.apiURL, h.prNumber)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, errors.Wrap(err, "Error creating request")
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating comment")
	}

	if res.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("Error creating comment: %s", res.Status)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	resData := bitbucketServerAPIComment{}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling response body")
	}

	return &bitbucketComment{
		id:        resData.ID,
		body:      resData.Text,
		createdAt: fmt.Sprint(resData.CreatedDate),
		url:       resData.URL(h.apiURL, h.prNumber),
	}, nil
}

// CallUpdateComment calls the Bitbucket Server API to update the body of a comment on the pull request.
func (h *bitbucketServerPRHandler) CallUpdateComment(ctx context.Context, comment Comment, body string) error {
	c, err := h.fetchServerComment(comment.Ref())
	if err != nil {
		return errors.Wrap(err, "Error retrieving comment version")
	}

	reqData, err := json.Marshal(map[string]string{
		"text":    body,
		"version": fmt.Sprint(c.Version),
	})
	if err != nil {
		return errors.Wrap(err, "Error marshaling comment body")
	}

	req, err := http.NewRequest("PUT", comment.Ref(), bytes.NewBuffer(reqData))
	if err != nil {
		return errors.Wrap(err, "Error creating request")
	}
	req.Header.Set("Content-Type", "application/json")

	res, _ := h.httpClient.Do(req)

	if res.Body != nil {
		defer res.Body.Close()
	}

	return err
}

// CallDeleteComment calls the Bitbucket Server API to delete the pull request comment.
func (h *bitbucketServerPRHandler) CallDeleteComment(ctx context.Context, comment Comment) error {
	c, err := h.fetchServerComment(comment.Ref())
	if err != nil {
		return errors.Wrap(err, "Error retrieving comment version")
	}

	url := fmt.Sprintf("%s?version=%d", comment.Ref(), c.Version)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "Error creating request")
	}

	res, err := h.httpClient.Do(req)

	if res.Body != nil {
		defer res.Body.Close()
	}

	return err
}

// CallHideComment calls the Bitbucket Server API to minimize the pull request comment.
func (h *bitbucketServerPRHandler) CallHideComment(ctx context.Context, comment Comment) error {
	return errors.New("Not implemented")
}

// AddMarkdownTag appends a tag to the end of the given string. Bitbucket
// doesn't support markdown comments.
func (h *bitbucketServerPRHandler) AddMarkdownTag(s string, tag string) string {
	return fmt.Sprintf("%s%s", s, bitbucketMarkdownTag(tag))
}

// fetchComment calls the Bitbucket Server API to retrieve a single comment.
func (h *bitbucketServerPRHandler) fetchServerComment(commentURL string) (*bitbucketServerAPIComment, error) {
	req, err := http.NewRequest("GET", commentURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating request")
	}

	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting comment")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Error getting comment: %s", res.Status)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	resData := bitbucketServerAPIComment{}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling response body")
	}

	return &resData, nil
}

func bitbucketMarkdownTag(tag string) string {
	return fmt.Sprintf("\n\n*(%s)*", tag)
}
