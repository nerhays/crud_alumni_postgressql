package repository

import (
	"crud_alumni/app/model"
	"crud_alumni/database"
	"fmt"
	"time"
)

// Ambil semua alumni
func GetAllAlumni() ([]model.Alumni, error) {
	rows, err := database.DB.Query(`SELECT id,user_id,nim,nama,jurusan,angkatan,tahun_lulus,email,no_telepon,alamat,created_at,updated_at FROM alumni`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Alumni
	for rows.Next() {
		var a model.Alumni
		rows.Scan(&a.ID,&a.UserID,&a.NIM,&a.Nama,&a.Jurusan,&a.Angkatan,&a.TahunLulus,&a.Email,&a.NoTelepon,&a.Alamat,&a.CreatedAt,&a.UpdatedAt)
		list = append(list,a)
	}
	return list,nil
}

// Tambah alumni
func CreateAlumni(a model.Alumni) (int,error) {
	var id int
	err := database.DB.QueryRow(`
		INSERT INTO alumni (nim,nama,jurusan,angkatan,tahun_lulus,email,no_telepon,alamat,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		a.NIM,a.Nama,a.Jurusan,a.Angkatan,a.TahunLulus,a.Email,a.NoTelepon,a.Alamat,time.Now(),time.Now()).Scan(&id)
	return id,err
}

// Update alumni
func UpdateAlumni(id int,a model.Alumni) error {
	_,err := database.DB.Exec(`
		UPDATE alumni SET nama=$1,jurusan=$2,angkatan=$3,tahun_lulus=$4,email=$5,no_telepon=$6,alamat=$7,updated_at=$8 WHERE id=$9`,
		a.Nama,a.Jurusan,a.Angkatan,a.TahunLulus,a.Email,a.NoTelepon,a.Alamat,time.Now(),id)
	return err
}

// Hapus alumni
func DeleteAlumni(id int) error {
	_,err := database.DB.Exec("DELETE FROM alumni WHERE id=$1",id)
	return err
}

//pagination sorting searching
func GetAlumniWithPagination(search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
    query := fmt.Sprintf(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
        FROM alumni
        WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1 OR email ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

    rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var alumni []model.Alumni
    for rows.Next() {
        var a model.Alumni
        if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt); err != nil {
            return nil, err
        }
        alumni = append(alumni, a)
    }
    return alumni, nil
}

func CountAlumni(search string) (int, error) {
    var total int
    err := database.DB.QueryRow(`SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1 OR email ILIKE $1`, "%"+search+"%").Scan(&total)
    return total, err
}
