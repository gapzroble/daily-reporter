package tempo

import (
	"context"
	"errors"
	"net/url"
	U "net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"

	"github.com/rroble/daily-reporter/lib/log"
	"github.com/rroble/daily-reporter/lib/schedule"
)

// Report func
func Report(doneTempo <-chan bool) ([]byte, error) {
	if !<-doneTempo {
		return nil, errors.New("Tempo not logged")
	}

	log.Debug("jira", "Creating screenshot")
	defer log.Debug("jira", "Done screenshot")
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer timeoutCancel()

	opts := []chromedp.ExecAllocatorOption{
		// show it to load tempo iframe
		// just hide it in i3 by:
		// for_window [instance="chromedp-runner"] move scratchpad
		chromedp.Flag("headless", false),
	}
	allocCtx, allocCancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer allocCancel()

	chromeCtx, chromeCancel := chromedp.NewContext(allocCtx)
	defer chromeCancel()

	r := strings.NewReplacer("{today}", schedule.Today, "{reportID}", reportID, "{jiraUser}", jiraUser)
	url := loginURL + U.PathEscape(r.Replace(jiraURL))

	var jwt string
	if err := chromedp.Run(chromeCtx, jiraLogin(url, &jwt)); err != nil {
		return nil, err
	}

	filterKey, err := getFilterKey(schedule.Today, jiraUser, jwt)
	if err != nil {
		return nil, err
	}

	pdf, err := exportReport(filterKey, jwt)
	if err != nil {
		return nil, err
	}

	return convertReport(pdf)
}

func jiraLogin(urlstr string, jwt *string) chromedp.Tasks {
	return chromedp.Tasks{
		debug("Opening browser."),
		chromedp.Navigate(urlstr),
		debug("Waiting login.."),
		chromedp.WaitVisible("#username", chromedp.ByID),
		chromedp.SendKeys("#username", email),
		debug("Continue"),
		chromedp.Click("#login-submit"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(2 * time.Second)
			return nil
		}),
		chromedp.WaitVisible("#password", chromedp.ByID),
		chromedp.SendKeys("#password", password),
		debug("Log in"),
		chromedp.Click("#login-submit"),

		debug("Waiting tempo.."),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				if tempCtx, src := getIframeContext(ctx); tempCtx != nil {
					*jwt = parseJwt(src)
					break
				}
				time.Sleep(1 * time.Second)
			}
			return nil
		}),
	}
}

func parseJwt(src string) string {
	rs, err := url.Parse(src)
	if err != nil {
		return ""
	}

	return rs.Query().Get("jwt")
}

func getIframeContext(ctx context.Context) (context.Context, string) {
	targets, _ := chromedp.Targets(ctx)
	var tgt *target.Info
	for _, t := range targets {
		if t.Type == "iframe" && strings.Contains(t.URL, "tempo.io") && t.TargetID != "" {
			tgt = t
			break
		}
	}
	if tgt != nil {
		ictx, _ := chromedp.NewContext(ctx, chromedp.WithTargetID(tgt.TargetID))
		return ictx, tgt.URL
	}

	return nil, ""
}

func debug(msg string, args ...interface{}) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		log.Debug("jira", msg, args...)
		return nil
	})
}
