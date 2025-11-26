package model

type Pekerjaan struct {
	ID                  int     `json:"id"`
	AlumniID            int     `json:"alumni_id"`
	NamaPerusahaan      string  `json:"nama_perusahaan"`
	PosisiJabatan       string  `json:"posisi_jabatan"`
	BidangIndustri      string  `json:"bidang_industri"`
	LokasiKerja         string  `json:"lokasi_kerja"`
	GajiRange           string  `json:"gaji_range,omitempty"`
	TanggalMulaiKerja   string  `json:"tanggal_mulai_kerja"`             // YYYY-MM-DD
	TanggalSelesaiKerja *string `json:"tanggal_selesai_kerja,omitempty"` // YYYY-MM-DD
	StatusPekerjaan     string  `json:"status_pekerjaan"`
	IsDellete           string  `json:"isdellete"`
	Deskripsi           string  `json:"deskripsi_pekerjaan,omitempty"`
}
type isdell struct {
	IsDellete string `json:"isdellete"`
}
type JumlahPekerjaanPerTahun struct {
	Tahun  int `json:"tahun"`
	Jumlah int `json:"jumlah"`
}

type PekerjaanResponse struct {
	Data []Pekerjaan `json:"data"`
	Meta MetaInfo    `json:"meta"`
}
