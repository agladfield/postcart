package jdb

func AddBlockedSendEmail(email string) error {
	//
	return nil
}

func CheckForBlockedOutboundEmail(email string) bool {
	return false
}

func CheckForBlockedInboundEmail(email string) bool {
	return false
}

func AddBlockedInboundEmail(email string) error {
	//
	return nil
}

func IncrementSentEmail(email string) (int, error) {
	//
	return 0, nil
}
