package rtwitter

import (
	"bytes"
	"context"
	"fmt"
	"go-stock-prices/model"
	"go-stock-prices/ports"
	"net/http"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
	"github.com/pkg/errors"
)

const (
	perfTemplate = `{{.Symbol}} returned {{.SimpleReturn}}% with a maximum drawdown of {{.MaxDrawdown}}%
({{.From}} - {{.To}})`
)

var duplicateContentRegex = regexp.MustCompile(".*not.*allowed.*duplicate.*content.*")

type perfTemplateData struct {
	Symbol       string
	SimpleReturn string
	MaxDrawdown  string
	From         string
	To           string
	Now          string
}

func newPerfTemplateData(p *model.Performance) *perfTemplateData {
	return &perfTemplateData{
		Symbol:       strings.ToUpper(p.Symbol),
		SimpleReturn: fmt.Sprintf("%.2f", p.SimpleReturn*100),
		MaxDrawdown:  fmt.Sprintf("%.2f", p.MaxDrawdown*100),
		From:         p.From.String(),
		To:           p.To.String(),
		Now:          time.Now().Format(time.DateTime),
	}
}

type twitterNotificationRepository struct {
	client       *gotwi.Client
	perfTemplate *template.Template
}

func NewTwitterNotificationRepository(client *gotwi.Client) ports.NotificationRepository {
	perfTemplate := template.Must(template.New("performance").Parse(perfTemplate))

	return &twitterNotificationRepository{
		client:       client,
		perfTemplate: perfTemplate,
	}
}

func (r *twitterNotificationRepository) NotifyPerformance(ctx context.Context, performance *model.Performance) error {
	text, err := r.createPerformanceTweetText(performance)
	if err != nil {
		return nil
	}

	p := &types.CreateInput{
		Text: gotwi.String(text),
	}

	_, err = managetweet.Create(ctx, r.client, p)
	if err != nil {
		ge, ok := err.(*gotwi.GotwiError)
		if ok {
			if isDuplicateContentError(ge) {
				fmt.Println("trying to post duplicate content. ignoring tweet.")
				return nil
			}
		}

		return errors.Wrap(err, "failed to post tweet")
	}

	return nil
}

func isDuplicateContentError(gotwiErr *gotwi.GotwiError) bool {
	if gotwiErr.StatusCode != http.StatusForbidden {
		return false
	}

	return duplicateContentRegex.Match([]byte(gotwiErr.Detail))
}

func (r *twitterNotificationRepository) createPerformanceTweetText(performance *model.Performance) (string, error) {
	data := newPerfTemplateData(performance)

	var buf bytes.Buffer
	err := r.perfTemplate.Execute(&buf, data)
	if err != nil {
		return "", errors.Wrap(err, "failed to create tweet text")
	}

	return buf.String(), nil
}
