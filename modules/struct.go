package modules

type ResIds struct {
	Time     string
	Attack   string
	Details  string
	Severity string
	Engine   string
	Type     string
	ConnInfo
}

type ResWaf struct {
	Time      string
	Client    string
	Rev       string
	Msg       string
	Attack    string
	Severity  string
	Maturity  int32
	Accuracy  int32
	Hostname  string
	Uri       string
	Unique_id string
	Ref       string
	Tags      []string
	Rule      WafAlertRule
	Version   string
	Operators string
	Type      string
	ConnInfo
}

type ResVds struct {
	Time             string
	Subfile          string
	Threatname       string
	Local_threatname string
	Attack           string
	Local_platfrom   string
	Local_vname      string
	Severity         string
	Local_enginetype string
	Local_logtype    string
	Local_engineip   string
	HttpUrl          string
	Filepath         string
	Type             string
	ConnInfo
}

type SrcVds struct {
	TimeAppend       string
	SeverityAppend   string
	Subfile          string
	Threatname       string
	Local_threatname string
	Local_vtype      string
	Attack           string
	Local_platfrom   string
	Local_vname      string
	Severity         string
	Local_enginetype string
	Local_logtype    string
	Local_engineip   string
	Xdr              []BackendObj
}

type SrcWaf struct {
	ProtoAppend    string
	Client         string
	Rev            string
	Msg            string
	Attack         string
	Severity       int32
	SeverityAppend string
	Maturity       int32
	Accuracy       int32
	Hostname       string
	Uri            string
	Unique_id      string
	Ref            string
	Tags           []string
	Rule           WafAlertRule
	Version        string
	Xdr            []BackendObj
}

type SrcIds struct {
	Time           uint64
	Src_ip         string
	Src_ip_info    IpInfo
	Src_port       uint16
	Dest_ip        string
	Dest_ip_info   IpInfo
	Dest_port      uint16
	Proto          uint8
	Byzoro_type    string
	Attack         string
	Attack_type    string
	Details        string
	Severity       uint32
	SeverityAppend string
	Engine         string
	Xdr            []BackendObj
}
type ConnInfo struct {
	Src_ip       string
	Src_port     uint32
	Src_ip_info  IpInfo
	Dest_ip      string
	Dest_port    uint32
	Dest_ip_info IpInfo
	Proto        string
}

type BackendObj struct {
	TimeAppend string
	Conn       Conn_backend `json:Conn`
	Http       struct {
		Url string `json:Url`
	}
	App struct {
		File string `json:"File"`
	} `json:"App"`
}

type Conn_backend struct {
	ProtoAppend string
	Sport       uint32 `json:Sport`
	Dport       uint32 `json:Dport`
	Sip         string `json:Sip`
	SipInfo     IpInfo `json:SipInfo`
	Dip         string `json:Dip`
	DipInfo     IpInfo `json:DipInfo`
}

type WafAlertRule struct {
	Ver  string
	Data string
	File string
	Line uint64
	Id   uint64
}

type IpInfo struct {
	Country  string `json:"Country"`
	Province string `json:"Province"`
	City     string `json:"City"`
	Lng      string `json:"Lng"`
	Lat      string `json:"Lat"`
}
