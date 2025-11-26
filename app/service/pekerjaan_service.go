package service

import (
	"crud_alumni/app/model"
	"crud_alumni/app/repository"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetAllPekerjaan(c *fiber.Ctx) error {
	data, err := repository.GetAllPekerjaan()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func GetPekerjaanByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	data, err := repository.GetPekerjaanByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID, _ := strconv.Atoi(c.Params("alumni_id"))
	data, err := repository.GetPekerjaanByAlumniID(alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func CreatePekerjaan(c *fiber.Ctx) error {
	var p model.Pekerjaan
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}
	if p.StatusPekerjaan == "" {
		p.StatusPekerjaan = "aktif"
	}
	if p.TanggalMulaiKerja == "" {
		p.TanggalMulaiKerja = time.Now().Format("2006-01-02")
	}
	id, err := repository.CreatePekerjaan(p)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	p.ID = id
	return c.Status(201).JSON(fiber.Map{"success": true, "data": p})
}

func UpdatePekerjaan(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var p model.Pekerjaan
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}
	if err := repository.UpdatePekerjaan(id, p); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func DeletePekerjaan(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := repository.DeletePekerjaan(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}
func GetPekerjaanByTahun(c *fiber.Ctx) error {
    tahun, _ := strconv.Atoi(c.Params("tahun"))
    data, err := repository.GetPekerjaanByTahun(tahun)
    if err != nil {
        return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"tahun": tahun, "jumlah": 0}})
    }
    return c.JSON(fiber.Map{"success": true, "data": data})
}

//pagination, sorting, searching
func GetPekerjaanPagination(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    sortBy := c.Query("sortBy", "id")
    order := c.Query("order", "asc")
    search := c.Query("search", "")

    offset := (page - 1) * limit
    whitelist := map[string]bool{"id": true, "nama_perusahaan": true, "posisi_jabatan": true, "bidang_industri": true, "lokasi_kerja": true, "status_pekerjaan": true}
    if !whitelist[sortBy] {
        sortBy = "id"
    }
    if strings.ToLower(order) != "desc" {
        order = "asc"
    }

    pekerjaan, err := repository.GetPekerjaanWithPagination(search, sortBy, order, limit, offset)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    total, _ := repository.CountPekerjaan(search)

    return c.JSON(fiber.Map{
        "data": pekerjaan,
        "meta": model.MetaInfo{
            Page:   page,
            Limit:  limit,
            Total:  total,
            Pages:  (total + limit - 1) / limit,
            SortBy: sortBy,
            Order:  order,
            Search: search,
        },
    })
}
func SoftDeletePekerjaan(c *fiber.Ctx) error {
    // ambil id pekerjaan dari param
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID pekerjaan tidak valid"})
    }

    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(int) // dari JWT

    // cek owner pekerjaan
    _, ownerUserID, err := repository.GetPekerjaanOwnerWithUser(id)
    if err == sql.ErrNoRows {
        return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
    }
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal cek pemilik"})
    }

    // kalau bukan admin → harus pemilik
    if role != "admin" && ownerUserID != userID {
        return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh hapus pekerjaan milik alumni lain"})
    }

    // lakukan soft delete
    rows, err := repository.SoftDeletePekerjaan(id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal hapus"})
    }
    if rows == 0 {
        return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan atau sudah dihapus"})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Pekerjaan berhasil dihapus (soft delete)",
    })
}
func RestorePekerjaan(c *fiber.Ctx) error {
    // ambil id pekerjaan dari param
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID pekerjaan tidak valid"})
    }

    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(int) // dari JWT

    // cek owner pekerjaan
    _, ownerUserID, err := repository.GetPekerjaanOwner(id)
    if err == sql.ErrNoRows {
        return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
    }
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal cek pemilik"})
    }

    // kalau bukan admin → harus pemilik
    if role != "admin" && ownerUserID != userID {
        return c.Status(403).JSON(fiber.Map{"error": "Tidak boleh restore pekerjaan milik alumni lain"})
    }

    // lakukan restore
    rows, err := repository.RestorePekerjaan(id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal Restore"})
    }
    if rows == 0 {
        return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan atau sudah di restore"})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Pekerjaan berhasil direstore",
    })
}

func Trash(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int) // dari JWT

	var data []model.Pekerjaan
	var err error

	if role == "admin" {
		// 🔹 Admin bisa lihat semua data yang isdellete='yes'
		data, err = repository.TrashAll()
	} else {
		// 🔹 User hanya bisa lihat data miliknya
		data, err = repository.TrashByUser(userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(data),
		"data":    data,
	})
}

func HardDeletePekerjaan(c *fiber.Ctx) error {
    // ambil id pekerjaan dari param
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID pekerjaan tidak valid"})
    }

    // ambil role & userID dari JWT (middleware sudah isi ini)
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(int)

    // cek apakah pekerjaan ada & siapa pemiliknya
    _, ownerUserID, err := repository.GetPekerjaanOwnerr(id)
    if err == sql.ErrNoRows {
        return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
    }
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal memeriksa pemilik pekerjaan"})
    }

    // jika bukan admin, maka harus pemilik datanya
    if role != "admin" && ownerUserID != userID {
        return c.Status(403).JSON(fiber.Map{"error": "Tidak diizinkan menghapus pekerjaan milik orang lain"})
    }

    // lakukan hard delete di database
    err = repository.HardDeletePekerjaan(id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Pekerjaan berhasil dihapus permanen",
    })
}



