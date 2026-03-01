package repository

import (
	"notes-project/internal/models"

	"gorm.io/gorm"
)

type NoteRepository interface {
	Create(note *models.Note) error
	GetAllNote() ([]models.Note, error)
}

type noteRepository struct {
	db *gorm.DB
}

// function untuk menginisialisasi repository
func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &noteRepository{db}
}

func (r *noteRepository) Create(note *models.Note) error {
	return r.db.Create(note).Error
}

func (r *noteRepository) GetAllNote() ([]models.Note, error) {
	var notes []models.Note
	err := r.db.Find(&notes).Error
	return notes, err
}
