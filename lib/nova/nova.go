package nova

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/rroble/daily-reporter/lib/log"
	"github.com/rroble/daily-reporter/lib/schedule"
)

// Log nova
func Log(loggable <-chan int64) ([]byte, error) {
	defer log.Debug("nova", "Done Nova")
	logSeconds := <-loggable
	// 0 is fine if holiday or leave
	if logSeconds <= 0 && !schedule.IsHolidayOrLeave() {
		return nil, errors.New("Already logged")
	}
	log.Debug("nova", "Logging Nova")
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer timeoutCancel()

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
	}
	allocCtx, allocCancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer allocCancel()

	chromeCtx, chromeCancel := chromedp.NewContext(allocCtx)
	defer chromeCancel()

	setNovaHours(logSeconds)

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
		// TODO: select date

		debugNova("Saving time report.."),
		chromedp.WaitVisible(`//span[contains(text(), "Save time report")]`),
		chromedp.Click(`//span[contains(text(), "Save time report")]/ancestor::button`),
		chromedp.WaitVisible(`//span[contains(text(), "Create new time report")]`),

		debugNova("Done report, preparing for screenshot.."),
		chromedp.ActionFunc(func(ctx context.Context) error {
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
		log.Debug("nova", msg, args...)
		return nil
	})
}

func setNovaHours(loggable int64) {
	novaHours = fmt.Sprintf("%f", float64(loggable)/3600)
	novaHours = strings.TrimRight(novaHours, "0")
	novaHours = strings.TrimRight(novaHours, ".")
	novaHours = strings.ReplaceAll(novaHours, ".", ",")
}
