package nova

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/rroble/daily-reporter/lib/log"
	"github.com/rroble/daily-reporter/lib/models"
	"github.com/rroble/daily-reporter/lib/schedule"
)

var chromeCtx context.Context
var chromeCancel, timeoutCancel, allocCancel context.CancelFunc
var simulation bool

func init() {
	if val := os.Getenv("SIMULATION"); val != "" {
		simulation = true
	}
}

// Init nova
func Init() {
	log.Debug("nova", "Loading Nova")
	var timeoutCtx context.Context
	timeoutCtx, timeoutCancel = context.WithTimeout(context.Background(), 120*time.Second)

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
	}
	var allocCtx context.Context
	allocCtx, allocCancel = chromedp.NewExecAllocator(timeoutCtx, opts...)

	chromeCtx, chromeCancel = chromedp.NewContext(allocCtx)
}

// End func
func End() {
	timeoutCancel()
	allocCancel()
	chromeCancel()
}

// LogFromTempo to nova
func LogFromTempo(worklog models.Worklog) error {
	tasks := append(login(), logReport(worklog))

	if len(tasks) == 0 {
		log.Debug("nova", "Nothing to log: %+v", worklog)
		return nil
	}

	if err := chromedp.Run(chromeCtx, tasks); err != nil {
		return err
	}

	return nil
}

func logReport(worklog models.Worklog) chromedp.Tasks {
	hours := toNovaHours(worklog.TimeSpentSeconds)
	project := "TIQ Internal time"
	details := worklog.Description
	category := "Billable hours"
	nextcat := false
	ionumber := ""

	switch worklog.Issue.Key {
	case "TIQ-957":
		project = "PN Cleveron Project"
		nextcat = true
		ionumber = "7014750/47200/Parcel Robotics 100% Randolph Roble"
	case "TIQ-1095":
		category = "Non billable hours"
	case "TIQ-621":
		project = "SO SysOps"
	case "TIQ-705":
		project = "SO Implementation NK"
		nextcat = true
	case "TIQ-1075":
		details = "Varner: " + details
		project = "SO Change requests"
	case "TIQ-684":
		project = "SO Byggmax"
	}

	button := "SAVE"
	if simulation {
		button = "CANCEL"
	}

	return chromedp.Tasks{
		chromedp.WaitVisible(`//span[contains(text(), "Create new time report")]`),
		debugNova("Creating new time report."),
		chromedp.Click(`//span[contains(text(), "Create new time report")]/ancestor::button`),
		chromedp.WaitVisible(`//span[contains(text(), "CREATE REPORT")]`),

		debugNova("Selecting %s project.", project),
		chromedp.WaitVisible(`.fm-combobox`),
		chromedp.Click(`.fm-combobox`),
		(func() chromedp.Tasks {
			if !nextcat {
				return chromedp.Tasks{}
			}
			return chromedp.Tasks{
				chromedp.Click(`.nextpage span`),
			}
		}()),
		chromedp.WaitVisible(`//div[contains(text(), "` + project + `")]`),
		chromedp.Click(`//div[contains(text(), "` + project + `")]`),
		chromedp.Click(`//span[contains(text(), "CREATE REPORT")]`),

		debugNova("Filling out report details.."),
		chromedp.Click(`//div[contains(text(), "Category")]/../..`),
		chromedp.Click(`//div[contains(text(), "` + category + `")]`),

		chromedp.Click(`//div[text()="Hours"]/../div`),
		chromedp.SendKeys(`[contentEditable=true]`, hours),

		(func() chromedp.Tasks {
			if ionumber == "" {
				return chromedp.Tasks{}
			}
			return chromedp.Tasks{
				chromedp.Click(`//div[contains(text(), "IO Number")]`),
				chromedp.Click(`//div[contains(text(), "` + ionumber + `")]`),
			}
		}()),

		chromedp.Click(`//div[contains(text(), "Write major activities done during the time period here")]/../div`),
		chromedp.SendKeys(`[contentEditable=true]`, details+"."),
		chromedp.Blur(`[contentEditable=true]`),

		// select date
		(func() chromedp.Tasks {
			return chromedp.Tasks{
				chromedp.Click(`//div[contains(@class, "fm-datefield")][1]`),
				chromedp.SendKeys(`[contentEditable=true]`, strings.ReplaceAll(schedule.Today, "-", "/")),
				chromedp.Blur(`[contentEditable=true]`),
				chromedp.Click(`//div[contains(@class, "fm-datefield")][2]`),
				chromedp.SendKeys(`[contentEditable=true]`, strings.ReplaceAll(schedule.Today, "-", "/")),
				chromedp.Blur(`[contentEditable=true]`),
			}
		}()),

		debugNova("Saving time report.."),
		chromedp.WaitVisible(`//span[contains(text(), "` + button + `")]`),
		chromedp.Click(`//span[contains(text(), "` + button + `")]/ancestor::button`),
		chromedp.WaitVisible(`//span[contains(text(), "Create new time report")]`),

		debugNova("Done report.."),
	}
}

var loggedIn bool

func login() chromedp.Tasks {
	if loggedIn {
		return chromedp.Tasks{}
	}

	return chromedp.Tasks{
		debugNova("Opening browser."),
		chromedp.Navigate(novaURL),

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

		chromedp.ActionFunc(func(ctx context.Context) error {
			loggedIn = true
			return nil
		}),
	}
}

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

func toNovaHours(loggable int64) string {
	var hours string

	hours = fmt.Sprintf("%f", float64(loggable)/3600)
	hours = strings.TrimRight(hours, "0")
	hours = strings.TrimRight(hours, ".")
	hours = strings.ReplaceAll(hours, ".", ",")

	return hours
}
