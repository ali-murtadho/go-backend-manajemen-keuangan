package models

func (Waktu) TableName() string {
	return "m_waktu"
}

type Waktu struct {
	Id         int    `json:"id"`
	BulanTahun string `json:"bulan_tahun"`
}

type WaktuRequest struct {
	BulanTahun string `json:"bulan_tahun"`
}

type WaktuResponse struct {
	Id         int    `json:"id"`
	BulanTahun string `json:"bulan_tahun"`
}

type WaktuUpdateRequest struct {
	BulanTahun string `json:"bulan_tahun"`
}

type WaktuUpdateResponse struct {
	Id         int    `json:"id"`
	BulanTahun string `json:"bulan_tahun"`
}

type WaktuDeleteResponse struct {
	Id         int    `json:"id"`
	BulanTahun string `json:"bulan_tahun"`
}
