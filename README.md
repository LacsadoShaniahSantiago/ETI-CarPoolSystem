# ETI-CarPoolSystem
## Design consideration of your microservices
1. User Management:
Separate Microservices: Create separate microservices for user account creation and management for passengers and car owners. This allows for independent scaling and deployment based on each user group's needs.
API Gateway: Implement an API gateway for user management API calls to provide a unified point of entry and enhance security.
Data Management: Store user data (passenger and car owner profiles) in a dedicated user database, potentially separate from trip data to improve data isolation and security.
Caching: Cache frequently accessed user data (e.g., basic profile info) to improve performance and reduce database load.

3. Trip Management:
Trip Service: Separate microservice for managing carpool trips, including publishing, searching, enrolment, and cancellation.
Real-time Updates: Utilize message queues or websockets to enable real-time trip updates (e.g., seat availability changes, cancellations) for seamless user experience.
Data Consistency: Implement database transactions or distributed locking mechanisms to ensure data consistency during concurrent trip operations.

5. Passenger Matching and Seat Assignment:
Matching Algorithm: Designed an efficient passenger matching algorithm considering trip preferences, location proximity, and first-come-first-serve basis.
Scalability: Use asynchronous processing for matching tasks to prevent blocking the main service and handle high demand during peak hours.
Notifications: Notify both passengers and car owners about successful enrolment and cancellation events.

4. Historical Trip Records:
Separate Microservice: Consider a separate microservice for retrieving historical trip records to decouple functionality and optimize querying performance.
Data Retention: Implement data retention policies to comply with the 1-year audit requirement while reducing storage costs.
Anonymization: Anonymize sensitive passenger and trip data after the retention period for privacy protection.

## Architecture diagram 
Components: 
- API Gateway: Acts as the single entry point for all API calls, managing authentication, authorization, and routing to specific microservices.
- User Management Microservice: Handles user account creation, management, and updates for both passengers and car owners.
- Passenger Profile Database: Stores passenger profile information (name, contact details).
- Car Owner Profile Database: Stores car owner profile information (driver's license, car plate).
- Trip Management Microservice: Manages carpool trip publishing, searching, enrolment, cancellation, and seat assignment.
- Message Queue: Facilitates asynchronous communication between microservices for real-time trip updates (e.g., seat availability changes, cancellations).
- Trip Database: Stores carpool trip information (pick-up locations, start time, destination, number of seats, enrolled passengers).
- Historical Trip Records Microservice: Manages retrieval of past carpool trip records.
- Historical Trip Records Database: Stores anonymized trip records after the retention period.
  
Connections:
- Users interact with the platform through the API Gateway.
- User Management Microservice interacts with the Passenger Profile Database and Car Owner Profile Database for profile data storage and retrieval.
- Trip Management Microservice interacts with the Geocoding Service for location-based search, Message Queue for real-time updates, and Trip Database for storing and 
  retrieving trip information.
- Passenger Matching and Seat Assignment are performed within the Trip Management Microservice.
- Historical Trip Records Microservice interacts with the Historical Trip Records Database to manage past trip data.

## Instructions for setting up and running your microservices
1. Opening Connection
   - In each folder, open Integraded Terminal in Vistual Studio Code and run each golang files in separate integrated terminals.
   - For mod folder, make sure connection is established which is indicated by "Listening to <port>" in the integrated terminal.
  
2. Establishing Database
   - Open sqlfile and run SQLDatabase file once in mySQL.
     * Remember to uncomment the first line of code to create carpooling_db database.
