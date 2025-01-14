# API Domped Digital

[API Domped Digital](https://domped-api.vercel.app/)

## Development

1. **Clone atau Pull Repository**

   ```bash
   git clone git@github.com:riz-it/domped-api.git
   ```

2. **Konfigurasi Environment Variables**

   - Salin file `env.example` ke `.env`, kemudian sesuaikan konfigurasi.

3. **Build Docker Image** _(opsional jika belum pernah membangun image, sesuaikan dengan OS yang digunakan)_

   - **Linux/Mac**:
     ```bash
     DOCKER_BUILDKIT=1 docker build -t github.com/riz-it/domped-api.git -f Dockerfile.dev .
     ```
   - **Windows**:
     ```bash
     docker build -t github.com/riz-it/domped-api.git .
     ```

4. **Jalankan Container**

   ```bash
   docker run --rm -it -v $(pwd):/app -w /app -p 9009:9009 domped-digital
   ```

5. **Mulai Pengembangan**

   - Setelah container berjalan, mulai pengembangan.

## Deployment

1. **Clone atau Pull Repository**

   ```bash
   git clone git@github.com:riz-it/domped-api.git
   ```

2. **Konfigurasi Environment Variables**

   - Salin file `env.example` ke `.env`, kemudian sesuaikan konfigurasi.

3. **Build Docker Image** _(opsional jika belum pernah membangun image, sesuaikan dengan OS yang digunakan)_

   - **Linux/Mac**:
     ```bash
     DOCKER_BUILDKIT=1 docker build -t github.com/riz-it/domped-api.git -f Dockerfile.prod .
     ```
   - **Windows**:
     ```bash
     docker build -t github.com/riz-it/domped-api.git .
     ```

4. **Jalankan Container**

   ```bash
   docker run --env-file .env --name domped-digital -d -p 9009:9009 github.com/riz-it/domped-api.git
   ```

5. **API Domped Digital Siap Digunakan**

   - Setelah container berjalan, mulai gunakan.

## Cleaning up Docker

- **Hentikan dan Hapus Container pmb** _(sesuaikan dengan OS yang digunakan)_

  - **Linux/Mac**:
    ```bash
    docker stop domped-digital || true && docker rm domped-digital || true
    ```
  - **Windows**:
    ```bash
    docker stop domped-digital
    docker rm domped-digital
    ```

- **Hapus semua container**

  ```bash
  docker rm -vf $(docker ps -aq)
  ```

- **Prune Docker system** _(menghapus semua volume, image, dan cache yang tidak digunakan)_

  ```bash
  docker system prune --volumes -af
  ```
