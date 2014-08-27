package deis

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var NotAllowed = errors.New("Not authorized")

// ReceiverAuth represents the structure sent to the controller when making a reciever
// authentication call.
type ReceiverAuth struct {
	Username        string `json:"receive_user"`
	Repo            string `json:"receive_repo"`
	Sha             string `json:"sha"`
	Fingerprint     string `json:"fingerprint"`
	SshConnection   string `json:"ssh_connection"`
	OriginalCommand string `json:"ssh_original_command"`
}

// String satisfies fmt.Stringer.
func (r *ReceiverAuth) String() string {
	res, err := json.Marshal(r)
	if err != nil {
		return ""
	}

	return string(res)
}

// CheckAuthForReceiver checks authentication in the Deis controller for permission
// to be able to make a build.
func (c *Controller) CheckAuthForReceiver(r *ReceiverAuth) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", c.GetURL(), strings.NewReader(r.String()))
	if err != nil {
		return err
	}

	req.Header.Add("X-Deis-Builder-Auth", c.BuildKey)

	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		return NotAllowed
	}

	return nil
}
