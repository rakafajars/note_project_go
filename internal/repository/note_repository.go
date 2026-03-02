package repository

import (
	"notes-project/internal/models"

	"gorm.io/gorm"
)

type NoteRepository interface {
	Create(note *models.Note) error
	GetAll(query string, limit, offset int) ([]models.Note, int64, error)
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

func (r *noteRepository) GetAll(query string, limit, offset int) ([]models.Note, int64, error) {
	var notes []models.Note
	var total int64

	db := r.db.Model(&models.Note{})

	if query != "" {
		db = db.Where("title ILIKE ?", "%"+query+"%")
	}

	db.Count(&total)

	err := db.Limit(limit).Offset(offset).Find(&notes).Error
	return notes, total, err
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
