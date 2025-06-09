package postmark

import (
	"fmt"
	"net/http"
)

const (
	listTemplatesPath    = "/templates"
	createTemplatePath   = listTemplatesPath
	validateTemplatePath = "/templates/validate"
)

// Shared Template Info

type TemplateInfo struct {
	Active         bool   `json:"Active"`
	TemplateID     int    `json:"TemplateId"`
	Name           string `json:"Name"`
	Alias          string `json:"Alias"`
	TemplateType   string `json:"TemplateType"`
	LayoutTemplate string `json:"LayoutTemplate,omitempty"`
}

// List Templates

type ListTemplatesResponse struct {
	TotalCount int            `json:"TotalCount"`
	Templates  []TemplateInfo `json:"Templates"`
}

func ListTemplates() (*ListTemplatesResponse, error) {
	var templatesRes ListTemplatesResponse

	reqErr := request(api, listTemplatesPath, http.MethodGet, nil, &templatesRes)
	if reqErr != nil {
		return nil, fmt.Errorf(postmarkAPIErrFmtStr, reqErr)
	}

	return &templatesRes, nil
}

// Create Template

type NewTemplate struct {
	Name           string `json:"Name"`
	Alias          string `json:"Alias"`
	HTMLBody       string `json:"HtmlBody"`
	TextBody       string `json:"TextBody"`
	Subject        string `json:"Subject"`
	TemplateType   string `json:"TemplateType,omitempty"`
	LayoutTemplate string `json:"LayoutTemplate,omitempty"`
}

func CreateTemplate(tmpl NewTemplate) (*TemplateInfo, error) {
	body, bodyErr := EncodeToStruct(tmpl)
	if bodyErr != nil {
		return nil, fmt.Errorf(postmarkAPIErrFmtStr, bodyErr)
	}

	var templateRes TemplateInfo

	reqErr := request(api, createTemplatePath, http.MethodPost, body, &templateRes)
	if reqErr != nil {
		return nil, fmt.Errorf(postmarkAPIErrFmtStr, reqErr)
	}

	return &templateRes, nil
}

// Validate Template

type TemplateValidationRequest[T any] struct {
	Subject                    string `json:"Subject"`
	HTMLBody                   string `json:"HtmlBody"`
	TextBody                   string `json:"TextBody"`
	TestRenderModel            T      `json:"TestRenderModel"`
	InlineCSSForHTMLTestRender bool   `json:"InlineCssForHtmlTestRender,omitempty"`
	TemplateType               string `json:"TemplateType,omitempty"`
	LayoutTemplate             string `json:"LayoutTemplate"`
}

type TemplateTargetValidationError struct {
	Message           string `json:"Message"`
	Line              int    `json:"Line"`
	CharacterPosition int    `json:"CharacterPosition"`
}

type TemplateTargetValidationResult struct {
	ContentIsValid   bool                            `json:"ContentIsValid"`
	ValidationErrors []TemplateTargetValidationError `json:"ValidationErrors"`
}

type TemplateValidationResponse struct {
	AllContentIsValid      bool                           `json:"AllContentIsValid"`
	TextBody               TemplateTargetValidationResult `json:"TextBody,omitempty"`
	HTMLBody               TemplateTargetValidationResult `json:"HtmlBody,omitempty"`
	Subject                TemplateTargetValidationResult `json:"Subject,omitempty"`
	SuggestedTemplateModel map[string]any                 `json:"SuggestedTemplateModel"`
}

func ValidateTemplate[T any](req TemplateValidationRequest[T]) (*TemplateValidationResponse, error) {
	body, bodyErr := EncodeToStruct(req)
	if bodyErr != nil {
		return nil, fmt.Errorf(postmarkAPIErrFmtStr, bodyErr)
	}

	var validateRes TemplateValidationResponse

	reqErr := request(api, validateTemplatePath, http.MethodPost, body, &validateRes)
	if reqErr != nil {
		return nil, fmt.Errorf(postmarkAPIErrFmtStr, reqErr)
	}

	return &validateRes, nil
}
