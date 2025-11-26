package service

import (
	"crud_alumni/app/model"
	"crud_alumni/app/repository"
	"crud_alumni/database"
	"database/sql"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAllAlumni(c *fiber.Ctx) error {
	data,err := repository.GetAllAlumni()
	if err!=nil { return c.Status(500).JSON(fiber.Map{"error":"Gagal ambil data"}) }
	return c.JSON(fiber.Map{"success":true,"data":data})
}

func CreateAlumni(c *fiber.Ctx) error {
	var a model.Alumni
	if err := c.BodyParser(&a); err!=nil {
		return c.Status(400).JSON(fiber.Map{"error":"Body tidak valid"})
	}
	id,err := repository.CreateAlumni(a)
	if err!=nil { return c.Status(500).JSON(fiber.Map{"error":"Gagal tambah"}) }
	a.ID = id
	return c.Status(201).JSON(fiber.Map{"success":true,"data":a})
}

func UpdateAlumni(c *fiber.Ctx) error {
	id,_ := strconv.Atoi(c.Params("id"))
	var a model.Alumni
	if err := c.BodyParser(&a); err!=nil {
		return c.Status(400).JSON(fiber.Map{"error":"Body tidak valid"})
	}
	if err := repository.UpdateAlumni(id,a); err!=nil {
		return c.Status(500).JSON(fiber.Map{"error":"Gagal update"})
	}
	return c.JSON(fiber.Map{"success":true})
}

func DeleteAlumni(c *fiber.Ctx) error {
	id,_ := strconv.Atoi(c.Params("id"))
	if err := repository.DeleteAlumni(id); err!=nil {
		return c.Status(500).JSON(fiber.Map{"error":"Gagal hapus"})
	}
	return c.JSON(fiber.Map{"success":true})
}

func GetAlumniByID(c *fiber.Ctx) error {
    id := c.Params("id")

    var a model.Alumni
    err := database.DB.QueryRow(`
        SELECT id, user_id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat
        FROM alumni WHERE id = $1`, id).
        Scan(&a.ID, &a.UserID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat)

    if err == sql.ErrNoRows {
        return c.Status(404).JSON(fiber.Map{
            "success": false,
            "message": "Alumni tidak ditemukan",
        })
    } else if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "error": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    a,
    })
}

//pagination sorting searching
func GetAlumniPagination(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    sortBy := c.Query("sortBy", "id")
    order := c.Query("order", "asc")
    search := c.Query("search", "")

    offset := (page - 1) * limit
    whitelist := map[string]bool{"id": true, "nim": true, "nama": true, "angkatan": true, "tahun_lulus": true, "email": true}
    if !whitelist[sortBy] {
        sortBy = "id"
    }
    if strings.ToLower(order) != "desc" {
        order = "asc"
    }

    alumni, err := repository.GetAlumniWithPagination(search, sortBy, order, limit, offset)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    total, _ := repository.CountAlumni(search)

    return c.JSON(model.AlumniResponse{
        Data: alumni,
        Meta: model.MetaInfo{
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
