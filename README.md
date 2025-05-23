# Receipt Processor

## Introduction

The Receipt Processor Challenge is a API application that takes in a JSON receipt and returns a JSON object with an ID generated by code.
The ID returned is the ID that can get the number of points the receipt was awarded.


## Installation

To run this project locally, follow these steps:

1. **Clone the repository:**

        git clone https://github.com/brian891019/receipt-processor-challenge.git
    

2. **Change to the project directory:**

        cd receipt-processor-challenge
    

3. **Install Dependencies**

        go mod tidy
    
4. **Running the Application**
   
        go run main.go
   
6. **The output should indicate the server is running:**

        Starting server on :8080

## API Endpoints

- **Upload a Receipt: (Send the POST Request using curl)**
  
        curl -X POST http://localhost:8080/receipts/process -H "Content-Type: application/json" -d @receipt.json

    (Note: change the receipt.json file if you want to process other receipt)


    You should receive a response with an ID like:

        {"id": "generated-id"}

          

- **Use the Returned ID to Get Points**

       curl -X GET http://localhost:8080/receipts/{your-generated-id}/points
  
   Replace {generated-id} with the actual ID returned from the POST request. This should return the points associated with the receipt:

       {"points": <calculated-points>}


# Rules

* One point for every alphanumeric character in the retailer name.
* 50 points if the total is a round dollar amount with no cents.
* 25 points if the total is a multiple of `0.25`.
* 5 points for every two items on the receipt.
* If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
* 6 points if the day in the purchase date is odd.
* 10 points if the time of purchase is after 2:00pm and before 4:00pm.


## Examples

```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```
```text
Total Points: 28
Breakdown:
     6 points - retailer name has 6 characters
    10 points - 5 items (2 pairs @ 5 points each)
     3 Points - "Emils Cheese Pizza" is 18 characters (a multiple of 3)
                item price of 12.25 * 0.2 = 2.45, rounded up is 3 points
     3 Points - "Klarbrunn 12-PK 12 FL OZ" is 24 characters (a multiple of 3)
                item price of 12.00 * 0.2 = 2.4, rounded up is 3 points
     6 points - purchase day is odd
  + ---------
  = 28 points
```

----

```json
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
```
```text
Total Points: 109
Breakdown:
    50 points - total is a round dollar amount
    25 points - total is a multiple of 0.25
    14 points - retailer name (M&M Corner Market) has 14 alphanumeric characters
                note: '&' is not alphanumeric
    10 points - 2:33pm is between 2:00pm and 4:00pm
    10 points - 4 items (2 pairs @ 5 points each)
  + ---------
  = 109 points
```
## Running Tests
The project includes unit tests for both the service and handler packages.

To run tests, use the following command:

```bash
 go test ./service ./handler
