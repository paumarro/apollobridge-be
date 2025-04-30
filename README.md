# **Apollo Bridge**
Apollo Bridge is a cutting-edge mobile app designed to revolutionize the gallery and museum experience. With Apollo Bridge, visitors can seamlessly interact with artworks by simply scanning them with their phones. The app provides instant access to curated information about each piece, delivered directly by the gallery or museum.

Whether you're an art enthusiast or a casual visitor, Apollo Bridge bridges the gap between art and technology, offering a deeper, more personalized understanding of the exhibits around you.

# ***Lets get Started

### **1. Prerequisites**
Ensure you have the following installed on your system:
- **Go**: Version `1.23.0` or higher. You can download it from [golang.org](https://golang.org/dl/).  
- **PostgreSQL**: Ensure PostgreSQL is installed and running.  
- **SQLite**: Ensure the SQLite library is installed.  
- **Git**: For cloning the repository.  

Verify your installations:
```bash
go version
git --version
```

---

### **2. Clone the Repository**
Clone your backend repository to your local environment:
```bash
git clone https://github.com/paumarro/apollo-be.git
cd apollo-be
```

---

### **3. Install Dependencies**
Run the following command to install all Go module dependencies specified in `go.mod`:
```bash
go mod tidy
```

This will download and organize all required dependencies.

---

### **4. Set Up Environment Variables**
The server depends on certain environment variables for configuration. Create a `.env` file in the root directory of the project and define the required variables.

```env
DATABASE_URL=postgres://username:password@localhost:5432/apollo_db
JWT_SECRET=your_secret_key
PORT=8080
```

> Replace `DATABASE_URL` with your actual database connection string. If you're using SQLite, you can specify the path to the SQLite database file, e.g., `DATABASE_URL=sqlite3://./apollo.db`.

---

### **5. Set Up the Database (if required)**
If your backend uses a database (like PostgreSQL or SQLite), ensure the database is set up:
1. **PostgreSQL**:
   - Create a database:
     ```bash
     createdb apollo_db
     ```
   - Update the `DATABASE_URL` in your `.env` file accordingly.

---

### **6. Run the Backend**
Start the backend server with the following command:
```bash
go run .
```

This will compile and run the application. By default, the server will run on the port specified in the `.env` file (e.g., `8080`).

---

### **7. Test the Backend**
Once the server is running, test it to ensure it’s working:
- Open a browser or use a tool like [Postman](https://www.postman.com/) or `curl` to make requests to the server.
- Example: Test the health endpoint (if implemented):
  ```bash
  curl http://localhost:8080/health
  ```

---

### **8. Troubleshooting**
- If you encounter issues with dependencies, try running:
  ```bash
  go clean -modcache
  go mod tidy
  ```
- Check the logs for errors related to environment variables or database connections.

---

### **9. Optional: Build the Application**
To create a binary for deployment, use:
```bash
go build -o apollo-be
```

You can then run the binary with:
```bash
./apollo-be
```

---


# **Thread Model**

![ThreadModelLowRes](https://github.com/user-attachments/assets/48a9a72c-fe28-4e86-8166-88595280a2b8)

# Cyber Security Measures
## **Cyber Security Measures Implemented by Railway**
1. **HTTPS Everywhere**  
   - Automatic provisioning of TLS/SSL certificates for deployed applications, ensuring HTTPS is enforced for all connections.

2. **Secure Defaults**  
   - Encourages secure configurations by default, such as isolating environments and providing secure connection strings for databases.

3. **Environment Variables**  
   - Provides a secure way to manage secrets and environment variables, reducing the risk of exposing sensitive information in code.

4. **Infrastructure Security**  
   - Handles server and network-level security, including firewalls and protection against common infrastructure-level attacks like DDoS.

5. **Content Delivery Network (CDN)**  
   - Offers integration with a CDN to help with caching and mitigating certain attacks like Distributed Denial of Service (DDoS).

6. **Access Control**  
   - Role-based access control (RBAC) ensures only authorized users can manage deployments and resources.

---

### **Security Measures Not Fully Handled by Railway**
1. **Input Validation**  
   - Ensures that all user-provided data is validated to prevent injection attacks, malformed data, or other vulnerabilities.

2. **Output Encoding**  
   - Encodes data before rendering it to users, preventing injection-based attacks like Cross-Site Scripting (XSS).

4. **Security Headers**  
   - Implements headers like `Content-Security-Policy`, `X-Frame-Options`, and `Strict-Transport-Security` to mitigate various browser-based attacks.

5. **Header Management**  
   - Ensures secure handling of common headers like `Authorization` and `Content-Type` to prevent header injection and ensure proper API communication.

6. **Authentication and Authorization**  
   - Enforces identity verification and access control through mechanisms such as Keycloak integration, JWT tokens, and validation of user credentials.

7. **Session Management**  
   - Implements secure session handling, including short-lived access tokens (e.g., 10-minute lifespan) and rotating refresh tokens to reduce the risk of token theft.

8. **Logging and Monitoring**  
   - Tracks system activity, detects anomalies, and alerts administrators to potential security incidents.

9. **Database Security**  
   - Enforces security measures like parameterized queries to prevent SQL injections, access controls, and regular audits to protect sensitive data stored in databases. ORM like `GORM` to interact with the database securely.
- Always use :

10. **Rate Limiting**  
   - Limits the number of requests a user or IP can make within a specific time frame, mitigating brute force and denial-of-service attacks.

11. **File Upload Security (in the Frontend and Storage Server)**  
   - Scans and validates uploaded files to prevent malicious files from being executed or stored on the server.

12. **Dependency Management**  
   - CI Pipeline for regularly updates and audits third-party libraries and dependencies to ensure they are free from known vulnerabilities with tools like GOSEC and GOVULNCHECK
