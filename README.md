# 📌 Bencode Parser in Go

A fast, efficient, and fully tested **Bencode encoder/decoder** written in Go.

![Bencode Encoding and Decoding]

---

## 🚀 Features  
✔️ Encode and decode **integers, strings, lists, and dictionaries**.  
✔️ Uses **generics** for flexible data handling.  
✔️ **Optimized** for performance with benchmarking.  
✔️ Fully **tested** with `testify`.  
✔️ **Modularized** for easy imports and contributions.  

---

## 📖 How to Use  
### 1️⃣ Install the Package  
Run this command in your Go project:  
```sh
go get github.com/joekingsleyMukundi/bencoder_go
```

### 2️⃣ Import the Package  
In your Go file:  
```go
package main

import (
	"fmt"
	"log"

	"github.com/joekingsleyMukundi/bencoder_go"
)

func main() {
	// Encode a dictionary
	data := map[string]interface{}{
		"name": "Alice",
		"age":  25,
	}
	encoded, err := bencode.Encode(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Encoded:", encoded) // Example: d3:agei25e4:name5:Alicee

	// Decode back
	decoded, err := bencode.Decode([]byte(encoded))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Decoded:", decoded)
}
```

---

## ⚙️ How It Works  
This package implements **Bencode encoding and decoding** as used in BitTorrent.

### 🔹 Encoding Rules  
| Type      | Format Example |
|-----------|----------------|
| **Integer** | `i123e` (Integer 123) |
| **String** | `4:Test` (`Test`) |
| **List** | `l4:spam4:eggse` (`["spam", "eggs"]`) |
| **Dictionary** | `d3:agei25e4:name5:Alicee` (`{"age": 25, "name": "Alice"}`) |

### 🔹 Decoding Rules  
- Reads a **byte stream** and determines the data type.  
- Uses **buffer-based parsing** for efficiency.  
- Supports **nested lists and dictionaries**.

---

## 🛠️ Development & Contribution  

### 1️⃣ Clone the Repo  
```sh
git clone https://github.com/joekingsleyMukundi/bencoder_go.git
cd bencoder_go
```

### 2️⃣ Install Dependencies  
```sh
go mod tidy
```

### 3️⃣ Run Tests  
```sh
go test ./...
```

### 5️⃣ Submit a Pull Request  
- **Fork** the repo  
- **Create** a new branch  
- **Make** your changes  
- **Open** a PR for review  

---

## 📷 Project Image  
Be sure to upload the image to your repo and replace `path-to-your-image.png` above with the correct path.

---

## 📜 License  
This project is licensed under the **MIT License**.

---

🚀 **Happy coding!** 🎉
