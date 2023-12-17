package eseis

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/idkw/eseisscrapper/pkg/infrastructure/chrome"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"text/template"
)

type pdfRes struct {
	buffer *[]byte
}

type ChromeDp struct {
	ctx *context.Context
}

func newChrome(URL string, username string, password string) (*chrome.Chrome, error) {
	c := chrome.NewChrome()
	err := login(c, URL, username, password)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func login(c *chrome.Chrome, URL string, username string, password string) error {
	tasks := chromedp.Tasks{
		chromedp.EmulateViewport(799, 799),
		chromedp.Navigate(URL),
		chromedp.WaitVisible("#login-username"),
		chromedp.SendKeys(`#login-username`, username+kb.Enter),
		chromedp.WaitVisible("#login-password"),
		chromedp.SendKeys(`#login-password`, password+kb.Enter),
		chromedp.WaitVisible(".sc-eHWfIC"), // wait for co-owner balance to be visible
	}
	return c.RunTasks(tasks)
}

func (e *EseisClient) SavePDF(URL string, outPath string) error {
	var pdfRes = pdfRes{}
	if err := e.chromeSession.RunTasks(printPdfTasks(URL, &pdfRes)); err != nil {
		logrus.Fatal(err)
	}
	if err := os.WriteFile(outPath, *pdfRes.buffer, 0o644); err != nil {
		logrus.Fatal(err)
	}
	return nil
}

func printPdfTasks(urlstr string, res *pdfRes) chromedp.Tasks {
	tasks := chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitReady(".sc-jQAxuV"), // wait for title to be visible
		chromedp.WaitReady(".sc-eDdKWq"), // wait for author to be visible
		chromedp.WaitReady(".sc-eHEENL"), // wait for report description to be visible
		chromedp.WaitReady(".sc-dWBRfb"), // wait for comment to be visible
	}
	//removeConsentAction := NewRemoveElementAction("#appconsent") // remove cookie banner
	//if removeConsentAction != nil {
	//	tasks = append(tasks, removeConsentAction)
	//}

	removeButtonsAction := NewRemoveElementAction("ul[role=\"navigation\"]") // remove action buttons
	if removeButtonsAction != nil {
		tasks = append(tasks, removeButtonsAction)
	}

	printPdfAction := chromedp.ActionFunc(func(ctx context.Context) error {
		buf, _, err := page.PrintToPDF().
			Do(ctx)
		if err != nil {
			return err
		}
		res.buffer = &buf
		return nil
	})
	tasks = append(tasks, printPdfAction)

	return tasks
}

func NewRemoveElementAction(selector string) chromedp.Action {
	tpl, err := template.New("js").
		Parse("let {{ .Node }} = document.querySelector('{{ .Selector }}'); {{ .Node }}.parentNode.removeChild({{ .Node }}); true;")
	if err != nil {
		logrus.Warnf("Failed to create template to remove element %s", selector)
		return nil
	}
	type ctx struct {
		Node, Selector string
	}
	buffer := bytes.Buffer{}
	c := ctx{
		Node:     fmt.Sprintf("node%d", rand.Int()),
		Selector: selector,
	}
	err = tpl.Execute(&buffer, c)
	if err != nil {
		logrus.Warnf("Failed to execute template to remove element %s", selector)
		return nil
	}
	javascript := buffer.String()
	return chromedp.Evaluate(javascript, nil)
}
