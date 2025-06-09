package postmark

import (
	"fmt"
	"net/http"
)

// create and delete inbound triggers for banning, unbanning users

const (
	createInboundTriggerPath       = "/triggers/inboundrules"
	deleteInboundTriggerPathFmtStr = "/triggers/inboundrules/%d"
)

type InboundTriggerRule struct {
	Rule string `json:"Rule"`
}

type InboundTriggerRuleResponse struct {
	ID   int    `json:"ID"`
	Rule string `json:"Rule"`
}

func CreateInboundTriggerRule(emailOrDomain string) (*InboundTriggerRuleResponse, error) {
	req := InboundTriggerRule{
		Rule: emailOrDomain,
	}
	body, bodyErr := EncodeToStruct(req)
	if bodyErr != nil {
		return nil, fmt.Errorf(postmarkAPIErrFmtStr, bodyErr)
	}

	var triggerRuleRes InboundTriggerRuleResponse

	reqErr := request(api, validateTemplatePath, http.MethodPost, body, &triggerRuleRes)
	if reqErr != nil {
		return nil, fmt.Errorf(postmarkAPIErrFmtStr, reqErr)
	}

	return &triggerRuleRes, nil
}

// Delete Inbound Trigger Rule

type DeleteInboundTriggerResponse struct {
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}

func DeleteInboundTriggerRule(ruleID int) (*DeleteInboundTriggerResponse, error) {
	deleteURLPath := fmt.Sprintf(deleteInboundTriggerPathFmtStr, ruleID)

	var deleteRes DeleteInboundTriggerResponse

	reqErr := request(api, deleteURLPath, http.MethodDelete, nil, &deleteRes)
	if reqErr != nil {
		return nil, fmt.Errorf(postmarkAPIErrFmtStr, reqErr)
	}

	return &deleteRes, nil
}
