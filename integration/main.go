package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/tiqqe/go-logger"
	"github.com/tiqqe/go-s3helper"
)

var bucket string

func init() {
	bucket = os.Getenv("Bucket")
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lctx, _ := lambdacontext.FromContext(ctx)
	logger.Init(lctx.AwsRequestID, os.Getenv("AWS_LAMBDA_FUNCTION_NAME"))

	defer handlePanic()

	logger.InfoStringf("headers: %+v", event.Headers)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), chromedp.Flag("headless", false))
	defer cancel()

	// create context
	chromeCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(logger.InfoStringf))
	defer cancel()

	var buf []byte

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(chromeCtx, fullScreenshot(`https://nova.fmi.filemaker-cloud.com/fmi/webd/nova%205`, 90, &buf)); err != nil {
		log.Fatal(err)
	}

	s3 := s3helper.New(ctx, bucket)
	dest := fmt.Sprintf("screenshot-%s.png", lctx.AwsRequestID)

	if err := s3.Upload(dest, buf); err != nil {
		logger.Error(&logger.LogEntry{
			Message:      "Failed to upload api data to s3",
			ErrorMessage: err.Error(),
			Keys: map[string]interface{}{
				"Bucket": bucket,
				"Key":    dest,
			},
		})
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil

}
func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible("#login_dialog_body", chromedp.ByID),
		chromedp.SendKeys("#login_name", "Randolph Roble"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

func main() {
	lambda.Start(handler)
}

func handlePanic() {
	msg := recover()
	if msg != nil {
		entry := &logger.LogEntry{
			Message:   "Go panic",
			ErrorCode: "GoPanic",
		}
		switch msg := msg.(type) {
		case string:
			entry.ErrorMessage = msg
		case error:
			entry.ErrorMessage = msg.Error()

		default:
			entry.ErrorCode = "Unknown error type"
			entry.SetKey("error", msg)
		}

		logger.Error(entry)
	}
}
