package modules

const ES_INDEX_ALERT = "apt"

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
	"Xdr.Conn.Dip",
	"Xdr.Conn.Dport",
	"Xdr.Conn.Sip",
	"Xdr.Conn.SipInfo.Country",
	"Xdr.Conn.SipInfo.Province",
	"Xdr.Conn.SipInfo.City",
}

var VdsItems = []string{
	"Subfile",
	"Log_time",
	"Threatname",
	"Local_threatname",
	"Local_vtype",
	"Local_platfrom",
	"Local_vname",
	"Local_extent",
	"Local_enginetype",
	"Local_logtype",
	"Local_engineip",
	"Xdr.Http.URL",
	"Xdr.App.File",
	"Xdr.Conn.Dip",
	"Xdr.Conn.Dport",
	"Xdr.Conn.Sip",
	"Xdr.Conn.SipInfo.Country",
	"Xdr.Conn.SipInfo.Province",
	"Xdr.Conn.SipInfo.City",
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
	"Src_ip",
	"Src_ip_info.Country",
	"Src_ip_info.Province",
	"Src_ip_info.City",
}
