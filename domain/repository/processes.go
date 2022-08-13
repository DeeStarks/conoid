package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

type ServiceProcessModel struct {
	Name          string        `json:"name"`
	Status        bool          `json:"status"`
	Type          string        `json:"type"`
	Listeners     []interface{} `json:"listeners"`
	RootDirectory string        `json:"root"`
	RemoteServer  string        `json:"server"`
	Tunnelled     bool          `json:"is_tunnelled"`
	CreatedAt     int64         `json:"created_at"`
}

type ServiceProcess struct {
	DB *os.File
}

// Retrive all running services
func (p ServiceProcess) RetrieveRunning() ([]ServiceProcessModel, error) {
	b := make([]byte, 1024*100) // Max read from DB: 100kb
	n, err := p.DB.Read(b)
	if err != nil && err != io.EOF {
		return nil, err
	}

	var data map[string]interface{}
	if err = json.Unmarshal(b[:n], &data); err != nil {
		fmt.Println(string(b))
		return nil, err
	}

	names := reflect.ValueOf(data).MapKeys()
	var processes []ServiceProcessModel
	for _, name := range names {
		// Check service is running, the add to list of processes
		if data[name.String()].(map[string]interface{})["status"].(bool) {
			var process ServiceProcessModel
			process.Name = name.String()

			// Get service data
			service := data[name.String()].(map[string]interface{})
			process.Status = service["status"].(bool)
			process.Type = service["type"].(string)
			process.Listeners = service["listeners"].([]interface{})
			process.RootDirectory = service["root"].(string)
			process.RemoteServer = service["server"].(string)
			process.Tunnelled = service["is_tunnelled"].(bool)
			process.CreatedAt = int64(service["created_at"].(float64))

			// Append
			processes = append(processes, process)
		}
	}
	return processes, nil
}

// Retrive all services
func (p ServiceProcess) RetrieveAll() ([]ServiceProcessModel, error) {
	b := make([]byte, 1024*100)
	n, err := p.DB.Read(b)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err = json.Unmarshal(b[:n], &data); err != nil {
		return nil, err
	}

	names := reflect.ValueOf(data).MapKeys()
	processes := make([]ServiceProcessModel, len(names))
	for i, name := range names {
		var process ServiceProcessModel
		process.Name = name.String()

		// Get service data
		service := data[name.String()].(map[string]interface{})
		process.Status = service["status"].(bool)
		process.Type = service["type"].(string)
		process.Listeners = service["listeners"].([]interface{})
		process.RootDirectory = service["root"].(string)
		process.RemoteServer = service["server"].(string)
		process.Tunnelled = service["is_tunnelled"].(bool)
		process.CreatedAt = int64(service["created_at"].(float64))

		// Append
		processes[i] = process
	}
	return processes, nil
}

func (p ServiceProcess) Create(data map[string]interface{}) (ServiceProcessModel, error) {
	b := make([]byte, 1024*100)
	n, err := p.DB.Read(b)
	if err != nil {
		return ServiceProcessModel{}, err
	}

	records := make(map[string]interface{})
	if err = json.Unmarshal(b[:n], &records); err != nil {
		return ServiceProcessModel{}, err
	}

	// Add to running services
	name := data["name"].(string)
	records[name] = data

	recordBytes, err := json.Marshal(records)
	if err != nil {
		return ServiceProcessModel{}, err
	}
	// Rewrite file
	p.DB.Truncate(0)
	p.DB.Seek(0, 0)
	if _, err = p.DB.Write(recordBytes); err != nil {
		return ServiceProcessModel{}, nil
	}

	// Serialize the data and return
	var service ServiceProcessModel
	dataBytes, _ := json.Marshal(data)
	json.Unmarshal(dataBytes, &service)

	// Return result
	return service, nil
}

func (p ServiceProcess) Update(name string, data map[string]interface{}) (ServiceProcessModel, error) {
	// Delete name, and created_at from data. This are read-only fields
	for _, f := range []string{"name", "created_at"} {
		delete(data, f)
	}

	var process ServiceProcessModel

	// Read initial data
	b := make([]byte, 1024*100)
	n, err := p.DB.Read(b)
	if err != nil {
		return ServiceProcessModel{}, err
	}

	records := make(map[string]interface{})
	if err = json.Unmarshal(b[:n], &records); err != nil {
		return ServiceProcessModel{}, err
	}

	// Get the service from records
	if rec, ok := records[name]; ok {
		dataKeys := reflect.ValueOf(data).MapKeys()
		// Update all fields passed
		for _, k := range dataKeys {
			rec.(map[string]interface{})[k.String()] = data[k.String()]
		}

		// Update record
		records[name] = rec
		recordsByte, _ := json.Marshal(records)
		// Rewrite file
		p.DB.Truncate(0)
		p.DB.Seek(0, 0)
		if _, err = p.DB.Write(recordsByte); err != nil {
			return ServiceProcessModel{}, nil
		}
		// Unmarshall the service record and return
		recByte, _ := json.Marshal(rec)
		json.Unmarshal(recByte, &process)
		return process, nil
	}
	return ServiceProcessModel{}, fmt.Errorf("service unknwon: %s", name)
}

func (p ServiceProcess) Get(name string) (ServiceProcessModel, error) {
	b := make([]byte, 1024*100)
	n, err := p.DB.Read(b)
	if err != nil {
		return ServiceProcessModel{}, err
	}

	var records map[string]interface{}
	if err = json.Unmarshal(b[:n], &records); err != nil {
		return ServiceProcessModel{}, err
	}
	if rec, ok := records[name]; ok {
		var process ServiceProcessModel
		dataByte, _ := json.Marshal(rec)
		json.Unmarshal(dataByte, &process)
		return process, nil
	}
	return ServiceProcessModel{}, fmt.Errorf("service unknwon: %s", name)
}
