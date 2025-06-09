package postcart

import (
	"github.com/agladfield/postcart/pkg/postmark"
)

const (
	deliveriesTemplateAlias     = "postcart-deliveries-template"
	deliveredTemplateAlias      = "postcart-delivered-template"
	bounceTemplateAlias         = "postcart-bounce-template"
	spamComplaintTemplateAlias  = "postcart-spamc-template"
	invalidBalanceTemplateAlias = "postcart-balance-template"
	layoutTemplateAlias         = "postcart-layout-template"
)

// layout template
// bounce template
// delivered template
// delivery template

const deliveryHTMLTemplate = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="x-apple-disable-message-reformatting" />
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title></title>
    <style type="text/css" rel="stylesheet" media="all">
    /* Base ------------------------------ */
	@import url("https://fonts.bunny.net/css?family=open-sans:400,700");
    body {
      width: 100% !important;
      height: 100%;
      margin: 0;
      padding: 0;
      -webkit-text-size-adjust: none;
      background-color: #FFFFFF;
      color: #51545E;
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
    body,
    td,
    th {
      font-family: "Open Sans", Helvetica, Arial, sans-serif;
    }

    p.sub {
      font-size: 13px;
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
          <img src="{{encoded_image}}" style="display: block; width: 100%; height: auto; margin: 0; padding: 0; border: 0; outline: none; text-decoration: none; -ms-interpolation-mode: bicubic;" />
        </td>
      </tr>
      <!-- Footer -->
      <tr>
        <td class="email-footer">
          <p class="f-fallback sub align-center">sent with <a href="https://postc.art/">postc.art!</a></p>
        </td>
      </tr>
    </table>
  </body>
</html>`
const deliveryTextTemplate = "{{ascii_text}}"
const deliverySubjecTemplate = "📪 You've Received a Postcard: {{subject}}"

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
		return tmplErr
	}

	return nil
}

func checkTemplatesAreAvailable() error {
	// we will look for the alias
	listTemplates, listErr := postmark.ListTemplates()
	if listErr != nil {
		return listErr
	}

	deliveriesIncluded := false

	for _, tmpl := range listTemplates.Templates {
		switch tmpl.Alias {
		case deliveriesTemplateAlias:
			deliveriesIncluded = true
		case deliveredTemplateAlias:
		case bounceTemplateAlias:
		case spamComplaintTemplateAlias:
		case invalidBalanceTemplateAlias:
		case layoutTemplateAlias:
		default:
			continue
		}
	}

	if !deliveriesIncluded {
		createDeliveriesTemplate()
	}

	return nil
}
