package main

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"math"
	"math/rand"
	"time"
)

type MineralType struct {
	Name          string
	Hardness      int
	MeltingPoint  float32
	FractureLimit int
}

func CreateMineralType(name string, hardness int, meltingPoint float32, fractureLimit int) MineralType {
	return MineralType{name, hardness, meltingPoint, fractureLimit}
}

type Mineral struct {
	UUID      uuid.UUID
	Type      MineralType
	State     string
	Fractures int
}


type Job struct {
	action  string
	mineral *Mineral
	Status  string
}

// Ready To be set by manager
func (job *Job) Ready() error {
	if job.Status != "NEW" {
		return errors.New("can on ready jobs with status NEW")
	}
	job.Status = "READY"
	return nil
}

// Started To by set by factory
func (job *Job) Started() error {
	if job.Status != "READY" {
		return errors.New("can on ready jobs with status READY")
	}
	job.Status = "STARTED"
	return nil
}

// Finished To be set by factory
func (job *Job) Finished() error {
	if job.Status != "STARTED" {
		return errors.New("can on ready jobs with status STARTED")
	}
	job.Status = "FINISHED"
	return nil
}

type JobQueue struct {
	jobs []*Job
	// Not implemented. Idea to maintain set of ids for faster access
	newJobIds      map[uuid.UUID]bool
	readyJobIds    map[uuid.UUID]bool
	finishedJobIds map[uuid.UUID]bool
}

func (jq *JobQueue) AddJob(job Job) {
	jq.jobs = append(jq.jobs, &job)
}

type Factory struct {
	jobQueue   *JobQueue
	currentJob Job
	active     bool
}

func (f Factory) FractureMineral(m *Mineral) error {
	log.Println("FractureMineral", m)
	if m.Fractures*2 > m.Type.FractureLimit {
		log.Println("Not able to fracture, over the limit")
		return errors.New("Problem")
	} else {
		m.Fractures *= 2
		return nil
	}
}

// Step checks if there are any new requests to process. Marks ready
func (f *Factory) Step() {
	if !f.active {
		log.Println("Factoring is offline")
		return
	}
	// Attempt to process a job

	log.Println("Factory. is making a step")
	for i, job := range f.jobQueue.jobs {
		log.Println(i, job)
		if job.Status == "READY" {
			log.Println("Factory. READY job found. Making Ready")
			err := job.Started()
			if err != nil {
				log.Println("ERROR", err)
				return
			}

			err = f.FractureMineral(job.mineral)
			if err != nil {
				log.Println("ERROR", err)
				return
			}

			log.Println("New status", job.Status)
		}
	}
}

type Manager struct {
	jobQueue *JobQueue
	active   bool
}

// Step checks if there are any new requests to process. Marks ready
func (m *Manager) Step() {
	log.Println("Manager. Step")
	log.Println("Manager. Iterating jobs")
	for i, job := range m.jobQueue.jobs {
		log.Println(i, job)
		if job.Status == "NEW" {
			log.Println("Manager. NEW job found. Making Ready")
			err := job.Ready()
			if err != nil {
				log.Println("ERROR", err)
			}
			log.Println("NEW STATUS", job.Status)
		}
	}

	log.Println("Manager. Checking new job requests")
}

type MineralTypeDB struct {
	mineralTypes map[string]MineralType
}

func (m *MineralTypeDB) AddMineralType(mineralType MineralType) {
	log.Println(mineralType, "CHECK")
	m.mineralTypes[mineralType.Name] = mineralType
}

func (m *MineralTypeDB) GetTypeByName(mineralName string) MineralType {
	log.Println("GetTypeByName")
	return m.mineralTypes[mineralName]
}

func (m MineralTypeDB) PrintAllMineralTypes() {
	log.Println(m.mineralTypes)
}

func CreateMineral(mineralTypeName MineralType, state string, fractures int) Mineral {
	return Mineral{
		UUID:      uuid.New(),
		Type:      mineralTypeName,
		State:     state,
		Fractures: fractures,
	}

}

func NewJob(action string, mineral *Mineral, status string) Job {
	return Job{
		action:  action,
		mineral: mineral,
		Status:  status,
	}

}

func main() {
	mineralTypeDB := MineralTypeDB{}
	mineralTypeDB.mineralTypes = make(map[string]MineralType)
	mineralTypeDB.AddMineralType(CreateMineralType("topaz", 200, 1000, 32))
	mineralTypeDB.AddMineralType(CreateMineralType("diamond", 1500, 5000, 8))

	jq := JobQueue{}
	supportedMineralTypeNames := []string{"diamond", "topaz"}

	for i := 0; i < 10; i++ {
		mineralTypeName := supportedMineralTypeNames[rand.Intn(len(supportedMineralTypeNames))]
		mineralType := mineralTypeDB.GetTypeByName(mineralTypeName)

		fractures := int(math.Pow(2, float64(rand.Intn(4)+1)))

		mineral := CreateMineral(mineralType, "fractured", fractures)
		job := NewJob("facture", &mineral, "NEW")
		jq.AddJob(job)
	}

	manager := &Manager{jobQueue: &jq}
	factory := &Factory{jobQueue: &jq, active: true}

	for {
		time.Sleep(time.Second)
		manager.Step()
		factory.Step()
	}
}
