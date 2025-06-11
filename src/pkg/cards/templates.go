package cards

import (
	"fmt"

	"github.com/agladfield/postcart/pkg/postmark"
)

const deliveriesTemplateAlias = "postcart-deliveries-template"

const deliveryHTMLTemplate = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="x-apple-disable-message-reformatting" />
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title></title>
    <style type="text/css" rel="stylesheet" media="all">
    /* Bunny Fonts Import */
    @import url("https://fonts.bunny.net/css?family=open-sans:400,700");

    /* Base ------------------------------ */
    body {
      width: 100% !important;
      height: 100%;
      margin: 0;
      padding: 0;
      -webkit-text-size-adjust: none;
      background-color: #FFFFFF;
      color: #51545E;
      font-family: "Open Sans", Helvetica, Arial, sans-serif;
    }

    a {
      color: #3869D4;
    }

    a img {
      border: none;
    }

    td {
      word-break: break-word;
      padding: 0;
    }

    table {
      border-collapse: collapse;
    }

    /* Type ------------------------------ */
    p.sub {
      font-size: 13px;
      font-family: "Open Sans", Helvetica, Arial, sans-serif;
    }

    /* Utilities ------------------------------ */
    .align-center {
      text-align: center;
    }

    /* Footer ------------------------------ */
    .email-footer {
      width: 100%;
      margin: 0;
      padding: 20px 0;
      background-color: #FFFFFF;
      text-align: center;
    }

    .email-footer p {
      color: #51545E;
      margin: 0;
      font-family: "Open Sans", Helvetica, Arial, sans-serif;
    }

    /* Media Queries ------------------------------ */
    @media only screen and (max-width: 600px) {
      .email-footer {
        width: 100% !important;
      }
    }
    </style>
    <!--[if mso]>
    <style type="text/css">
      .f-fallback {
        font-family: Arial, sans-serif;
      }
    </style>
    <![endif]-->
  </head>
  <body>
    <table width="100%" cellpadding="0" cellspacing="0" role="presentation" style="margin: 0; padding: 0;">
      <!-- Image as main content -->
      <tr>
        <td>
          <img src="{{image_url}}" style="display: block; width: 100%; height: auto; margin: 0; padding: 0; border: 0; outline: none; text-decoration: none; -ms-interpolation-mode: bicubic;" />
        </td>
      </tr>
      <!-- Footer -->
      <tr>
        <td class="email-footer">
          <p class="f-fallback sub align-center">sent with <a href="https://postc.art">postc.art</a></p>
        </td>
      </tr>
    </table>
  </body>
</html>`
const deliveryTextTemplate = "{{ascii_text}}"
const deliverySubjecTemplate = "📪 {{subject}}"

const templatesErrFmtStr = "cards templates err: %w"

func createDeliveriesTemplate() error {
	deliveriesTmpl := postmark.NewTemplate{
		Name:     "Deliveries Template",
		Alias:    deliveriesTemplateAlias,
		HTMLBody: deliveryHTMLTemplate,
		TextBody: deliveryTextTemplate,
		Subject:  deliverySubjecTemplate,
	}
	_, tmplErr := postmark.CreateTemplate(deliveriesTmpl)
	if tmplErr != nil {
		return fmt.Errorf(templatesErrFmtStr, tmplErr)
	}

	return nil
}

const (
	listCount  = 100
	listOffset = 0
)

func checkTemplatesAreAvailable() error {
	// we will look for the alias
	listTemplates, listErr := postmark.ListTemplates(listCount, listOffset)
	if listErr != nil {
		return fmt.Errorf(templatesErrFmtStr, listErr)
	}

	deliveriesIncluded := false

	for _, tmpl := range listTemplates.Templates {
		switch tmpl.Alias {
		case deliveriesTemplateAlias:
			deliveriesIncluded = true
		default:
			continue
		}
	}

	if !deliveriesIncluded {
		deliveryTempErr := createDeliveriesTemplate()
		if deliveryTempErr != nil {
			return fmt.Errorf(templatesErrFmtStr, deliveryTempErr)
		}
	}

	return nil
}

// © Arthur Gladfield
