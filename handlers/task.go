package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/AnasSabbag/task-manager/models"
)

// In-memory tasks
var tasks = []models.Task{
	{ID: "1", Title: "Sample Task", Description: "This is a sample task", Status: "Pending"},
}

// GetTasks returns all tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// CreateTask adds a new task and broadcasts it
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	tasks = append(tasks, newTask)

	// Broadcast new task to all WebSocket clients
	SendMessageToClients(models.Message{
		Type: "taskAdded",	
		Task: newTask,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var updatedTask models.Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	taskID := vars["id"]

	found := false
	for i, t := range tasks {
		if t.ID == taskID {
			tasks[i] = updatedTask
			found = true

			// Broadcast update
			SendMessageToClients(models.Message{Type: "taskUpdated", Task: updatedTask})
			break
		}
	}

	if !found {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}



func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	found := false
	for i, t := range tasks {
		if t.ID == taskID {
			// Remove the task
			tasks = append(tasks[:i], tasks[i+1:]...)
			found = true

			// Broadcast delete
			SendMessageToClients(models.Message{Type: "taskDeleted", Task: t})
			break
		}
	}

	if !found {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

