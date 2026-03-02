package repository

import (
	"notes-project/internal/models"

	"gorm.io/gorm"
)

type NoteRepository interface {
	Create(note *models.Note) error
	GetAllNote() ([]models.Note, error)
	Delete(id uint) error
	Update(id uint, note *models.Note) error
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

func (r *noteRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Note{}, id)

	// Cek Jika terjadi error koneksi/query
	if result.Error != nil {
		return result.Error
	}

	// Jika tidak ada baris yang terhapus (ID tidak ketemu)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Error bawaan GORM untuk "Data Tidak Ditemukan"
	}

	return nil
}

func (r *noteRepository) Update(id uint, note *models.Note) error {
	// mencari data berdasarkan ID, lalu diperbaiki fieldnya
	result := r.db.Model(&models.Note{}).Where("id = ?", id).Updates(note)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
