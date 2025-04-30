# **Apollo Bridge**

Apollo Bridge is a cutting-edge mobile app designed to revolutionize the gallery and museum experience. With Apollo Bridge, visitors can seamlessly interact with artworks by simply scanning them with their phones. The app provides instant access to curated information about each piece, delivered directly by the gallery or museum.

Whether you're an art enthusiast or a casual visitor, Apollo Bridge bridges the gap between art and technology, offering a deeper, more personalized understanding of the exhibits around you.

---

## **Let's Get Started**

### **1. Prerequisites**
Ensure you have the following installed on your system:
- **Docker**: Ensure Docker is installed and running. You can download it from [docker.com](https://www.docker.com/).
- **Docker Compose**: Comes bundled with Docker Desktop or can be installed separately. Verify installation:
  ```bash
  docker compose --version
  ```
- **Git**: For cloning the repository.

---

### **2. Clone the Repository**
Clone your backend repository to your local environment:
```bash
git clone https://github.com/paumarro/apollo-be.git
cd apollo-be
```

---

### **3. Configure Environment Variables**
Apollo Bridge uses environment variables to manage configuration securely. Create a `.env` file in the root directory of the project (or copy the provided `.env.example` file) and define the following variables:

```env
# Database Configuration
POSTGRES_USER=myuser
POSTGRES_PASSWORD=securepassword
POSTGRES_DB=art
ART_DB_URL=postgres://myuser:securepassword@db:5432/art

# Keycloak Configuration
KEYCLOAK_ADMIN=admin
KEYCLOAK_ADMIN_PASSWORD=secureadminpassword
KEYCLOAK_DOMAIN=http://keycloak:8080
KEYCLOAK_CLIENT_ID=apollo-client
KEYCLOAK_CLIENT_SECRET=your-client-secret
JWKS_URL=http://keycloak:8080/realms/apollo/protocol/openid-connect/certs
```

> **Note:** Replace `securepassword` and `secureadminpassword` with strong passwords. The `KEYCLOAK_CLIENT_SECRET` will be generated in the Keycloak Admin Console (explained below).

---

### **4. Set Up and Run with Docker**
The backend and Keycloak services are configured to run using Docker Compose. Follow these steps:

1. **Update the `docker-compose.yml` File with the defined env variables**  


2. **Start the Services**  
   Run the following command to start all services:
   ```bash
   docker compose up --build
   ```

   This will:
   - Build and run the backend service (`art-service`) on port `3000`.
   - Set up a PostgreSQL database (`db`) on port `5432`.
   - Start Keycloak (`keycloak`) on port `8080`.

3. **Access Keycloak**  
   - Open your browser and navigate to `http://localhost:8080`.
   - Log in to the Keycloak Admin Console using the credentials from your `.env` file (`KEYCLOAK_ADMIN` and `KEYCLOAK_ADMIN_PASSWORD`).

---

### **5. Configure Keycloak**
1. **Create a Realm**  
   - In the Keycloak Admin Console, create a new realm named `apollo`.

2. **Set Up a Client**  
   - In the `apollo` realm, navigate to the **Clients** section and click **Create**.
   - Set the following:
     - **Client ID**: `apollo-client`
     - **Client Protocol**: `openid-connect`
     - **Access Type**: `confidential`
   - Save the client, then go to the **Credentials** tab to copy the `Client Secret`. Add this to your `.env` file as `KEYCLOAK_CLIENT_SECRET`.

3. **Configure Redirect URIs**  
   - In the client settings, under the **Valid Redirect URIs** field, add the callback URL for your app:
     ```
     http://localhost:3000/auth/callback
     ```

4. **Set Up Users**  
   - Navigate to the **Users** section and create test users for authentication. Assign them appropriate roles if needed.

---

### **6. Test the Backend**
Once all services are running, test the backend to ensure it’s working:
- Open a browser or use a tool like [Postman](https://www.postman.com/) or `curl` to make requests to the server.

---

### **7. Troubleshooting**
- If Keycloak is inaccessible, ensure the `KEYCLOAK_ADMIN` and `KEYCLOAK_ADMIN_PASSWORD` in the `.env` file match the values in `docker-compose.yml`.
- If the backend cannot connect to the database, verify the `ART_DB_URL` format in the `.env` file.

---

### **8. Stopping the Services**
To stop all running services, use:
```bash
docker compose down
```

This will stop and remove all containers, but data in volumes (`db-data` and `keycloak-data`) will persist.

---

### **9. Optional: Build the Backend Locally**
If you prefer running the backend locally instead of Docker, follow the steps in the original instructions to set up Go, dependencies, and the database.

