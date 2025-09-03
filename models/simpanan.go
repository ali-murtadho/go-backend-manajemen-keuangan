package models

// CREATE TABLE m_simpanan (
//     id SERIAL PRIMARY KEY,
//     jenis_simpanan VARCHAR(50) NOT NULL, -- bank konv, bank digital, ewallet
//     nama VARCHAR(100) NOT NULL,
//     warna VARCHAR(20)
//     isMain BOOLEAN DEFAULT FALSE
// );

func (Simpanan) TableName() string {
	return "m_simpanan"
}

type Simpanan struct {
	JenisSimpanan string `json:"jenis_simpanan"`
	Nama          string `json:"nama"`
	Warna         string `json:"warna"`
	IsMain        bool   `json:"is_main"`
}

type SimpananResponse struct {
	Id            int    `json:"id"`
	JenisSimpanan string `json:"jenis_simpanan"`
	Nama          string `json:"nama"`
	Warna         string `json:"warna"`
	IsMain        bool   `json:"is_main"`
}

type SimpananCreateRequest struct {
	Jenis  string `json:"jenis"`
	Nama   string `json:"nama"`
	Warna  string `json:"warna"`
	IsMain bool   `json:"is_main"`
}
