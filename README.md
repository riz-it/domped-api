# Domped Digital API

[API Domped Digital](https://domped-api.vercel.app/)

## Development

1. **Clone or Pull the Repository**

   ```bash
   git clone git@github.com:riz-it/domped-api.git
   ```

2. **Configure Environment Variables**

   - Copy the `env.example` file to `.env`, then adjust the configurations as needed.

3. **Build Docker Image** _(optional if you have not built the image before, adjust according to your OS)_

   - **Linux/Mac**:
     ```bash
     DOCKER_BUILDKIT=1 docker build -t github.com/riz-it/domped-api.git -f Dockerfile.dev .
     ```
   - **Windows**:
     ```bash
     docker build -t github.com/riz-it/domped-api.git .
     ```

4. **Run the Container**

   ```bash
   docker run --rm -it -v $(pwd):/app -w /app -p 9009:9009 domped-digital
   ```

5. **Start Development**

   - Once the container is running, begin your development process.

## Deployment

1. **Clone or Pull the Repository**

   ```bash
   git clone git@github.com:riz-it/domped-api.git
   ```

2. **Configure Environment Variables**

   - Copy the `env.example` file to `.env`, then adjust the configurations as needed.

3. **Build Docker Image** _(optional if you have not built the image before, adjust according to your OS)_

   - **Linux/Mac**:
     ```bash
     DOCKER_BUILDKIT=1 docker build -t github.com/riz-it/domped-api.git -f Dockerfile.prod .
     ```
   - **Windows**:
     ```bash
     docker build -t github.com/riz-it/domped-api.git .
     ```

4. **Run the Container**

   ```bash
   docker run --env-file .env --name domped-digital -d -p 9009:9009 github.com/riz-it/domped-api.git
   ```

5. **Domped Digital API is Ready to Use**

   - Once the container is running, you can start using the API.

## Cleaning up Docker

- **Stop and Remove the `domped-digital` Container** _(adjust according to your OS)_

  - **Linux/Mac**:
    ```bash
    docker stop domped-digital || true && docker rm domped-digital || true
    ```
  - **Windows**:
    ```bash
    docker stop domped-digital
    docker rm domped-digital
    ```

- **Remove All Containers**

  ```bash
  docker rm -vf $(docker ps -aq)
  ```

- **Prune Docker System** _(removes all unused volumes, images, and cache)_

  ```bash
  docker system prune --volumes -af
  ```
