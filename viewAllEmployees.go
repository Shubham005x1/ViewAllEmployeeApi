package content

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)

type Employee struct {
	ID        string `firestore:"id" json:"id"`
	FirstName string `firestore:"firstname" json:"firstname"`
	LastName  string `firestore:"lastname" json:"lastname"`
	Email     string `firestore:"email" json:"email"`
	Password  string `firestore:"password" json:"-"`
	PhoneNo   string `firestore:"phoneNo" json:"phoneNo"`
	Role      string `firestore:"role" json:"role"`
}

var (
	client     *firestore.Client
	onceClient sync.Once
)

func InitializeFirestore() {
	onceClient.Do(func() {
		ctx := context.Background()

		// Initialize Firestore with the service account key
		var err error
		client, err = firestore.NewClient(ctx, "takeoff-task-3")
		if err != nil {
			log.Fatalf("Failed to create Firestore client: %v", err)
		}
	})
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/employees", ViewAllEmployees).Methods("GET")
	log.Println("Server started on :8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
func ViewAllEmployees(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	InitializeFirestore()

	iter := client.Collection("employees").Documents(ctx)
	defer iter.Stop()

	var employees []Employee

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve employees: %v", err), http.StatusInternalServerError)
			return
		}

		var employee Employee
		if err := doc.DataTo(&employee); err != nil {
			http.Error(w, fmt.Sprintf("Failed to convert data: %v", err), http.StatusInternalServerError)
			return
		}

		employees = append(employees, employee)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}
