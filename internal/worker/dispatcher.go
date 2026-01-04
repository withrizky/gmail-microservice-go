package worker

import (
	"encoding/json"
	"log"
	"sync"
	"sync/atomic"

	"gmail_microservice/internal/mailer"
	"gmail_microservice/internal/model"
)

type Dispatcher struct {
	JobQueue       chan model.EmailPayload
	MaxWorkers     int
	Accounts       []model.GmailAccount // Daftar Akun
	accountCounter uint32               // Penghitung untuk rotasi
	wg             sync.WaitGroup
}

// NewDispatcher menerima JSON string daftar akun
func NewDispatcher(maxWorkers int, bufferSize int, accountsJSON string) *Dispatcher {
	var accounts []model.GmailAccount

	// Parsing JSON Config dari .env ke Struct
	err := json.Unmarshal([]byte(accountsJSON), &accounts)
	if err != nil {
		log.Fatalf("Gagal load akun Gmail: %v", err)
	}

	if len(accounts) == 0 {
		log.Fatal("Tidak ada akun Gmail yang dikonfigurasi!")
	}

	return &Dispatcher{
		JobQueue:   make(chan model.EmailPayload, bufferSize),
		MaxWorkers: maxWorkers,
		Accounts:   accounts,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		d.wg.Add(1)
		go func(workerID int) {
			defer d.wg.Done()
			log.Printf("Worker #%d siap.", workerID)

			for job := range d.JobQueue {
				d.processJob(workerID, job)
			}

			log.Printf("Worker #%d berhenti.", workerID)
		}(i)
	}
}

func (d *Dispatcher) processJob(workerID int, job model.EmailPayload) {
	// --- LOGIKA ROTASI AKUN (ROUND ROBIN) ---
	// Ambil index selanjutnya, lalu di-modulo dengan jumlah akun
	idx := atomic.AddUint32(&d.accountCounter, 1) % uint32(len(d.Accounts))
	selectedAccount := d.Accounts[idx]

	// Kirim menggunakan akun yang terpilih
	err := mailer.Send(job, selectedAccount)

	if err != nil {
		log.Printf("[Worker-%d] ❌ Gagal (Sender: %s) -> %s: %v", workerID, selectedAccount.Email, job.To, err)
	} else {
		log.Printf("[Worker-%d] ✅ Sukses (Sender: %s) -> %s", workerID, selectedAccount.Email, job.To)
	}
}

func (d *Dispatcher) Stop() {
	close(d.JobQueue)
	d.wg.Wait()
}
