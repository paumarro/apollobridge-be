### **Thread Model**

![ThreadModelLowRes](https://github.com/user-attachments/assets/48a9a72c-fe28-4e86-8166-88595280a2b8)


### **Cyber Security Measures Implemented by Railway**
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
