# KaggleCarAPI

This repository contains the code for the KaggleCarAPI project, implemented in Go. The KaggleCarAPI is a RESTful API that provides access to a dataset of car information, allowing users to retrieve data such as car makes, models, and specifications.

## Dataset
https://www.kaggle.com/datasets/peshimaammuzammil/2023-car-model-dataset-all-data-you-need?resource=download

## Getting Started

To get started with the KaggleCarAPI, follow the steps below:

1. **Clone the repository:**

   ```shell
   git clone https://github.com/phllpmcphrsn/KaggleCarAPI.git
   ```

2. **Navigate to the project directory:**

   ```shell
   cd KaggleCarAPI
   ```

3. **Build the project:**

   ```shell
   go build
   ```

4. **Run the server:**

   ```shell
   ./KaggleCarAPI
   ```

   The API server will start running on `http://localhost:8080`.

## API Endpoints

The KaggleCarAPI provides the following endpoints:

- **GET /cars**

  Retrieves a list of all cars in the dataset.

- **GET /cars/{make}**

  Retrieves a list of cars for a specific make.

- **GET /cars/{make}/{model}**

  Retrieves information for a specific car make and model.

- **POST /cars**

  Adds a new car to the dataset. Requires the car make, model, and specifications in the request body.

- **PUT /cars/{make}/{model}**

  Updates the information of a specific car make and model. Requires the updated specifications in the request body.

- **DELETE /cars/{make}/{model}**

  Deletes a specific car from the dataset.

For detailed information about each endpoint and the expected request/response formats, please refer to the API documentation.

## Data Format

The car dataset used by the KaggleCarAPI is stored in the `data/cars.csv` file. Each row represents a car and contains information such as the make, model, year, engine type, and horsepower.

## Contributing

Contributions to the KaggleCarAPI project are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

The KaggleCarAPI project is licensed under the MIT License. Please see the `LICENSE` file for more information.

## Contact

If you have any questions or inquiries about the KaggleCarAPI project, you can reach out to the project owner, Philip McPherson, via email at philip.mcpherson@example.com.

Feel free to explore the repository and start using the KaggleCarAPI!