# Grocademy Web Monolith
Website untuk pemenuhan tugas Seleksi LabPro 3

## Run (local)

Buat .env di root berisi (misal):
```shell
APP_PORT=8080

DB_HOST=localhost
DB_USER=user
DB_PASSWORD=password
DB_NAME=grodb
DB_PORT=5432
DB_SSLMODE=disable

CLOUDINARY_URL=<cloudinary-api>
JWT_SECRET_KEY=secret-key
```
Lalu jalankan perintah berikut:
```shell
make build_app # make sure docker and make is available
```
Kunjungi http://localhost:8080


## Design Pattern
1. Dependency Injection (DI), untuk menginjek objek service ke handler.
3. Repository Pattern, memisahkan data access dari logika bisnis. Kelas service enggunakan GORM.
4. Middleware Pattern, untuk menggunakan objek yang sama dalam mengatur alur request dan respons API. Misalnya auth middleware.

## Tech Stack
* Go v1.24.2
* Gin v1.10.1
* GORM v1.30.1
* Postgres

## Endpoint
- auth
  - POST /auth/login
  - POST /auth/register
  - GET /auth/self

- courses
  - GET /courses
  - POST /courses
  - GET /courses/{id}
  - PUT /courses/{id}
  - DELETE /courses/{id}
  
- modules
  - GET /courses/{courseId}/modules
  - POST /courses/{courseId}/modules
  - PATCH /courses/{courseId}/modules/reorder
  - GET /modules/{id}
  - PUT /modules/{id}
  - DELETE /modules/{id}
  - PATCH /modules/{id}/complete

- users
  - GET /users
  - POST /users
  - GET /users/{id}
  - PUT /users/{id}
  - DELETE /users/{id}
  - POST /users/{id}/balance
 
## Bonus
- B2 - [Deployment](https://grocademy-monolith-production.up.railway.app/)
- B3 - Polling (short-polling untuk course)
- B6 - Responsive UI
- B7 - [Dokumentasi API](https://grocademy-monolith-production.up.railway.app/docs/index.html)
- B8 - SOLID
  - Single Responsibility Principle
SRP menyatakan bahwa modul, kelas, atau fungsi hanya boleh memiliki satu tanggung jawab. Dalam Go, kita gunakan struct. Setiap model (User, Course, Module) hanya bertanggung jawab untuk merepresentasikan struktur data dan skema basis datanya. Setiap layanan (UserService, CourseService, ModuleService) bertanggung jawab atas logika bisnis yang terkait dengan entitasnya masing-masing. Setiap handler (misalnya, UserHandler, CourseHandler) bertanggung jawab untuk mengurai permintaan HTTP, memanggil servis yang sesuai, dan memformat respons HTTP.
  - Open Closed Principle
OCP menyatakan bahwa entitas perangkat lunak harus terbuka untuk ekstensi, tetapi tertutup untuk modifikasi. Dalam Go, interface adalah contoh penerapan ini. Misal struct UserService mengimplementasi UserServicer. Struct ini tidak bisa dapat menambahkan fungsi-fungsi baru untuk mengektensi fungsionalitas, tetapi tetap harus mengikuti interface.
  - Prinsip Substitusi Liskov menyatakan bahwa objek dari superkelas harus dapat digantikan dengan objek dari subkelasnya tanpa merusak aplikasi. Dalam Go, hal ini berlaku untuk interface. Jika suatu tipe mengimplementasikan suatu interface, tipe tersebut harus dapat digantikan dengan interface tersebut
  - Interface Segregation Principle menyatakan bahwa klien tidak boleh dipaksa bergantung pada antarmuka yang tidak mereka gunakan. Hal ini mendorong antarmuka yang lebih kecil dan lebih terfokus.
Contohnya, interface UserServicer, CourseServicer, ModuleServicer sudah cukup terperinci. UserServicer hanya berfokus pada operasi pengguna, CourseServicer pada operasi kursus, dll. Antarmuka ini tidak mengandung metode asing yang tidak terkait dengan domainnya.
  - Dependency Inversion Principle menyatakan bahwa modul tingkat tinggi tidak boleh bergantung pada modul tingkat bawah, keduanya harus bergantung pada abstraksi. Misal, UserHandler bergantung pada antarmuka UserServicer (sebuah abstraksi), bukan struct UserService konkret. Hal yang sama berlaku untuk AuthHandler, CourseHandler, dan ModuleHandler.
- B11 - Bucket, menggunakan Cloudinary

## Screenshot
![img1](https://github.com/L4mbads/grocademy-monolith/blob/6d59493847ed5040eb5a38143db5d7b01933951b/assets/Screenshot%202025-08-25%20004750.png)
![img2](https://github.com/L4mbads/grocademy-monolith/blob/6d59493847ed5040eb5a38143db5d7b01933951b/assets/Screenshot%202025-08-25%20004826.png)
![img3](https://github.com/L4mbads/grocademy-monolith/blob/6d59493847ed5040eb5a38143db5d7b01933951b/assets/Screenshot%202025-08-25%20005008.png)
![img4](https://github.com/L4mbads/grocademy-monolith/blob/6d59493847ed5040eb5a38143db5d7b01933951b/assets/Screenshot%202025-08-25%20005027.png)

## Identitas
- Fachriza Ahmad Setiyono (13523162)
