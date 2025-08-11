Tokopedia Handphone Scraper
Deskripsi
Program ini adalah utilitas berbasis Go untuk mengekstrak 100 produk teratas dari kategori Mobile Phones / Handphone di Tokopedia, lalu menyimpannya ke file CSV.

Informasi produk yang dikumpulkan:

Name of Product

Description

Image Link

Price

Rating (out of 5 stars)

Name of store or merchant

Scraper ini menggunakan Chromedp untuk meng-handle konten dinamis Tokopedia dan GoQuery untuk parsing HTML.

Persyaratan Sistem
Go 1.18 atau lebih baru

Google Chrome terinstal

Koneksi internet stabil

Instalasi
Clone repository:

bash
Copy
Edit
git clone https://github.com/username/tokopedia-handphone-scraper.git
cd tokopedia-handphone-scraper
Install dependensi:

bash
Copy
Edit
go mod tidy
Cara Menjalankan
bash
Copy
Edit
go run main.go
Program akan:

Membuka halaman Tokopedia kategori Handphone

Menjalankan auto-scroll dan navigasi multi-halaman

Menyimpan data ke file tokopedia_handphone.csv

Struktur Output CSV
File tokopedia_handphone.csv memiliki format kolom berikut:

pgsql
Copy
Edit
Name,Description,ImageLink,Price,Rating,Merchant
Contoh Output CSV
arduino
Copy
Edit
Name,Description,ImageLink,Price,Rating,Merchant
Samsung Galaxy A14,,https://example.com/image1.jpg,Rp 2.999.000,4.8,Samsung Official Store
Infinix Note 12,,https://example.com/image2.jpg,Rp 2.499.000,4.7,Infinix Official Store
OPPO A57,,https://example.com/image3.jpg,Rp 2.399.000,4.9,OPPO Official Store
Xiaomi Redmi Note 11,,https://example.com/image4.jpg,Rp 2.699.000,4.8,Xiaomi Official Store
Vivo Y16,,https://example.com/image5.jpg,Rp 1.799.000,4.7,Vivo Official Store
Catatan: Kolom "Description" bisa kosong karena tidak semua listing menampilkan deskripsi singkat di halaman kategori.

Catatan Teknis
Dynamic Content: Tokopedia memuat produk dengan JavaScript, sehingga scraping dilakukan setelah halaman selesai di-render oleh Chrome Headless.

Selector: Selector CSS dapat berubah sewaktu-waktu, sehingga perlu diperbarui jika scraping gagal mendapatkan data.

Rate Limiting: Untuk menghindari pemblokiran, script memberi jeda beberapa detik di antara scroll dan pergantian halaman.

Description: Beberapa produk tidak menampilkan deskripsi singkat di listing, sehingga kolom ini bisa kosong.

Lisensi
MIT License
