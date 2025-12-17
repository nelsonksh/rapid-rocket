# Andamio Protocol Explorer (Go + HTMX + Tailwind)

Welcome to your rewritten **Andamio Protocol Explorer**! This project replaces the original Next.js/React stack with a high-performance, lightweight stack using **Go**, **HTMX**, and **Tailwind CSS**.

## 🚀 Tech Stack

-   **Backend**: Go 1.25.2 (Standard Library `net/http`) - Handles routing, API logic, and HTML rendering.
-   **Frontend Interactivity**: [HTMX](https://htmx.org/) - Enables dynamic partial page updates without complex Client-Side JavaScript.
-   **Styling**: [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework (loaded via CDN for simplicity).
-   **Documentation**: [Scalar](https://scalar.com/) - Beautiful Open Source API references.

## 📂 Project Structure

```
AntigravityPlayground/
├── main.go              # THE MONOLITH: Server, Router, Data Handlers (Live & Mock), and Structs.
├── go.mod               # Go module definition.
├── views/               # HTML Templates
│   ├── index.html       # The Main Dashboard (Hero, Search, Layout).
│   ├── analytics.html   # Stats Grid Fragment.
│   ├── transactions.html# Recent Transactions Fragment.
│   ├── contributions.html# Recent Contributions Fragment.
│   ├── search.html      # Search Results Fragment.
│   └── docs.html        # Scalar API Documentation wrapper.
└── assets/
    └── openapi.yaml     # OpenAPI 3.0 Specification for your API.
```

## 🧠 Key Concepts & Architecture

### 1. Server-Side Rendering (SSR) & Fragments
Unlike the previous React App which fetched JSON and built DOM on the client, this app renders HTML on the server.
-   **Full Pages**: The `/` route returns the complete `index.html`.
-   **HTML Fragments**: API endpoints like `/api/analytics` or `/search` return *partial* HTML (e.g., just the list of transactions), not JSON.

### 2. HTMX Binding
The magic happens in `views/index.html`. We use HTML attributes to drive interactivity:

-   **Lazy Loading**:
    ```html
    <div hx-get="/api/analytics" hx-trigger="load" hx-swap="innerHTML">...</div>
    ```
    This tells the browser: *"When this loads, fetch content from /api/analytics and replace my inner HTML with the result."*

-   **Live Search**:
    ```html
    <input hx-get="/search" hx-trigger="keyup changed delay:200ms" hx-target="#search-results" ... >
    ```
    This sends the input value to `/search?q=...` 200ms after you stop typing, then puts the resulting HTML into `#search-results`.

### 3. Real Data & Mock Data Strategy
-   **Real Data**:
    -   **Analytics**: Fetches live stats from `preprod.andamioscan.andamio.space/v2/transactions/count`.
    -   **Search**: Fetches transaction details from `preprod.andamioscan.andamio.space/v2/transactions/{hash}`.
-   **Mock Data**: Recent Transactions and Contributions lists currently use in-memory mock data (defined in `main.go`).

## 🛠️ How to Develop

### Running the Server
```bash
go run main.go
```
Visit **http://localhost:8080** to see the app.

### API Documentation
Visit **http://localhost:8080/docs** to see the interactive Scalar API docs.
-   Edit `assets/openapi.yaml` to update the specs.
-   The docs auto-update when you refresh the page.

### Modifying Content
-   **Change UI/Layout**: Edit `views/index.html` or fragment templates in `views/*.html`.
-   **Change Data/Logic**: Edit `main.go`.

## ⏭️ Next Steps

1.  **Complete Data Integration**: Replace the mock data in `transactionsHandler` and `contributionsHandler` with real API calls.
2.  **Add Pages**: Create new handlers in `main.go` and corresponding templates to show full details.
3.  **Deployment**: Build the binary with `go build -o server main.go` and deploy!

Enjoy your new blazing fast explorer! 🚀
