package main

import (
	"belajar-go/database"
	"belajar-go/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"belajar-go/repositories"
	"belajar-go/services"
	"belajar-go/handlers"

	"github.com/spf13/viper"
)
var produk = []models.Product{
	{ID: 1, Name: "Produk A", Price: 10000, Stock: 50},
	{ID: 2, Name: "Produk B", Price: 20000, Stock: 30},
	{ID: 3, Name: "Produk C", Price: 15000, Stock: 20},
}


// Handler untuk health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status":  "OK",
		"message": "Cashier API is running",
	}
	json.NewEncoder(w).Encode(response)
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
	var produkBaru models.Product
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

	var produkUpdate models.Product
	json.NewDecoder(r.Body).Decode(&produkUpdate)

	for i, p := range produk {
		if p.ID == produkUpdate.ID {
			produk[i] = produkUpdate
			w.WriteHeader(http.StatusOK)

			fmt.Println("Produk ID", p.ID, "berhasil diupdate")
			json.NewEncoder(w).Encode(produkUpdate)
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusMethodNotAllowed)
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// Jalankan server
func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		fmt.Println("Gagal terhubung ke database:", err)
		return
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Routing
	http.HandleFunc("/api/produk", productHandler.HandleProducts)


	http.HandleFunc("/produk", ambilProduk)
	http.HandleFunc("/produk/", ambilDetailProduk)
	http.HandleFunc("/tambah", tambahProduk)
	http.HandleFunc("/hapus", hapusProduk)
	http.HandleFunc("/update", updateProduk)
	http.HandleFunc("/health", healthHandler)

	fmt.Println("Server kasir berjalan di localhost:" + config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Gagal menjalankan server")
	}
}
