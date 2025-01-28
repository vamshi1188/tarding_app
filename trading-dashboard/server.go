package main

import (
	"fmt"
	"net/http"
)

// handler for the homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Trading Dashboard</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				margin: 0;
				padding: 0;
				background-color: #f4f4f9;
			}
			header {
				background-color: #333;
				color: white;
				padding: 1rem 0;
				text-align: center;
			}
			main {
				padding: 2rem;
				text-align: center;
			}
			button {
				background-color: #007bff;
				border: none;
				color: white;
				padding: 0.5rem 1rem;
				cursor: pointer;
				font-size: 1rem;
				border-radius: 5px;
			}
			button:hover {
				background-color: #0056b3;
			}
		</style>
	</head>
	<body>
		<header>
			<h1>Welcome to the Trading Dashboard</h1>
		</header>
		<main>
			<p>Here you can view real-time trading calls and market data.</p>
			<button onclick="alert('Fetching latest market data...')">Get Market Data</button>
		</main>
	</body>
	</html>
	`

	fmt.Fprint(w, html)
}

func main() {
	http.HandleFunc("/", homeHandler) // Route for the homepage

	fmt.Println("Starting server at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
