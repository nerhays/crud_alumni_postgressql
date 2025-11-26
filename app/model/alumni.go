package model

import "time"

	type Alumni struct {
		ID         int       `json:"id"`
		UserID         int       `json:"userid"`
		NIM        string    `json:"nim"`
		Nama       string    `json:"nama"`
		Jurusan    string    `json:"jurusan"`
		Angkatan   int       `json:"angkatan"`
		TahunLulus int       `json:"tahun_lulus"`
		Email      string    `json:"email"`
		NoTelepon  string    `json:"no_telepon,omitempty"`
		Alamat     string    `json:"alamat,omitempty"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	type MetaInfo struct {
	    Page   int    `json:"page"`
	    Limit  int    `json:"limit"`
	    Total  int    `json:"total"`
	    Pages  int    `json:"pages"`
	    SortBy string `json:"sortBy"`
	    Order  string `json:"order"`
	    Search string `json:"search"`
	}

	type AlumniResponse struct {
	    Data []Alumni `json:"data"`
	    Meta MetaInfo `json:"meta"`
	}
