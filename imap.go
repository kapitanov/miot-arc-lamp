package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mxk/go-imap/imap"
)

var (
	imapAddr     string
	imapUsername string
	imapPassword string
	imapLabel    string

	imapClient     *imap.Client
	lastImapStatus *imapStatus
)

func init() {
	imapAddr = os.Getenv("IMAP_ADDR")
	imapUsername = os.Getenv("IMAP_USERNAME")
	imapPassword = os.Getenv("IMAP_PASSWORD")
	imapLabel = os.Getenv("IMAP_LABEL")
}

type imapStatus struct {
	unreadCount uint32
}

func fetchImapStatus() (*imapStatus, error) {
	if imapClient == nil {
		err := imapConnect()
		if err != nil {
			return nil, err
		}
	}

	cmd, err := check(imapClient.Status(imapLabel))
	if err != nil {
		fmt.Fprintf(os.Stderr, "imap: status error. %s\n", err)
		imapClient = nil
		return nil, err
	}

	var count uint32
	for _, result := range cmd.Data {
		mailboxStatus := result.MailboxStatus()
		if mailboxStatus != nil {
			count += mailboxStatus.Unseen
		}
	}

	status := &imapStatus{count}
	if lastImapStatus == nil || lastImapStatus.unreadCount != status.unreadCount {
		fmt.Fprintf(os.Stdout, "imap: %d unread messages\n", status.unreadCount)
	}

	lastImapStatus = status
	return status, nil
}

func imapConnect() error {
	var err error

	fmt.Fprintf(os.Stdout, "imap: dialing %s\n", imapAddr)
	if strings.HasSuffix(imapAddr, ":993") {
		imapClient, err = imap.DialTLS(imapAddr, nil)
	} else {
		imapClient, err = imap.Dial(imapAddr)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "imap: failed to dial. %s\n", err)
		imapClient = nil
		return err
	}

	if imapClient.Caps["STARTTLS"] {
		fmt.Fprintf(os.Stdout, "imap: starting tls\n")
		_, err = check(imapClient.StartTLS(nil))
		if err != nil {
			fmt.Fprintf(os.Stderr, "imap: tls error. %s\n", err)
			imapClient = nil
			return err
		}
	}

	fmt.Fprintf(os.Stdout, "imap: logging in as %s\n", imapUsername)
	_, err = check(imapClient.Login(imapUsername, imapPassword))
	if err != nil {
		fmt.Fprintf(os.Stderr, "imap: login error. %s\n", err)
		imapClient = nil
		return err
	}

	return nil
}

func check(cmd *imap.Command, err error) (*imap.Command, error) {
	if err != nil {
		return nil, err
	}

	_, err = cmd.Result(imap.OK)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
