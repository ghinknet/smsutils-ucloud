# smsutils-ucloud

UCloud SMS (USMS) driver for [smsutils](https://go.gh.ink/smsutils).

This module implements the `model.Driver` / `model.Client` interfaces of the core
`smsutils` library on top of the UCloud `ucloud-sdk-go` USMS service.

Driver name: `ucloud`

## Installation

```bash
go get go.gh.ink/smsutils/ucloud/v3
```

## Usage

Blank-import this package so it registers itself, then configure it via the core client.

```go
import (
	"go.gh.ink/smsutils/v3/client"
	"go.gh.ink/smsutils/v3/model"

	_ "go.gh.ink/smsutils/ucloud/v3"
)

clients, err := client.NewClient(model.Config{
	Credentials: model.C{
		"ucloud": {
			"publicKey":  "<your-public-key>",
			"privateKey": "<your-private-key>",
			"projectId":  "<your-project-id>",
		},
	},
})
if err != nil {
	panic(err)
}

// "sender" is the SigContent; "UTA..." is the USMS TemplateId.
err = clients["ucloud"].SendMessage("+8617601205205", "<sign-content>", "UTA12345", model.Vars{
	{Key: "code", Value: "1234"},
})
```

## Credentials

| Key | Constant | Required | Description |
|-----|----------|----------|-------------|
| `publicKey` | `PublicKey` | Yes | UCloud API public key |
| `privateKey` | `PrivateKey` | Yes | UCloud API private key |
| `projectId` | `ProjectId` | Yes | UCloud project ID |

Missing any of these yields `errors.ErrDriverCredentialInvalid`.

## Behavior

- The destination is normalized via `utils.ProcessNumberForChinese`, then reformatted into
  UCloud's `(countryCode)nationalNumber` form (e.g. `(86)17601205205`).
- `template` maps to the USMS `TemplateId`; `sender` maps to `SigContent`.
- `vars` are passed **positionally** as `TemplateParams` (in the order supplied), so the order
  of `model.Vars` must match the placeholders in your template. The `Key` fields are ignored.
- A USMS `RetCode` error is mapped to `errors.ErrDriverSendFailed`, decorated with the
  provider code, message, request UUID and raw response. Other SDK server errors are returned
  as-is.
- The SDK log level is set to `Fatal` to suppress noisy output.

## Requirements

- Go 1.25.0+
- `go.gh.ink/smsutils/v3`
- UCloud SDK: `github.com/ucloud/ucloud-sdk-go`

## License

[Apache License 2.0](LICENSE)
