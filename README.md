# Mistyped - Server

This is the server-side component of the Mistyped application, built with Golang. It handles the processing of URLs and checks for similar URLs by considering possible typos.

## Features

- **URL Processing:** Receives URLs from the client-side application and processes them for typo analysis.
- **Typo Analysis:** Analyzes possible typos to find similar URLs.

## Getting Started

To run the Mistyped server locally, follow these steps:

### Prerequisites

- [Go](https://golang.org/) (for Golang backend)

### Installation

1. Clone the repository:
`git clone https://github.com/your-username/mistyped-server.git`
2. Navigate to the project directory:
`cd mistyped-server`
3. Install dependencies:
`go get`
4. Build the server:
`go build`
5. Run the server:
`./mistyped-server`

## API Documentation

- **POST /check-url**
  - Description: Receives a URL from the client and returns a list of similar URLs.
  - Request body: `{ "url": "example.com" }`
  - Response body: `{ "similarUrls": ["exmaple.com", "exampl.com", ...] }`

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)