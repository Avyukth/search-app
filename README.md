
# USPTO Design Patent Search Engine

## Background

The United States Patent and Trademark Office (USPTO) provides a dataset of design patents, which includes information about various design patents granted by the USPTO. This project aims to create a search engine that allows users to search for design patents based on various criteria.

## Requirements

- **Data Parsing**: The application downloads and parses the USPTO Design Patent dataset, available in XML format. The dataset can be obtained from [USPTO's official site](https://bulkdata.uspto.gov/data/patent/grant/redbook/2023/).

- **Database**: The parsed data is stored efficiently in a database. The choice of the database system was made considering performance and scalability.

- **Search Functionality**: Users can search for design patents based on:
  - Patent Title
  - Patent Number
  - Inventor(s) Name
  - Assignee (Owner) Name
  - Application Date
  - Issue Date
  - Design Class (if available)

- **Documentation**: Clear postman documentation is provided on how to set up and run the search engine, including instructions for downloading and parsing the USPTO Design Patent dataset.

- **Performance Optimization**: Optimizations have been implemented to ensure efficient performance, even with a large dataset.

## Getting Started

### Prerequisites

- GoLang installed on your machine.
- Access to the USPTO Design Patent dataset.

### Setup and Running

1. **Clone the Repository**:
   ```sh
   cd search-app
   ```

2. **Build the Project**:
   ```sh
   make build
   ```

3. **Run Tests**:
   ```sh
   make test
   ```

4. **Run the Project**:
   ```sh
   make run
   ```

5. **Clean Up**:
   ```sh
   make clean
   ```

## Documentation

For a detailed guide on how to use the search engine, including endpoints and example requests, refer to the provided Postman documentation.

## Database Schema and Data Parsing

The database schema has been designed to efficiently store and retrieve patent data. The primary tables include `Patents`, `Inventors`, `Assignees`, and `DesignClasses`. Relationships between these tables ensure data integrity and efficient querying.

Data parsing from the XML dataset involves extracting relevant fields and transforming them into a format suitable for database storage. Special attention was given to handle inconsistencies in the dataset and ensure accurate data extraction.

## Evaluation Criteria

**Candidate POV**

- [X] **Functionality**: All Functionality Met
- [X] **Efficiency** search is done in Efficient manner
- [X] **Code Quality** COde Quality is Go Lang Best Practices
- [X] **Postman Documentation** [api.Html](api.html) is Documented
- [X] **Database Design** Mongo is used to store the data
- [X] **Error Handling** Error Handling is done
- [X] **Overall Impression**: General assessment of the project as a whole is good , took time to parse the data efficiently.

## Deadline

The project is expected to be completed and submitted within 5 days from the start date.

## License

This project is licensed under the MIT License. See the LICENSE.md file for details.


