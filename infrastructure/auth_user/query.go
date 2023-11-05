package auth_user

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"lambda_list/infrastructure/db"
)

func FindByStudentId(studentId uuid.UUID) (uuid.UUID, error) {
	conn, err := db.ConnectRDB()
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Close()
	dbStudentId := uuid.Nil
	rawStudentId, err := studentId.MarshalBinary()
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "fail marshal student id")
	}
	err = conn.QueryRow("SELECT student_id FROM students WHERE student_id = ?", rawStudentId).Scan(&dbStudentId)
	if err != nil || dbStudentId == uuid.Nil {
		return uuid.Nil, errors.Wrap(err, "fail find student id")
	}
	return dbStudentId, nil
}
