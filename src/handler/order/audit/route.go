package audit

import "github.com/cookieY/yee"

func AuditRestFulAPis() yee.RestfulAPI {
	return yee.RestfulAPI{
		Get:  AuditOSCFetchAndKillApis,
		Post: AuditOrderApis,
		Put:  AuditOrRecordOrderFetchApis,
	}
}

func OpenAuditRestFulAPis() yee.RestfulAPI {
	return yee.RestfulAPI{
		Post: OpenAuditOrderApis,
	}
}
