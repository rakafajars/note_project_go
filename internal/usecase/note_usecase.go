package usecase

import (
	"errors"
	"notes-project/internal/models"
	"notes-project/internal/repository"
	"time"

	"gorm.io/gorm"
)

// interface untuk kontrak logika bisnis
type NoteUsecase interface {
	GetAllNotes() ([]models.Note, error)
	CreateNote(title, content string) (*models.Note, error)
	DeleteNote(id uint) error
}

type noteUsecase struct {
	repo repository.NoteRepository
}

// constructor untuk menyutikann (inject) repository ke usecase
func NewTodoUsecase(r repository.NoteRepository) NoteUsecase {
	return &noteUsecase{repo: r}
}

func (u *noteUsecase) CreateNote(title, content string) (*models.Note, error) {
	// contoh logika Bisnis: judul tidak boleh kosong
	if title == "" {
		return nil, errors.New("Judul tidak boleh kosong")
	}

	if content == "" {
		return nil, errors.New("Content tidak boleh kosong")
	}

	note := &models.Note{
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}

	// Memanggil repository untuk simpan ke DB
	err := u.repo.Create(note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (u *noteUsecase) GetAllNotes() ([]models.Note, error) {
	return u.repo.GetAllNote()
}

func (u *noteUsecase) DeleteNote(id uint) error {
	// Validasi Sederhana: Pastikan ID tidak nol
	if id == 0 {
		return errors.New("ID catatan tidak valid")
	}

	err := u.repo.Delete(id)
	if err != nil {
		// jika errornya adalah recordnotFound, kita berikan pesan spesifik
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("catatan tidak ditemukan")
		}

		return err
	}

	return nil
}
