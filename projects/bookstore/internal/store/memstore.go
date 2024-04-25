package store

import (
	"sync"
	"yuxing.code/bookstore/store"
	"yuxing.code/bookstore/store/factory"
)

var MemStoreProviderName = "mem"

// 利用 init函数 注册存储实现
func init() {
	factory.Register(MemStoreProviderName, &MemStore{
		books: make(map[string]*store.Book),
	})
}

// MemStore 基于内存的存储实现
type MemStore struct {
	sync.RWMutex
	books map[string]*store.Book
}

func (m *MemStore) Get(id string) (store.Book, error) {
	m.RLock()
	defer m.RUnlock()
	book, exit := m.books[id]
	if !exit {
		return store.Book{}, store.ErrNotFound
	}
	return *book, nil
}

func (m *MemStore) GetAll() ([]store.Book, error) {
	m.RLock()
	defer m.RUnlock()

	books := make([]store.Book, len(m.books))
	for _, book := range m.books {
		books = append(books, *book)
	}
	return books, nil
}

func (m *MemStore) Delete(id string) error {
	m.Lock()
	defer m.Unlock()

	if _, exit := m.books[id]; !exit {
		return store.ErrNotFound
	}
	delete(m.books, id)
	return nil
}

func (m *MemStore) Create(book *store.Book) error {
	m.Lock()
	defer m.Unlock()

	if _, exit := m.books[book.Id]; exit {
		return store.ErrExist
	}
	nBook := *book
	m.books[book.Id] = &nBook
	return nil
}

func (m *MemStore) Update(book *store.Book) error {
	m.Lock()
	defer m.Unlock()

	oldBook, exit := m.books[book.Id]
	if !exit {
		return store.ErrNotFound
	}

	nBook := *oldBook
	if book.Name != "" {
		nBook.Name = book.Name
	}

	if book.Authors != nil {
		nBook.Authors = book.Authors
	}

	if book.Press != "" {
		nBook.Press = book.Press
	}
	m.books[book.Id] = &nBook
	return nil
}
