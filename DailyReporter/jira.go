package main

import (
	"context"
	"errors"
	"log"
	"net/url"
	u "net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

var (
	reportID  = "f629bffe-eb81-4097-bc4b-c29c2f563090"
	jiraURL   = "https://arcanys.atlassian.net/plugins/servlet/ac/io.tempo.jira/tempo-app#!/reports/logged-time/{reportID}?columns=WORKED_COLUMN&dateDisplayType=days&from={today}&groupBy=project&groupBy=issue&groupBy=worklog&periodType=FIXED&subPeriodType=MONTH&to={today}&viewType=TIMESHEET&workerId={jiraUser}"
	exportURL = "https://app.tempo.io/rest/tempo-timesheets/4/worklogs/export/{filterKey}?format=pdf&title=Daily%2520Report&jwt={jwt}&groupBy=project,issue,worklog&columns=worked"
)

func jiraScreenshot(doneTempo <-chan bool) ([]byte, error) {
	if !<-doneTempo {
		return nil, errors.New("Tempo not logged")
	}

	debug("jira", "Creating screenshot")
	defer debug("jira", "Done screenshot")
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

	chromeCtx, chromeCancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer chromeCancel()

	r := strings.NewReplacer("{today}", today, "{reportID}", reportID, "{jiraUser}", jiraUser)
	url := "https://id.atlassian.com/login?continue=" + u.PathEscape(r.Replace(jiraURL))

	var jwt string
	if err := chromedp.Run(chromeCtx, jiraLogin(url, &jwt)); err != nil {
		return nil, err
	}

	filterKey, err := getFilterKey(today, jiraUser, jwt)
	if err != nil {
		return nil, err
	}

	return exportReport(filterKey, jwt)
}

func jiraLogin(urlstr string, jwt *string) chromedp.Tasks {
	return chromedp.Tasks{
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	debug("jira", "Resize..")
		// 	// force viewport emulation
		// 	err := emulation.SetDeviceMetricsOverride(width, height, 1, false).
		// 		WithScreenOrientation(&emulation.ScreenOrientation{
		// 			Type:  emulation.OrientationTypePortraitPrimary,
		// 			Angle: 0,
		// 		}).
		// 		Do(ctx)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	return nil
		// }),
		debugJira("Opening browser."),
		chromedp.Navigate(urlstr),
		debugJira("Waiting login.."),
		chromedp.WaitVisible("#username", chromedp.ByID),
		chromedp.SendKeys("#username", email),
		debugJira("Continue"),
		chromedp.Click("#login-submit"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(2 * time.Second)
			return nil
		}),
		chromedp.WaitVisible("#password", chromedp.ByID),
		chromedp.SendKeys("#password", jiraPassword),
		debugJira("Log in"),
		chromedp.Click("#login-submit"),

		debugJira("Waiting tempo.."),
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

// func createScreenshot(quality int64, res *[]byte) chromedp.Tasks {
// 	return chromedp.Tasks{
// 		debugJira("Check total hours.."),
// 		chromedp.WaitVisible(`//span[@class="tempoTotalHoursTotal"]`),
// 		chromedp.ActionFunc(func(ctx context.Context) error {
// 			var totalHours string
// 			chromedp.Text(`//span[@class="tempoTotalHoursTotal"]`, &totalHours, chromedp.NodeVisible)
// 			if totalHours == "0h" {
// 				return errors.New("No logs found")
// 			}
// 			debug("jira", "Got total hours: %s", totalHours)
// 			return nil
// 		}),

// 		debugJira("Preparing screenshot.."),
// 		chromedp.ActionFunc(func(ctx context.Context) error {
// 			// wait for the page to render properly
// 			time.Sleep(waitScreenshot * time.Second)

// 			x, y := int64(56), int64(100)

// 			// capture screenshot
// 			var err error
// 			debug("jira", "Capture..")
// 			*res, err = page.CaptureScreenshot().
// 				WithQuality(quality).
// 				WithClip(&page.Viewport{
// 					X:      float64(x),
// 					Y:      float64(y),
// 					Width:  float64(width - x),
// 					Height: float64(height - y),
// 					Scale:  1,
// 				}).Do(ctx)
// 			if err != nil {
// 				return err
// 			}
// 			return nil
// 		}),
// 	}
// }

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

func debugJira(msg string, args ...interface{}) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		debug("jira", msg, args...)
		return nil
	})
}
