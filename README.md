
# USPTO Design Patent Search Engine

## Background

The United States Patent and Trademark Office (USPTO) provides a dataset of design patents, which includes information about various design patents granted by the USPTO. This project aims to create a search engine that allows users to search for design patents based on various criteria.

## Requirements

- **Data Parsing**: The application downloads and parses the USPTO Design Patent dataset, available in XML format. The dataset can be obtained from [USPTO's official site](https://bulkdata.uspto.gov/data/patent/grant/redbook/2023/).

- **Database**: The parsed data is stored efficiently in a MongoDB database, considering performance and scalability.

- **Search Functionality**: Users can search for design patents based on:
  - Patent Title
  - Patent Number
  - Inventor(s) Name
  - Assignee (Owner) Name
  - Application Date
  - Issue Date
  - Design Class (if available)

- **Documentation**: Clear documentation is provided on how to set up and run the search engine, including instructions for downloading and parsing the USPTO Design Patent dataset.

- **Performance Optimization**: Optimizations have been implemented to ensure efficient performance, even with a large dataset.

## Data Model and System Design

The system is designed to be scalable and efficient in handling large datasets. The data model is structured to optimize the storage and retrieval of patent data. Below is a link to the system design and data model diagram:

![Data Model Diagram](data-model.png)

![System Design Diagram](design.png)


## Getting Started

### Prerequisites

- GoLang installed on your machine.
- Docker and Docker Compose installed on your machine.
- Access to the USPTO Design Patent dataset.

### Setup and Running

1. **Clone the Repository**:
   ```sh
   git clone <repository-url> search-app
   cd search-app
   ```

2. **Build the Database Docker Image and Start the Database Services**:
   ```sh
   make build-db
   ```

3. **Build the Project**:
   ```sh
   make build
   ```

4. **Run Tests**:
   ```sh
   make test
   ```

5. **Run the Project using Docker Compose**:
   ```sh
   make docker-compose-up
   ```

6. **View Logs**:
   ```sh
   make docker-log
   ```

7. **Clean Up**:
   ```sh
   make clean
   ```

8. **Stop the Docker Compose Services**:
   ```sh
   make docker-compose-down
   ```

## Documentation

For a detailed guide on how to use the search engine, including endpoints and example requests, refer to the provided Postman documentation available at [api.html](api.html).

## Evaluation Criteria

**Candidate POV**

- [X] **Functionality**: All Functionality Met
- [X] **Efficiency**: Search is done in an efficient manner using optimized queries.
- [X] **Code Quality**: Code adheres to Go Lang Best Practices.
- [X] **Documentation**: [api.html](api.html) is well documented.
- [X] **Database Design**: MongoDB is used to store the data, ensuring scalability and performance.
- [X] **Error Handling**: Comprehensive error handling is implemented to manage potential issues.
- [X] **Overall Impression**: The project is well-executed, with efficient data parsing and retrieval mechanisms in place.

## Deadline

The project is expected to be completed and submitted within 5 days from the start date.

## License

This project is licensed under the MIT License. See the LICENSE.md file for details.

