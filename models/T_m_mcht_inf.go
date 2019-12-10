package models

type T_m_mcht_inf struct {
	Mcht_cd                   string
	Outside_mcht_cd           string
	Ins_id_cd                 string
	Name                      string
	Name_busi                 string
	Prov_cd                   string
	City_cd                   string
	County_cd                 string
	Reg_place                 string
	Reg_addr                  string
	Busi_lice_no              string
	Busi_rang                 string
	Busi_deadline             string
	Org_cd                    string
	Expiry_date               string
	Certif_policyholders_type string
	Policyholders_name        string
	Certif_type               string
	Certif_expiry_date        string
	Certif_no                 string
	Hotline                   string
	Mcht_type                 string
	Busi_type                 string
	Busi_type_name            string
	Contact                   string
	Mobile                    string
	Often_email               string
	Manager_name              string
	Emp_no                    string
	Rec_opr_id                string
	Rec_upd_opr               string
	Rec_crt_ts                string
	Rec_upd_ts                string
	Systemflag                string
	Status                    string
	Channel_no                string
	Wx_channel_id             string
	Head_ins_cd               string
}

func (t T_m_mcht_inf) TableName() string {
	return "t_m_mcht_inf"
}
