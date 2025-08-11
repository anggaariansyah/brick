# Handphone Scraper

## Deskripsi  
Program ini adalah utilitas berbasis **Go** untuk mengekstrak **100 produk teratas** dari kategori **Mobile Phones / Handphone** di Tokopedia, lalu menyimpannya ke file CSV.  

Informasi produk yang dikumpulkan:  
1. **Name of Product**  
2. **Description**  
3. **Image Link**  
4. **Price**  
5. **Rating** (out of 5 stars)  
6. **Name of store or merchant**  

Scraper ini menggunakan **Chromedp** untuk meng-handle konten dinamis Tokopedia dan **GoQuery** untuk parsing HTML.  

---

## Persyaratan Sistem
- Go 1.18 atau lebih baru  
- Google Chrome terinstal  
- Koneksi internet stabil  

---

## Instalasi
Clone repository:
```bash
git clone https://github.com/anggaariansyah/brick.git
cd brick
go mod tidy
go run main.go

