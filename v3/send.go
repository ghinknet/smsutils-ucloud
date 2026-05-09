package ucloud

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	uerr "github.com/ucloud/ucloud-sdk-go/ucloud/error"
	smsutilsError "go.gh.ink/smsutils/v3/errors"
	"go.gh.ink/smsutils/v3/model"
	"go.gh.ink/smsutils/v3/utils"
)

func (c Client) SendMessage(dest string, sender string, template string, vars model.Vars) error {
	// Try to parse number
	_, countryCode, nationalNumber, _, err := utils.ProcessNumberForChinese(dest)
	if err != nil {
		return err
	}
	dest = strings.Join([]string{
		"(", strconv.FormatInt(countryCode, 10), ")",
		strconv.FormatInt(nationalNumber, 10),
	}, "")

	// Preprocess vars
	params := make([]string, len(vars))
	for k, v := range vars {
		params[k] = v.Value
	}

	req := c.Client.NewSendUSMSMessageRequest()
	req.SigContent = ucloud.String(sender)
	req.TemplateId = ucloud.String(template)
	req.PhoneNumbers = []string{
		dest,
	}
	req.TemplateParams = params
	resp, err := c.Client.SendUSMSMessage(req)
	if err != nil {
		if e, ok := errors.AsType[uerr.ServerError](err); ok {
			if e.Name() == uerr.ErrRetCode {
				return smsutilsError.ErrDriverSendFailed.
					WithDriverName(Name).
					WithDriverCode(strconv.Itoa(resp.GetRetCode())).
					WithDriverMessage(resp.GetMessage()).
					WithDriverRequestID(resp.GetRequestUUID()).
					WithDriverResponse(resp)
			}
			return e
		}
		return err
	}

	return nil
}
