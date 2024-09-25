package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // Importa o driver MySQL
)

// Estrutura do paciente
type Paciente struct {
	FullName        string `json:"fullName"`
	CPF             string `json:"cpf"`
	Address         string `json:"address"`
	City            string `json:"city"`
	State           string `json:"state"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Hospital        string `json:"hospital"`
	CardNo          string `json:"cardNo"`
	AppointmentDate string `json:"appointmentDate"`
}

// Função para conectar ao banco de dados MySQL
func ConectarBanco() (*sql.DB, error) {
	// String de conexão para MySQL
	// formato: "<usuário>:<senha>@tcp(<host>:<porta>)/<nome-do-banco>"
	dsn := "user:password@tcp(localhost:3306)/pacientes" // Substitua com seus dados

	// Abre a conexão com o banco de dados MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Erro de conexão com o banco de dados: %v", err)
	}

	// Verifica se a conexão foi bem-sucedida
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Erro de conexão com o banco de dados: %v", err)
	}

	// Criar a tabela se não existir
	sqlStmt := `CREATE TABLE IF NOT EXISTS pacientes (
		id INT AUTO_INCREMENT PRIMARY KEY,
		fullName VARCHAR(255),
		cpf VARCHAR(20),
		address VARCHAR(255),
		city VARCHAR(100),
		state VARCHAR(100),
		phone VARCHAR(20),
		email VARCHAR(255),
		hospital VARCHAR(255),
		cardNo VARCHAR(20),
		appointmentDate VARCHAR(20)
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("Erro de criar tabela: %v", err)
	}

	return db, nil
}

func ArmazenarPaciente(db *sql.DB, paciente Paciente) error {
	query := `
	INSERT INTO pacientes (fullName, cpf, address, city, state, phone, email, hospital, cardNo, appointmentDate)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := db.Exec(query, paciente.FullName, paciente.CPF, paciente.Address, paciente.City, paciente.State, paciente.Phone, paciente.Email, paciente.Hospital, paciente.CardNo, paciente.AppointmentDate)
	if err != nil {
		return fmt.Errorf("Erro inserindo paciente: %v", err)
	}
	return nil
}

func handleCadastroPaciente(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var paciente Paciente
	err := json.NewDecoder(r.Body).Decode(&paciente)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	db, err := ConectarBanco()
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = ArmazenarPaciente(db, paciente)
	if err != nil {
		http.Error(w, "Erro ao armazenar paciente no banco de dados", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Paciente recebido e armazenado com sucesso: %s\n", paciente.FullName)
}
func main() {
	http.HandleFunc("/cadastro-paciente", handleCadastroPaciente)
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
