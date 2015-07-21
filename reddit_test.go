package geddit

import (
	"fmt"
	"github.com/toqueteos/webbrowser"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func TestSubmit(t *testing.T) {

	session, err := NewLoginSession(
		"redditgolang",
		"apitest11",
		"tester",
	)
	if err != nil {
		t.Error(err)
	}

	subreddit, err := session.AboutSubreddit("mybottester")
	if err != nil {
		t.Error(err)
	}

	needsCaptcha, err := session.NeedsCaptcha()
	if err != nil {
		t.Error(err)
	}

	t.Log("Needs captcha:", needsCaptcha)

	if needsCaptcha {
		iden, err := session.NewCaptchaIden()
		if err != nil {
			t.Error(err)
		}

		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			fmt.Println("data is being piped to stdin")
		} else {
			fmt.Println("stdin is from a terminal")
		}

		//_, err = session.CaptchaImage(iden)
		//if err != nil {
		//	t.Error(err)
		//}
		webbrowser.Open("http://www.reddit.com/captcha/" + iden)
		fmt.Println("Enter captcha in captcha.txt")

		ind0 := strings.Index(os.Args[0], "github.com")
		ind1 := strings.Index(os.Args[0], "/_test")
		p0 := os.Getenv("GOPATH")
		p0 = path.Join(p0, "src", os.Args[0][ind0:ind1], "captcha.txt")
		fmt.Println(p0)
		time.Sleep(time.Second * 7)
		for {
			if _, err = os.Stat(p0); err != nil {
				time.Sleep(time.Second * 3)
			} else {
				break
			}
		}
		bb, err := ioutil.ReadFile(p0)
		if err != nil {
			t.Error(err)
		}
		captcha := strings.TrimSpace(string(bb))
		fmt.Println("CAPTCHA:", captcha)
		os.Remove(p0)

		err = session.Submit(NewTextSubmission(subreddit.Name, "CAPTCHA TESTING TEXT", "TEST TEXT", true, &Captcha{iden, captcha}))
		if err != nil {
			t.Error(err)
		}

		// there is a submission delay!
		// needs to find a way to test both without the time delay, maybe two accounts
		//
		//err = session.Submit(NewLinkSubmission(subreddit.Name, "CAPTCHA TESTING LINK", "https://github.com/jzelinskie/geddit", true, &Captcha{iden, captcha}))
		//if err != nil {
		//	t.Error(err)
		//}

	} else {

		err = session.Submit(NewTextSubmission(subreddit.Name, "TESTING TEXT", "TEST TEXT", true, &Captcha{}))
		if err != nil {
			t.Error(err)
		}

		// there is a submission delay!
		// needs to find a way to test both without the time delay, maybe two accounts
		//
		//err = session.Submit(NewLinkSubmission(subreddit.Name, "TESTING LINK", "https://github.com/jzelinskie/geddit", true, &Captcha{}))
		//if err != nil {
		//	t.Error(err)
		//}
	}

}

func TestListings(t *testing.T) {
	session, err := NewLoginSession(
		"redditgolang",
		"apitest11",
		"tester",
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = session.MySaved(NewSubmissions, "")
	if err != nil {
		t.Fatal(err)
	}
}
