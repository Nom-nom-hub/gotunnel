package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := "3000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>GoTunnel Example</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; }
        .info { background: #e3f2fd; padding: 15px; border-radius: 4px; margin: 20px 0; }
        .code { background: #f5f5f5; padding: 10px; border-radius: 4px; font-family: monospace; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 GoTunnel Example Server</h1>
        <div class="info">
            <strong>Status:</strong> Running successfully!<br>
            <strong>Port:</strong> %s<br>
            <strong>Path:</strong> %s<br>
            <strong>Method:</strong> %s<br>
            <strong>Remote Address:</strong> %s
        </div>
        
        <h2>How to use GoTunnel:</h2>
        <ol>
            <li>Start the GoTunnel server: <span class="code">./gotunnel-server --port 443 --cert cert.pem --key key.pem --allowed-tokens "your-token"</span></li>
            <li>Start the GoTunnel client: <span class="code">./gotunnel-client --server your-domain.com --subdomain myapp --local-port 3000 --token "your-token"</span></li>
            <li>Access your app at: <span class="code">https://myapp.your-domain.com</span></li>
        </ol>
        
        <h2>Request Headers:</h2>
        <div class="code">
`, port, r.URL.Path, r.Method, r.RemoteAddr)
		
		for name, values := range r.Header {
			for _, value := range values {
				fmt.Fprintf(w, "%s: %s\n", name, value)
			}
		}
		
		fmt.Fprintf(w, `        </div>
    </div>
</body>
</html>`)
	})

	addr := ":" + port
	fmt.Printf("Starting example server on port %s...\n", port)
	fmt.Printf("Visit http://localhost%s to see the example\n", addr)
	fmt.Printf("Use Ctrl+C to stop the server\n")
	
	log.Fatal(http.ListenAndServe(addr, nil))
} 