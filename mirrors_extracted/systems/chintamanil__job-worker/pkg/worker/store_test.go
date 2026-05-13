package worker

import (
	"sync"
	"testing"
)

func TestJobStore_AddAndGet(t *testing.T) {
	store := NewJobStore()
	job := &Job{
		ID:      "job-1",
		Command: "/bin/echo",
		Args:    []string{"hello"},
		Owner:   "user1",
		state:   JobState{Status: JobStatusRunning},
	}

	store.Add(job)

	// Owner can access their own job
	got, err := store.Get("user1", "job-1", false)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.ID != job.ID {
		t.Errorf("Get().ID = %s, want %s", got.ID, job.ID)
	}
	if got.Owner != job.Owner {
		t.Errorf("Get().Owner = %s, want %s", got.Owner, job.Owner)
	}
}

func TestJobStore_GetNonexistent(t *testing.T) {
	store := NewJobStore()

	_, err := store.Get("user1", "nonexistent", false)
	if err != ErrJobNotFound {
		t.Fatalf("Get() error = %v, want ErrJobNotFound", err)
	}
}

func TestJobStore_GetWithOwnership(t *testing.T) {
	store := NewJobStore()
	job := &Job{
		ID:    "job-1",
		Owner: "user1",
	}
	store.Add(job)

	tests := []struct {
		name    string
		userID  string
		jobID   string
		isAdmin bool
		wantErr bool
	}{
		{
			name:    "owner can access own job",
			userID:  "user1",
			jobID:   "job-1",
			isAdmin: false,
			wantErr: false,
		},
		{
			name:    "non-owner cannot access job",
			userID:  "user2",
			jobID:   "job-1",
			isAdmin: false,
			wantErr: true,
		},
		{
			name:    "admin can access any job",
			userID:  "admin",
			jobID:   "job-1",
			isAdmin: true,
			wantErr: false,
		},
		{
			name:    "nonexistent job returns error",
			userID:  "user1",
			jobID:   "nonexistent",
			isAdmin: false,
			wantErr: true,
		},
		{
			name:    "admin cannot access nonexistent job",
			userID:  "admin",
			jobID:   "nonexistent",
			isAdmin: true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.Get(tt.userID, tt.jobID, tt.isAdmin)
			if tt.wantErr {
				if err != ErrJobNotFound {
					t.Errorf("Get() error = %v, want ErrJobNotFound", err)
				}
				if got != nil {
					t.Errorf("Get() returned non-nil job for unauthorized access")
				}
			} else {
				if err != nil {
					t.Errorf("Get() unexpected error = %v", err)
				}
				if got == nil {
					t.Errorf("Get() returned nil for authorized access")
				}
			}
		})
	}
}

func TestJobStore_GetByID(t *testing.T) {
	store := NewJobStore()
	job := &Job{
		ID:    "job-1",
		Owner: "user1",
	}

	store.Add(job)

	// GetByID doesn't check ownership
	got, ok := store.GetByID("job-1")
	if !ok {
		t.Fatal("GetByID() returned false for existing job")
	}
	if got.ID != job.ID {
		t.Errorf("GetByID().ID = %s, want %s", got.ID, job.ID)
	}
}

func TestJobStore_GetByID_Nonexistent(t *testing.T) {
	store := NewJobStore()

	_, ok := store.GetByID("nonexistent")
	if ok {
		t.Fatal("GetByID() returned true for nonexistent job")
	}
}

func TestJobStore_Delete(t *testing.T) {
	store := NewJobStore()
	job := &Job{
		ID:    "job-1",
		Owner: "user1",
	}

	store.Add(job)
	store.Delete("job-1")

	_, ok := store.GetByID("job-1")
	if ok {
		t.Fatal("GetByID() returned true after Delete()")
	}
}

func TestJobStore_DeleteNonexistent(t *testing.T) {
	store := NewJobStore()
	// Should not panic
	store.Delete("nonexistent")
}

func TestJobStore_ConcurrentAccess(t *testing.T) {
	store := NewJobStore()
	var wg sync.WaitGroup

	// Concurrent writers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			job := &Job{
				ID:    string(rune('a'+id%26)) + string(rune('0'+id)),
				Owner: "user1",
			}
			store.Add(job)
		}(i)
	}

	// Concurrent readers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.Count()
		}()
	}

	wg.Wait()
}

func TestJobStore_Count(t *testing.T) {
	store := NewJobStore()

	if store.Count() != 0 {
		t.Fatalf("Count() = %d, want 0", store.Count())
	}

	store.Add(&Job{ID: "job-1", Owner: "user1"})
	if store.Count() != 1 {
		t.Fatalf("Count() = %d, want 1", store.Count())
	}

	store.Add(&Job{ID: "job-2", Owner: "user2"})
	if store.Count() != 2 {
		t.Fatalf("Count() = %d, want 2", store.Count())
	}

	store.Delete("job-1")
	if store.Count() != 1 {
		t.Fatalf("Count() = %d, want 1", store.Count())
	}
}

func TestJobStore_Update(t *testing.T) {
	store := NewJobStore()
	job := &Job{
		ID:    "job-1",
		Owner: "user1",
		state: JobState{Status: JobStatusRunning},
	}

	store.Add(job)

	// Modify the job and verify it's reflected (using lock since state is protected)
	got, _ := store.GetByID("job-1")
	got.mu.Lock()
	got.state.Status = JobStatusCompleted
	got.state.ExitCode = 0
	got.mu.Unlock()

	// Get again to verify update
	updated, _ := store.GetByID("job-1")
	state := updated.State()
	if state.Status != JobStatusCompleted {
		t.Errorf("Status = %v, want %v", state.Status, JobStatusCompleted)
	}
	if state.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", state.ExitCode)
	}
}

func TestJobStore_ResourceLimits(t *testing.T) {
	store := NewJobStore()
	job := &Job{
		ID:          "job-1",
		Owner:       "user1",
		CPULimit:    0.5,
		MemoryLimit: 1024 * 1024 * 512, // 512 MiB
		IOWeight:    500,               // Higher than default (100)
	}

	store.Add(job)

	got, err := store.Get("user1", "job-1", false)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.CPULimit != 0.5 {
		t.Errorf("CPULimit = %v, want 0.5", got.CPULimit)
	}
	if got.MemoryLimit != 1024*1024*512 {
		t.Errorf("MemoryLimit = %v, want %v", got.MemoryLimit, 1024*1024*512)
	}
	if got.IOWeight != 500 {
		t.Errorf("IOWeight = %v, want 500", got.IOWeight)
	}
}

func TestJobStore_AddDuplicateID(t *testing.T) {
	store := NewJobStore()

	job1 := &Job{
		ID:      "job-1",
		Owner:   "user1",
		Command: "/bin/first",
	}
	job2 := &Job{
		ID:      "job-1", // Same ID
		Owner:   "user2",
		Command: "/bin/second",
	}

	store.Add(job1)
	store.Add(job2) // Should replace

	got, ok := store.GetByID("job-1")
	if !ok {
		t.Fatal("GetByID() returned false")
	}

	// Second job should have replaced first
	if got.Owner != "user2" {
		t.Errorf("Owner = %s, want user2 (job replaced)", got.Owner)
	}
	if got.Command != "/bin/second" {
		t.Errorf("Command = %s, want /bin/second", got.Command)
	}

	// Count should still be 1
	if store.Count() != 1 {
		t.Errorf("Count() = %d, want 1", store.Count())
	}
}

func TestJobStore_ConcurrentGetAndDelete(t *testing.T) {
	store := NewJobStore()

	// Add jobs
	for i := range 100 {
		store.Add(&Job{
			ID:    string(rune('a' + i%26)),
			Owner: "user1",
		})
	}

	var wg sync.WaitGroup

	// Concurrent deleters
	for i := range 50 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			store.Delete(string(rune('a' + id%26)))
		}(i)
	}

	// Concurrent readers (with auth check)
	for i := range 50 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Get may succeed or fail - either is fine, just no panic/race
			store.Get("user1", string(rune('a'+id%26)), false)
		}(i)
	}

	wg.Wait()
}

func TestJobStore_SecurityNoEnumeration(t *testing.T) {
	// Verify that unauthorized access returns same error as nonexistent
	// This prevents attackers from enumerating job IDs
	store := NewJobStore()
	store.Add(&Job{
		ID:    "secret-job",
		Owner: "user1",
	})

	// Non-owner trying to access existing job
	_, errUnauthorized := store.Get("attacker", "secret-job", false)

	// Anyone trying to access nonexistent job
	_, errNonexistent := store.Get("attacker", "nonexistent", false)

	// Both should return same error (ErrJobNotFound)
	if errUnauthorized != ErrJobNotFound {
		t.Errorf("Unauthorized access error = %v, want ErrJobNotFound", errUnauthorized)
	}
	if errNonexistent != ErrJobNotFound {
		t.Errorf("Nonexistent job error = %v, want ErrJobNotFound", errNonexistent)
	}
	if errUnauthorized != errNonexistent {
		t.Error("Unauthorized and nonexistent should return identical errors")
	}
}

func TestJobStore_AllJobStatuses(t *testing.T) {
	store := NewJobStore()

	statuses := []JobStatus{
		JobStatusUnspecified,
		JobStatusRunning,
		JobStatusCompleted,
		JobStatusFailed,
		JobStatusStopped,
	}

	for i, status := range statuses {
		job := &Job{
			ID:    string(rune('a' + i)),
			Owner: "user1",
			state: JobState{Status: status},
		}
		store.Add(job)

		got, err := store.Get("user1", job.ID, false)
		if err != nil {
			t.Fatalf("Get() error = %v for status %v", err, status)
		}
		if got.State().Status != status {
			t.Errorf("Status = %v, want %v", got.State().Status, status)
		}
	}
}
