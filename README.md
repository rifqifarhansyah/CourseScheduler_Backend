# CourseScheduler - BackEnd
<h2 align="center">
   Course Scheduler Web
</h2>
<hr>

## Table of Contents
1. [General Info](#general-information)
2. [Creator Info](#creator-information)
3. [Features](#features)
4. [Technologies Used](#technologies-used)
5. [Setup](#setup)
6. [Usage](#usage)
7. [Algorithm](#usage)
8. [Video Capture](#videocapture)
9. [Screenshots](#screenshots)
10. [Structure](#structure)
11. [Project Status](#project-status)
12. [Room for Improvement](#room-for-improvement)
13. [Acknowledgements](#acknowledgements)
14. [Contact](#contact)

<a name="general-information"></a>

## General Information
Sebuah aplikasi berbasis website sederhana yang dapat digunakan untuk melakukan `pencarian data mata kuliah` agar memperoleh nilai maksimal dengan memanfaatkan `algoritma dynamic programming`, `menambah data mata kuliah` via `text field dan JSON file`, `menampilkan data mata kuliah` dan `melakukan penghapusan data mata kuliah`. Website ini dibangun dengan `React` sebagai `Framework FrontEnd`, `Golang` sebagai `Framework BackEnd`, dan `MySQL` sebagai database. Repository Backend dan Frontend terpisah satu sama lain. Tugas ini digunakan untuk memenuhi `Tugas 8 Seleksi Lab IRK tahun 2023`.
 
<a name="creator-information"></a>

## Creator Information

| Nama                        | NIM      | E-Mail                      |
| --------------------------- | -------- | --------------------------- |
| Mohammad Rifqi Farhansyah   | 13521166 | 13521166@std.stei.itb.ac.id |

<a name="features"></a>

## Features
- `Mencari data mata kuliah` pada `SearchPage` agar mendapatkan nilai maksimal `dengan algoritma dynamic programming`
- `Menambahkan data mata kuliah` dengan `TextField atau JSON File` pada `Add Page`
- `Melihat data mata kuliah` dan `menghapus data mata kuliah` pada `Data Page`

<a name="technologies-used"></a>

## Technologies Used
* [Golang](https://go.dev/learn/) - versi 1.20.0
* [React](https://react.dev/) - versi 18.2.0
* [MySQL](https://www.mysql.com/) - versi 5.7.30
* [Docker](https://www.docker.com/) - versi 20.10.21

> Note: The version of the libraries above is the version that we used in this project. You can use the latest version of the libraries.

<a name="setup"></a>

## Setup
1. Clone Repository ini dengan menggunakan command berikut
   ```sh
   git clone https://github.com/rifqifarhansyah/CourseScheduler_Backend.git
   ```
2. Buka Folder "coursescheduler_backend" di Terminal
3. Jalankan script docker dengan menjalankan
    ```sh
    docker compose build
    ```
4. Lalu jalankan program dengan perintah
    ```sh
    docker compose up
    ```
5. Apabila pada endpoint tertentu muncul error `table coursescheduler.courses` maka jalankan perintah pada langkah 5-8 berikut:
    ```sh
    docker-compose exec db mysql -u root -p
    ```
6. Masukkan password sesuai yang diminta
7. Setelah berhasil masuk ke `MySQL server`, masukkan query berikut:
    ```sh
    CREATE DATABASE IF NOT EXISTS course_scheduler;
    USE course_scheduler;
    CREATE TABLE IF NOT EXISTS courses (
    id INT PRIMARY KEY AUTO_INCREMENT,
    namaMk VARCHAR(255),
    jurusan VARCHAR(255),
    fakultas VARCHAR(255),
    jumlahSks INT,
    semesterMin INT,
    prediksiNilai VARCHAR(255)
    );
    ```
8. Lalu keluarlah dari `MySQL server` dengan memasukkan command:
    ```sh
    exit
    ```
9. Program backend akan otomatis dijalankan pada localhost `(default PORT:5001)`

<a name="usage"></a>

## Usage
1. Untuk `mencari mata kuliah yang menghasilkan nilai maksimal`, masuk ke `Search Page` lalu isi seluruh `text field yang diminta`, kemudian tekan tombol `search`
2. Untuk `menambahkan mata kuliah`, masuk ke `Add Page` lalu pilih metode yang hendak digunakan, apabila menggunakan metode `text field`, isi seluruh `text field yang diminta`, kemudian tekan tombol `add`
3. Apabila `penambahan mata kuliah` dilakukan via `JSON File`, maka terlebih dahulu cek `JSON Example`, sesuaikan format JSON masukan, lalu tekan `Choose File` untuk memilih file JSON yang hendak digunakan
4. Untuk `melihat data mata kuliah yang ada`, masuklah ke `Data Page`
5. Sementara itu, `penghapusan data mata kuliah` dapat dilakukan dengan menekan `icon kotak sampah` di samping tiap data mata kuliah pada `data page`

<a name="algorithm"></a>

## Algorithm
### Deskripsi Umum Algoritma Dynamic Programming
Algoritma `Dynamic Programming` adalah sebuah metode untuk memecahkan masalah yang bisa dipecahkan dengan cara memecahkannya menjadi submasalah yang lebih kecil, kemudian menyimpan solusi dari setiap submasalah yang telah dipecahkan untuk menghindari perhitungan yang berulang. Teknik ini umumnya digunakan untuk mengoptimalkan masalah di mana solusi optimal dapat dicapai dengan menggabungkan solusi dari submasalah yang lebih kecil.
### Penerapan Algoritma Dynamic Programming pada Program
Pada kode program ini, terdapat fungsi `searchCourses` yang menggunakan `algoritma Dynamic Programming` untuk menyelesaikan masalah `Knapsack (rucksack) dengan sedikit variasi`. Masalah Knapsack umumnya adalah masalah di mana kita memiliki beberapa item dengan nilai tertentu dan kapasitas tertentu di dalam knapsack, dan kita perlu memilih kombinasi item yang menghasilkan nilai maksimum tanpa melebihi kapasitas knapsack.

Dalam kasus kode di bawah, kita memiliki `sejumlah mata kuliah dengan bobot SKS tertentu` dan `nilai prediksi untuk setiap mata kuliah`. Tujuan dari algoritma ini adalah untuk memilih kombinasi mata kuliah yang memiliki bobot SKS maksimum yang tidak melebihi batas maksimum yang telah ditentukan, dengan nilai prediksi total yang maksimal.
### Analisis Penerapan Algoritma Dynamic Programming
1. Langkah Persiapan
- Mata kuliah yang relevan dengan jurusan dan fakultas masukan diambil menggunakan fungsi `getCoursesByJurusanFakultas`.
Kemudian, mata kuliah yang sesuai dengan semester yang diinginkan (semesterPengambilan) dipilih dan disimpan dalam `validCourses`.
2. Inisialisasi DP Table
- `dp` adalah tabel dua dimensi untuk menyimpan nilai prediksi maksimum yang dapat dicapai dengan memilih kombinasi mata kuliah tertentu, dengan jumlah SKS yang telah ditentukan. dp[i][w] akan berisi nilai prediksi maksimum dengan mempertimbangkan i mata kuliah pertama dan memiliki total SKS maksimum w.
2. Dynamic Programming - Proses Pengisian DP Table
- Melalui iterasi untuk setiap mata kuliah dan total SKS yang mungkin, tabel dp diisi dengan nilai prediksi maksimum.
Jika SKS dari mata kuliah ke-i (validCourses[i-1].JumlahSks) lebih kecil dari atau sama dengan total SKS w, maka ada dua pilihan:
Memasukkan mata kuliah ke-i dalam kombinasi dengan mempertimbangkan nilai prediksi getKonversiValue(validCourses[i-1].PrediksiNilai) atau tidak memasukkan mata kuliah ke-i dalam kombinasi.Pilihan terbaik diambil menggunakan fungsi max untuk memperbarui nilai dp[i][w].
3. Memilih Mata Kuliah yang Menghasilkan Nilai Prediksi Maksimum
- Setelah tabel dp terisi, kita dapat menemukan kombinasi mata kuliah dengan nilai prediksi maksimum yang tidak melebihi batas SKS (maxSKS).
Mata kuliah yang terpilih disimpan dalam `selectedCourses` dengan membalik iterasi melalui tabel dp.
4. Hasil
- `selectedCourses` berisi daftar mata kuliah yang dipilih yang akan dikembalikan sebagai respons dari fungsi `searchCoursesAPI`.
### Analisis Kompleksitas
Untuk menganalisis kompleksitas algoritma di atas, kita perlu mempertimbangkan dua tahap utama dari algoritma tersebut:

1. Inisialisasi DP Table: Algoritma melakukan inisialisasi tabel dp dengan ukuran n+1 (jumlah mata kuliah + 1) baris dan maxSKS+1 kolom, di mana n adalah jumlah mata kuliah yang sesuai dengan semester yang diinginkan. Oleh karena itu, kompleksitas inisialisasi tabel adalah O(n * maxSKS).

2. Proses Pengisian DP Table: Algoritma menggunakan dua loop bersarang untuk mengisi nilai-nilai tabel dp. Loop pertama memiliki n iterasi, dan loop kedua memiliki maxSKS iterasi. Dalam setiap iterasi, terdapat operasi konstan yang dilakukan (cek kondisi, pemilihan nilai maksimum dengan fungsi max, dan operasi aritmatika lainnya). Oleh karena itu, kompleksitas pengisian tabel adalah O(n * maxSKS).

Oleh karena kedua tahap tersebut dilakukan secara terpisah (tidak bersarang), kompleksitas keseluruhan algoritma adalah penjumlahan dari kedua kompleksitas tersebut:

Kompleksitas Total = O(n * maxSKS) + O(n * maxSKS) = O(n * maxSKS)

Jadi, kompleksitas algoritma Dynamic Programming di atas adalah O(n * maxSKS), dengan n adalah jumlah mata kuliah yang sesuai dengan semester yang diinginkan dan maxSKS adalah kapasitas maksimum SKS yang diizinkan.

<a name="videocapture"></a>

## Video Capture
<nl>

![CourseScheduler Gif](https://github.com/rifqifarhansyah/CourseScheduler_Backend/blob/main/img/CourseScheduler.gif?raw=true)

<a name="screenshots"></a>

## Screenshots
<p>
  <p>Gambar 1. Search Page</p>
  <img src="/img/SS1.png/">
  <nl>
  <p>Gambar 2. Add Page via Text Field</p>
  <img src="/img/SS2.png/">
  <nl>
  <p>Gambar 3. Add Page via JSON File</p>
  <img src="/img/SS3.png/">
  <nl>
   <p>Gambar 4. JSON File Example</p>
   <img src="/img/SS4.png/">
   <nl>
    <p>Gambar 5. Data and Remove Page</p>
   <img src="/img/SS5.png/">
   <nl>
    <p>Gambar 6. Searching Result</p>
   <img src="/img/SS6.png/">
   <nl>
</p>

<a name="structure"></a>

## Structure
```bash
│   .env
│   docker-compose.yml
│   Dockerfile
│   go.mod
│   go.sum
│   main
│   main.go
│   README.md
│
├───img
│       CourseScheduler.gif
│       SS1.png
│       SS2.png
│       SS3.png
│       SS4.png
│       SS5.png
│       SS6.png
│
└───tes
        tes1.json
        tes2.json
```

<a name="project-status">

## Project Status
Project is: _complete_

<a name="room-for-improvement">

## Room for Improvement
Perbaikan yang dapat dilakukan pada program ini adalah:
- Meningkatkan efektifitas algoritma dynamic programming

<a name="acknowledgements">

## Acknowledgements
- Terima kasih kepada Tuhan Yang Maha Esa

<a name="contact"></a>

## Contact
<h4 align="center">
  Kontak Saya : mrifki193@gmail.com<br/>
  2023
</h4>
<hr>
