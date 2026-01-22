package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Definisi data produk
type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

// Inisialisasi data produk
var produk = []Produk{
	{ID: 1, Nama: "Beras Merah 5kg", Harga: 65000, Stok: 25},
	{ID: 2, Nama: "Gula Pasir 1kg", Harga: 15000, Stok: 20},
	{ID: 3, Nama: "Minyak Goreng 1L", Harga: 17000, Stok: 30},
	{ID: 4, Nama: "Telur Ayam 1kg", Harga: 32000, Stok: 25},
	{ID: 5, Nama: "Tepung Terigu 1kg", Harga: 12000, Stok: 20},
}

// Fungsi untuk mengambil data produk
func ambilProduk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Hanya GET yang diizinkan", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produk)
}

// Fungsi untuk mengambil data detail produk berdasarkan ID
func ambilDetailProduk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Hanya GET yang diizinkan", http.StatusMethodNotAllowed)
		return
	}

	idString := strings.TrimPrefix(r.URL.Path, "/produk/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusMethodNotAllowed)
}

// Fungsi untuk menambahkan data produk
func tambahProduk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Hanya POST yang diizinkan", http.StatusMethodNotAllowed)
		return
	}
	var produkBaru Produk
	if err := json.NewDecoder(r.Body).Decode(&produkBaru); err != nil {
		http.Error(w, "Body request tidak valid", http.StatusBadRequest)
		return
	}

	produk = append(produk, produkBaru)
	fmt.Println("Produk berhasil ditambahkan")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(produkBaru)
}

// Fungsi untuk menghapus data produk
func hapusProduk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Hanya DELETE yang diizinkan", http.StatusMethodNotAllowed)
		return
	}

	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)

			fmt.Println("Produk berhasil dihapus")
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusMethodNotAllowed)
}

// Fungsi untuk mengupdate data produk
func updateProduk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Hanya PUT yang diizinkan", http.StatusMethodNotAllowed)
		return
	}

	var produkUpdate Produk
	json.NewDecoder(r.Body).Decode(&produkUpdate)

	for i, p := range produk {
		if p.ID == produkUpdate.ID {
			produk[i] = produkUpdate
			w.WriteHeader(http.StatusOK)

			fmt.Println("Produk ID %d berhasil diupdate", p.ID)
			json.NewEncoder(w).Encode(produkUpdate)
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusMethodNotAllowed)
}

// Jalankan server
func main() {
	http.HandleFunc("/produk", ambilProduk)
	http.HandleFunc("/produk/", ambilDetailProduk)
	http.HandleFunc("/tambah", tambahProduk)
	http.HandleFunc("/hapus", hapusProduk)
	http.HandleFunc("/update", updateProduk)

	fmt.Println("Server kasir berjalan di port:8080")
	http.ListenAndServe(":8080", nil)
}
