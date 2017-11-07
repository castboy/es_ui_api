package modules

type IdsAlert struct {
	Time         uint64
	Src_ip       string
	Src_port     uint32
	Src_ip_info  IpInfo
	Dest_ip      string
	Dest_port    uint32
	Dest_ip_info IpInfo
	Proto        uint32
	Byzoro_type  string
	Attack_type  string
	Details      string
	Severity     uint32
	Engine       string
}

type WafAlert struct {
	Client    string
	Rev       string
	Msg       string
	Attack    string
	Severity  int32
	Maturity  int32
	Accuracy  int32
	Hostname  string
	Uri       string
	Unique_id string
	Ref       string
	Tags      []string
	Rule      WafAlertRule
	Version   string
}

type VdsAlert struct {
	Subfile          string
	Threatname       string
	Local_threatname string
	Local_vtype      string
	Local_platfrom   string
	Local_vname      string
	Local_extent     string
	Local_enginetype string
	Local_logtype    string
	Local_engineip   string
}

type BackendObj struct {
	Time uint64       `json:Time`
	Conn Conn_backend `json:Conn`
	Http struct {
		Url string `json:Url`
	}
	App struct {
		File string `json:"File,omitempty"`
	} `json:"App,omitempty"`
}

type Conn_backend struct {
	Proto   uint8  `json:Proto`
	Sport   uint16 `json:Sport`
	Dport   uint16 `json:Dport`
	Sip     string `json:Sip`
	SipInfo IpInfo `json:SipInfo`
	Dip     string `json:Dip`
	DipInfo IpInfo `json:DipInfo`
}

type WafAlertRule struct {
	Ver  string
	Data string
	File string
	Line uint64
	Id   uint64
}

type IpInfo struct {
	Country  string `json:"Country,omitempty"`
	Province string `json:"Province,omitempty"`
	City     string `json:"City,omitempty"`
	Lng      string `json:"Lng,omitempty"`
	Lat      string `json:"Lat,omitempty"`
}

type WafSource struct {
	WafAlert
	Xdr []BackendObj
}

type VdsSource struct {
	VdsAlert
	Xdr []BackendObj
}

type IdsSource struct {
	IdsAlert
}

type ApiWafRes struct {
	WafAlert
	Proto        uint8
	Time         uint64
	Dest_ip      string
	Dest_port    uint16
	Dest_ip_info IpInfo
	Src_ip       string
	Src_port     uint16
	Src_ip_info  IpInfo
	Operators    string
}

type ApiVdsRes struct {
	VdsAlert
	Time         uint64
	Proto        uint8
	Severity     string
	HttpUrl      string
	Filepath     string
	Dest_ip      string
	Dest_port    uint16
	Dest_ip_info IpInfo
	Src_ip       string
	Src_port     uint16
	Src_ip_info  IpInfo
}

type ApiIdsRes struct {
	IdsAlert
	Dest_ip     string
	Dest_port   uint32
	Src_ip      string
	Src_ip_info IpInfo
}
