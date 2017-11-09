package modules

const ES_INDEX_ALERT = "apt"

const (
	SUCCESS                      ResCode = 0
	WRONG                        ResCode = 1
	ERR_EXPRESS                  ResCode = 2
	NOT_CLOSE_QUOTES_SINGLE      ResCode = 3
	NOT_CLOSE_QUOTES_DOUBLE      ResCode = 4
	NOT_FOUND_QUOTES_SINGLE_NEXT ResCode = 5
	TOKEN_IS_NOT_EXPRESS         ResCode = 6
	CodeEnd                      ResCode = 7
)

type ResCode int
type ExprErr string

var ErrExprType = [CodeEnd]ExprErr{
	NOT_CLOSE_QUOTES_SINGLE:      `not close '`,
	NOT_CLOSE_QUOTES_DOUBLE:      `not close "`,
	NOT_FOUND_QUOTES_SINGLE_NEXT: `no found single' next`,
	TOKEN_IS_NOT_EXPRESS:         `token isn't expr`,
}

const (
	Waf     AlertType = 0
	Vds     AlertType = 1
	Ids     AlertType = 2
	Multi   AlertType = 3
	TypeEnd AlertType = 4
)

type AlertType int

var EsType = [TypeEnd]string{
	Waf:   "waf_alert",
	Vds:   "vds_alert",
	Ids:   "ids_alert",
	Multi: "",
}

var WafItems = []string{
	"Client",
	"Rev",
	"Msg",
	"Attack",
	"Severity",
	"Maturity",
	"Accuracy",
	"Hostname",
	"Uri",
	"Unique_id",
	"Ref",
	"Tags",
	"Rule",
	"Version",
	"Xdr.Conn.Proto",
	"Xdr.Conn.Dip",
	"Xdr.Conn.Dport",
	"Xdr.Conn.DipInfo.Country",
	"Xdr.Conn.DipInfo.Province",
	"Xdr.Conn.DipInfo.City",
	"Xdr.Conn.DipInfo.Lat",
	"Xdr.Conn.DipInfo.Lng",
	"Xdr.Conn.Sip",
	"Xdr.Conn.Sport",
	"Xdr.Conn.SipInfo.Country",
	"Xdr.Conn.SipInfo.Province",
	"Xdr.Conn.SipInfo.City",
	"Xdr.Conn.SipInfo.Lat",
	"Xdr.Conn.SipInfo.Lng",
}

var VdsItems = []string{
	"Subfile",
	"Threatname",
	"Local_threatname",
	"Local_vtype",
	"Local_platfrom",
	"Local_vname",
	"Local_extent",
	"Local_enginetype",
	"Local_logtype",
	"Local_engineip",
	"Xdr.Time",
	"Xdr.Proto",
	"Xdr.Http.Url",
	"Xdr.App.File",
	"Xdr.Conn.Dip",
	"Xdr.Conn.Dport",
	"Xdr.Conn.DipInfo.Country",
	"Xdr.Conn.DipInfo.Province",
	"Xdr.Conn.DipInfo.City",
	"Xdr.Conn.DipInfo.Lat",
	"Xdr.Conn.DipInfo.Lng",
	"Xdr.Conn.Sip",
	"Xdr.Conn.Sport",
	"Xdr.Conn.SipInfo.Country",
	"Xdr.Conn.SipInfo.Province",
	"Xdr.Conn.SipInfo.City",
	"Xdr.Conn.SipInfo.Lat",
	"Xdr.Conn.SipInfo.Lng",
}

var IdsItems = []string{
	"Time",
	"Proto",
	"Byzoro_type",
	"Attack_type",
	"Details",
	"Severity",
	"Engine",
	"Dest_ip",
	"Dest_port",
	"Dest_ip_info.Country",
	"Dest_ip_info.Province",
	"Dest_ip_info.City",
	"Dest_ip_info.Lat",
	"Dest_ip_info.Lng",
	"Src_ip",
	"Src_port",
	"Src_ip_info.Country",
	"Src_ip_info.Province",
	"Src_ip_info.City",
	"Src_ip_info.Lat",
	"Src_ip_info.Lng",
}
