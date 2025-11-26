package repository

import (
	"crud_alumni/app/model"
	"crud_alumni/database"
)

func FindUserByUsernameOrEmail(identifier string) (*model.User, string, error) {
    var user model.User
    var passwordHash string

    err := database.DB.QueryRow(`
        SELECT id, username, email, password_hash, role, created_at
        FROM users WHERE username=$1 OR email=$1
    `, identifier).Scan(
        &user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.CreatedAt,
    )

    if err != nil {
        return nil, "", err
    }

    return &user, passwordHash, nil
}

