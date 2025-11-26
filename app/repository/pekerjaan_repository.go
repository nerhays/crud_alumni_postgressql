package repository

import (
	"crud_alumni/app/model"
	"crud_alumni/database"
	"fmt"
	"time"
)

// Ambil semua pekerjaan
func GetAllPekerjaan() ([]model.Pekerjaan, error) {
	rows, err := database.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, isdellete
		FROM pekerjaan_alumni
		ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		var tMulai time.Time
		var tSelesai *time.Time

		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &tMulai, &tSelesai,
			&p.StatusPekerjaan, &p.Deskripsi, &p.IsDellete,
		); err != nil {
			return nil, err
		}

		p.TanggalMulaiKerja = tMulai.Format("2006-01-02")
		if tSelesai != nil {
			val := tSelesai.Format("2006-01-02")
			p.TanggalSelesaiKerja = &val
		}

		result = append(result, p)
	}
	return result, nil
}

// Ambil pekerjaan by ID
func GetPekerjaanByID(id int) (*model.Pekerjaan, error) {
	var p model.Pekerjaan
	var tMulai time.Time
	var tSelesai *time.Time

	err := database.DB.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, isdellete
		FROM pekerjaan_alumni WHERE id = $1`, id).
		Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &tMulai, &tSelesai,
			&p.StatusPekerjaan, &p.Deskripsi, &p.IsDellete)
	if err != nil {
		return nil, err
	}

	p.TanggalMulaiKerja = tMulai.Format("2006-01-02")
	if tSelesai != nil {
		val := tSelesai.Format("2006-01-02")
		p.TanggalSelesaiKerja = &val
	}

	return &p, nil
}

// Ambil pekerjaan by AlumniID
func GetPekerjaanByAlumniID(alumniID int) ([]model.Pekerjaan, error) {
	rows, err := database.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, isdellete
		FROM pekerjaan_alumni WHERE alumni_id = $1`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		var tMulai time.Time
		var tSelesai *time.Time

		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &tMulai, &tSelesai,
			&p.StatusPekerjaan, &p.Deskripsi, &p.IsDellete,
		); err != nil {
			return nil, err
		}

		p.TanggalMulaiKerja = tMulai.Format("2006-01-02")
		if tSelesai != nil {
			val := tSelesai.Format("2006-01-02")
			p.TanggalSelesaiKerja = &val
		}

		result = append(result, p)
	}
	return result, nil
}

// Tambah pekerjaan
func CreatePekerjaan(p model.Pekerjaan) (int, error) {
	tMulai, err := time.Parse("2006-01-02", p.TanggalMulaiKerja)
	if err != nil {
		return 0, err
	}

	var tSelesai *time.Time
	if p.TanggalSelesaiKerja != nil && *p.TanggalSelesaiKerja != "" {
		parsed, err := time.Parse("2006-01-02", *p.TanggalSelesaiKerja)
		if err != nil {
			return 0, err
		}
		tSelesai = &parsed
	}

	var lastInsertID int
	err = database.DB.QueryRow(`
		INSERT INTO pekerjaan_alumni (
			alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
			lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
			status_pekerjaan, deskripsi_pekerjaan
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id
	`,
		p.AlumniID, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri,
		p.LokasiKerja, p.GajiRange, tMulai, tSelesai,
		p.StatusPekerjaan, p.Deskripsi,
	).Scan(&lastInsertID)

	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

// Update pekerjaan
func UpdatePekerjaan(id int, p model.Pekerjaan) error {
	tMulai, err := time.Parse("2006-01-02", p.TanggalMulaiKerja)
	if err != nil {
		return err
	}

	var tSelesai *time.Time
	if p.TanggalSelesaiKerja != nil && *p.TanggalSelesaiKerja != "" {
		parsed, err := time.Parse("2006-01-02", *p.TanggalSelesaiKerja)
		if err != nil {
			return err
		}
		tSelesai = &parsed
	}

	_, err = database.DB.Exec(`
		UPDATE pekerjaan_alumni
		SET alumni_id=$1, nama_perusahaan=$2, posisi_jabatan=$3, bidang_industri=$4,
		    lokasi_kerja=$5, gaji_range=$6, tanggal_mulai_kerja=$7, tanggal_selesai_kerja=$8,
		    status_pekerjaan=$9, deskripsi_pekerjaan=$10, updated_at=NOW()
		WHERE id=$11
	`,
		p.AlumniID, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri,
		p.LokasiKerja, p.GajiRange, tMulai, tSelesai,
		p.StatusPekerjaan, p.Deskripsi, id,
	)
	return err
}

// Hapus pekerjaan
func DeletePekerjaan(id int) error {
	_, err := database.DB.Exec(`DELETE FROM pekerjaan_alumni WHERE id=$1`, id)
	return err
}

func GetPekerjaanByTahun(tahun int) (model.JumlahPekerjaanPerTahun, error) {
    var result model.JumlahPekerjaanPerTahun
    err := database.DB.QueryRow(`
        SELECT EXTRACT(YEAR FROM tanggal_mulai_kerja)::int AS tahun, COUNT(*)
        FROM pekerjaan_alumni
        WHERE EXTRACT(YEAR FROM tanggal_mulai_kerja) = $1
        GROUP BY tahun
    `, tahun).Scan(&result.Tahun, &result.Jumlah)
    if err != nil {
        result.Tahun = tahun
        result.Jumlah = 0
        return result, err
    }
    return result, nil
}

//pagination sorting searching
func GetPekerjaanWithPagination(search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
    query := fmt.Sprintf(`
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan
        FROM pekerjaan_alumni
        WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1 OR lokasi_kerja ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

    rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var pekerjaan []model.Pekerjaan
    for rows.Next() {
        var p model.Pekerjaan
        if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.Deskripsi); err != nil {
            return nil, err
        }
        pekerjaan = append(pekerjaan, p)
    }
    return pekerjaan, nil
}

func CountPekerjaan(search string) (int, error) {
    var total int
    err := database.DB.QueryRow(`SELECT COUNT(*) FROM pekerjaan_alumni WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1 OR lokasi_kerja ILIKE $1`, "%"+search+"%").Scan(&total)
    return total, err
}

func GetPekerjaanOwnerWithUser(id int) (int, int, error) {
    var alumniID, ownerUserID int
    err := database.DB.QueryRow(`
        SELECT p.alumni_id, a.user_id
        FROM pekerjaan_alumni p
        JOIN alumni a ON p.alumni_id = a.id
        WHERE p.id = $1 AND p.isdellete = 'no'`,
        id,
    ).Scan(&alumniID, &ownerUserID)

    if err != nil {
        return 0, 0, err
    }
    return alumniID, ownerUserID, nil
}

func SoftDeletePekerjaan(id int) (int64, error) {
    res, err := database.DB.Exec(`
        UPDATE pekerjaan_alumni
        SET isdellete = 'yes', updated_at = NOW()
        WHERE id = $1 AND isdellete = 'no'`,
        id,
    )
    if err != nil {
        return 0, err
    }
    return res.RowsAffected()
}
func RestorePekerjaan(id int) (int64, error) {
    res, err := database.DB.Exec(`
        UPDATE pekerjaan_alumni
        SET isdellete = 'no', updated_at = NOW()
        WHERE id = $1 AND isdellete = 'yes'`,
        id,
    )
    if err != nil {
        return 0, err
    }
    return res.RowsAffected()
}
func GetPekerjaanOwner(id int) (int, int, error) {
    var alumniID, ownerUserID int
    err := database.DB.QueryRow(`
        SELECT p.alumni_id, a.user_id
        FROM pekerjaan_alumni p
        JOIN alumni a ON p.alumni_id = a.id
        WHERE p.id = $1 AND p.isdellete = 'yes'`,
        id,
    ).Scan(&alumniID, &ownerUserID)

    if err != nil {
        return 0, 0, err
    }
    return alumniID, ownerUserID, nil
}

func TrashAll() ([]model.Pekerjaan, error) {
	rows, err := database.DB.Query(`
		SELECT p.id, p.alumni_id, p.nama_perusahaan, p.posisi_jabatan, p.bidang_industri,
		       p.lokasi_kerja, p.gaji_range, p.tanggal_mulai_kerja, p.tanggal_selesai_kerja,
		       p.status_pekerjaan, p.deskripsi_pekerjaan, p.isdellete
		FROM pekerjaan_alumni p
		WHERE p.isdellete = 'yes'
		ORDER BY p.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		var tMulai time.Time
		var tSelesai *time.Time

		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &tMulai, &tSelesai,
			&p.StatusPekerjaan, &p.Deskripsi, &p.IsDellete,
		); err != nil {
			return nil, err
		}

		p.TanggalMulaiKerja = tMulai.Format("2006-01-02")
		if tSelesai != nil {
			val := tSelesai.Format("2006-01-02")
			p.TanggalSelesaiKerja = &val
		}
		result = append(result, p)
	}
	return result, nil
}
func TrashByUser(userID int) ([]model.Pekerjaan, error) {
	rows, err := database.DB.Query(`
		SELECT p.id, p.alumni_id, p.nama_perusahaan, p.posisi_jabatan, p.bidang_industri,
		       p.lokasi_kerja, p.gaji_range, p.tanggal_mulai_kerja, p.tanggal_selesai_kerja,
		       p.status_pekerjaan, p.deskripsi_pekerjaan, p.isdellete
		FROM pekerjaan_alumni p
		JOIN alumni a ON p.alumni_id = a.id
		WHERE p.isdellete = 'yes' AND a.user_id = $1
		ORDER BY p.id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		var tMulai time.Time
		var tSelesai *time.Time

		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &tMulai, &tSelesai,
			&p.StatusPekerjaan, &p.Deskripsi, &p.IsDellete,
		); err != nil {
			return nil, err
		}

		p.TanggalMulaiKerja = tMulai.Format("2006-01-02")
		if tSelesai != nil {
			val := tSelesai.Format("2006-01-02")
			p.TanggalSelesaiKerja = &val
		}
		result = append(result, p)
	}
	return result, nil
}

func HardDeletePekerjaan(id int) error {
    query := `DELETE FROM pekerjaan_alumni WHERE id = $1`
    _, err := database.DB.Exec(query, id)
    return err
}
func GetPekerjaanOwnerr(id int) (int, int, error) {
    var alumniID, ownerUserID int
    err := database.DB.QueryRow(`
        SELECT p.alumni_id, a.user_id
        FROM pekerjaan_alumni p
        JOIN alumni a ON p.alumni_id = a.id
        WHERE p.id = $1
    `, id).Scan(&alumniID, &ownerUserID)

    if err != nil {
        return 0, 0, err
    }
    return alumniID, ownerUserID, nil
}




