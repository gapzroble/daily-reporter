package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const waitScreenshot = 15 // seconds

var (
	novaURL   = "https://nova.fmi.filemaker-cloud.com/fmi/webd/nova%205"
	novaHours = "8,5"
)

func logNova(hours <-chan float64) ([]byte, error) {
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer timeoutCancel()

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
	}
	allocCtx, allocCancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer allocCancel()

	chromeCtx, chromeCancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer chromeCancel()

	setNovaHours(<-hours)

	var buf []byte
	tasks := logAndScreenshot(novaURL, 90, &buf)

	if err := chromedp.Run(chromeCtx, tasks); err != nil {
		return nil, err
	}

	return buf, nil
}

func logAndScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		debugNova("Opening browser."),
		chromedp.Navigate(urlstr),

		debugNova("Waiting login.."),
		chromedp.WaitVisible("#login_dialog_body", chromedp.ByID),
		chromedp.SendKeys("#login_name", username),
		chromedp.SendKeys("#login_pwd", password),
		debugNova("Login."),
		chromedp.Click("#login_nonguest"),

		debugNova("Waiting dashboard.."),
		chromedp.WaitVisible(`//span[contains(text(), "TIME REPORTING")]`),
		debugNova("Going to TIME REPORTING."),
		chromedp.Click(`//span[contains(text(), "TIME REPORTING")]/ancestor::button`),

		chromedp.WaitVisible(`//span[contains(text(), "Create new time report")]`),
		debugNova("Creating new time report."),
		chromedp.Click(`//span[contains(text(), "Create new time report")]/ancestor::button`),

		debugNova("Selecting %s project.", project),
		chromedp.WaitVisible(`.fm-combobox`),
		chromedp.Click(`.fm-combobox`),
		chromedp.WaitVisible(`//div[contains(text(), "` + project + `")]`),
		chromedp.Click(`//div[contains(text(), "` + project + `")]`),
		chromedp.WaitVisible(`//span[contains(text(), "CREATE REPORT")]`),
		chromedp.Click(`//span[contains(text(), "CREATE REPORT")]`),

		debugNova("Filling out report details.."),
		chromedp.Click(`//div[text()="Hours"]/../div`),
		chromedp.SendKeys(`[contentEditable=true]`, novaHours),
		chromedp.Click(`//div[contains(text(), "Write major activities done during the time period here")]/../div`),
		chromedp.SendKeys(`[contentEditable=true]`, details+"."),
		chromedp.Click(`//div[contains(text(), "Category")]/../..`),
		chromedp.Click(`//div[contains(text(), "Billable hours")]`),

		debugNova("Saving time report.."),
		chromedp.WaitVisible(`//span[contains(text(), "Save time report")]`),
		chromedp.Click(`//span[contains(text(), "Save time report")]/ancestor::button`),
		chromedp.WaitVisible(`//span[contains(text(), "Create new time report")]`),

		debugNova("Done report, preparing for screenshot.."),
		chromedp.ActionFunc(func(ctx context.Context) error {

			width, height := int64(1395), int64(985)

			// force viewport emulation
			err := emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// wait for the page to render properly
			time.Sleep(waitScreenshot * time.Second)

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      15,
					Y:      165,
					Width:  1000,
					Height: 61,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

func debugNova(msg string, args ...interface{}) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		log.Printf(msg+"\n", args...)
		return nil
	})
}

func setNovaHours(hours float64) {
	novaHours = fmt.Sprintf("%f", hours)
	novaHours = strings.TrimRight(novaHours, "0")
	novaHours = strings.TrimRight(novaHours, ".")
	novaHours = strings.ReplaceAll(novaHours, ".", ",")
}
