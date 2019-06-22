package astislack

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/http"
	"github.com/pkg/errors"
)

type Slack struct {
	legacyToken string
	s           *astihttp.Sender
}

func New(c Configuration) *Slack {
	return &Slack{
		legacyToken: c.LegacyToken,
		s:           astihttp.NewSender(c.Sender),
	}
}

func (s *Slack) send(method, url string, ps url.Values, payloadOut interface{}) (err error) {
	// Add token to parameters
	ps.Set("token", s.legacyToken)

	// Create request
	var req *http.Request
	if req, err = http.NewRequest(method, url+"?"+ps.Encode(), nil); err != nil {
		err = errors.Wrap(err, "astislack: creating request failed")
		return
	}

	// Send
	var resp *http.Response
	if resp, err = s.s.Send(req); err != nil {
		err = errors.Wrap(err, "astislack: sending failed")
		return
	}
	defer resp.Body.Close()

	// Unmarshal
	if payloadOut != nil {
		if err = json.NewDecoder(resp.Body).Decode(payloadOut); err != nil {
			err = errors.Wrap(err, "astislack: unmarshaling failed")
			return
		}
	}
	return
}

type Response struct {
	Files    ResponseFiles    `json:"files"`
	Messages ResponseMessages `json:"messages"`
	OK       bool             `json:"ok"`
	Team     string           `json:"team"`
	TeamID   string           `json:"team_id"`
	User     string           `json:"user"`
	UserID   string           `json:"user_id"`
}

type ResponseFiles struct {
	Matches    []ResponseFile     `json:"matches"`
	Pagination ResponsePagination `json:"pagination"`
	Paging     ResponsePaging     `json:"paging"`
	Total      int                `json:"total"`
}

type ResponseFile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ResponseMessages struct {
	Matches    []ResponseMessage  `json:"matches"`
	Pagination ResponsePagination `json:"pagination"`
	Paging     ResponsePaging     `json:"paging"`
	Total      int                `json:"total"`
}

type ResponseMessage struct {
	Channel ResponseChannel `json:"channel"`
	Text    string          `json:"text"`
	TS      string          `json:"ts"`
}

type ResponseChannel struct {
	ID   string `json:"id"`
	Name string `json:"string"`
}

type ResponsePagination struct {
	First      int `json:"first"`
	Last       int `json:"last"`
	Page       int `json:"page"`
	PageCount  int `json:"page_count"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
}

type ResponsePaging struct {
	Count int `json:"count"`
	Page  int `json:"page"`
	Pages int `json:"pages"`
	Total int `json:"total"`
}

func (s *Slack) Delete(query string) (err error) {
	// Loop
	var resp Response
	for (resp.Files.Pagination.Page == 0 || resp.Files.Pagination.Page <= resp.Files.Pagination.PageCount) || (resp.Messages.Pagination.Page == 0 || resp.Messages.Pagination.Page <= resp.Messages.Pagination.PageCount) {
		// Create parameters
		ps := url.Values{}
		ps.Set("query", query)

		// Add page
		if resp.Files.Pagination.Page > 0 {
			ps.Set("page", strconv.Itoa(resp.Files.Pagination.Page+1))
		}

		// Send
		if err = s.send(http.MethodGet, "https://slack.com/api/search.all", ps, &resp); err != nil {
			err = errors.Wrap(err, "astislack: sending failed")
			return
		}

		// Loop through files
		for _, f := range resp.Files.Matches {
			// Delete
			astilog.Infof("astislack: deleting file '%s'", f.Name)
			if err = s.deleteFile(f.ID); err != nil {
				err = errors.Wrap(err, "astislack: deleting file failed")
				return
			}
		}

		// Loop through messages
		for _, m := range resp.Messages.Matches {
			// Delete
			astilog.Infof("astislack: deleting message '%s' sent to '%s'", m.Text, m.Channel.Name)
			if err = s.deleteMessage(m.Channel.ID, m.TS); err != nil {
				err = errors.Wrap(err, "astislack: deleting message failed")
				return
			}
		}
	}
	return
}

func (s *Slack) deleteFile(id string) (err error) {
	// Create parameters
	ps := url.Values{}
	ps.Set("file", id)

	// Send
	var resp Response
	if err = s.send(http.MethodPost, "https://slack.com/api/files.delete", ps, &resp); err != nil {
		err = errors.Wrap(err, "astislack: sending failed")
		return
	}
	return
}

func (s *Slack) deleteMessage(channel, ts string) (err error) {
	// Create parameters
	ps := url.Values{}
	ps.Set("channel", channel)
	ps.Set("ts", ts)

	// Send
	var resp Response
	if err = s.send(http.MethodPost, "https://slack.com/api/chat.delete", ps, &resp); err != nil {
		err = errors.Wrap(err, "astislack: sending failed")
		return
	}
	return
}

type Me struct {
	Team   string
	TeamID string
	User   string
	UserID string
}

func (s *Slack) Me() (m Me, err error) {
	// Send
	var resp Response
	if err = s.send(http.MethodPost, "https://slack.com/api/auth.test", url.Values{}, &resp); err != nil {
		err = errors.Wrap(err, "astislack: sending failed")
		return
	}

	// Create me
	m = Me{
		Team:   resp.Team,
		TeamID: resp.TeamID,
		User:   resp.User,
		UserID: resp.UserID,
	}
	return
}
