package leaderelection

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/lib/pq"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/forks/k8s.io/client-go/tools/leaderelection"
	"gitlab.com/flattrack/flattrack/internal/forks/k8s.io/client-go/tools/leaderelection/resourcelock"
)

type Lock struct {
	id   string
	name string
	db   *sql.DB

	lease *resourcelock.LeaderElectionRecord
}

func NewLock(db *sql.DB) *Lock {
	return &Lock{
		id:    common.RandStringRunes(5),
		name:  "default",
		db:    db,
		lease: nil,
	}
}

func (l *Lock) Get(ctx context.Context) (ler *resourcelock.LeaderElectionRecord, rb []byte, err error) {
	sqlStatement := `select holderIdentity, leaseDurationSeconds, acquireTime, renewTime, leaderTransitions from leader_election where name = $1 limit 1`
	rows, err := l.db.Query(sqlStatement, l.name)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		ler = &resourcelock.LeaderElectionRecord{}
		var acquireTime, renewTime int64
		if err = rows.Scan(&ler.HolderIdentity, &ler.LeaseDurationSeconds, &acquireTime, &renewTime, &ler.LeaderTransitions); err != nil {
			return nil, nil, err
		}
		ler.AcquireTime = metav1.NewTime(time.Unix(acquireTime, 0))
		ler.RenewTime = metav1.NewTime(time.Unix(renewTime, 0))
	}
	if ler == nil {
		l.lease = nil
		return &resourcelock.LeaderElectionRecord{}, nil, nil
	}
	l.lease = ler
	rb, err = json.Marshal(ler)
	if err != nil {
		return nil, nil, err
	}
	return ler, rb, nil
}

func (l *Lock) Create(ctx context.Context, ler resourcelock.LeaderElectionRecord) error {
	sqlStatement := `insert into leader_election (name, holderIdentity, leaseDurationSeconds, acquireTime, renewTime, leaderTransitions)
                     values ($1, $2, $3, $4, $5, $6)
                     returning holderIdentity, leaseDurationSeconds, acquireTime, renewTime, leaderTransitions`
	rows, err := l.db.Query(sqlStatement, l.name, ler.HolderIdentity, ler.LeaseDurationSeconds, ler.AcquireTime.Unix(), ler.RenewTime.Unix(), ler.LeaderTransitions)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error: failed to close rows: %v\n", err)
		}
	}()
	for rows.Next() {
		if l.lease == nil {
			l.lease = &resourcelock.LeaderElectionRecord{}
		}
		var acquireTime, renewTime int64
		if err = rows.Scan(&l.lease.HolderIdentity, &l.lease.LeaseDurationSeconds, &acquireTime, &renewTime, &l.lease.LeaderTransitions); err != nil {
			return err
		}
		l.lease.AcquireTime = metav1.NewTime(time.Unix(acquireTime, 0))
		l.lease.RenewTime = metav1.NewTime(time.Unix(renewTime, 0))
	}
	return nil
}

func (l *Lock) Update(ctx context.Context, ler resourcelock.LeaderElectionRecord) error {
	if l.lease == nil {
		return l.Create(ctx, resourcelock.LeaderElectionRecord{})
	}
	sqlStatement := `update leader_election
                     set holderIdentity = $2, leaseDurationSeconds = $3, acquireTime = $4, renewTime = $5, leaderTransitions = $6
                     where name = $1`
	if _, err := l.db.Exec(sqlStatement,
		l.name, ler.HolderIdentity, ler.LeaseDurationSeconds, ler.AcquireTime.Unix(), ler.RenewTime.Unix(), ler.LeaderTransitions); err != nil {
		return err
	}
	return nil
}

func (l *Lock) RecordEvent(s string) {
	log.Println("leader election event:", s)
}

func (l *Lock) Describe() string {
	return l.name
}

func (l *Lock) Identity() string {
	return l.id
}

func (l *Lock) Run(fn func() error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	leaderelector, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock:            l,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {},
			OnStoppedLeading: func() {
				log.Println("no longer the leader, staying inactive.")
			},
			OnNewLeader: func(currentID string) {
				if currentID == l.id {
					return
				}
				log.Printf("new/current leader is %s\n", currentID)
			},
		},
	})
	if err != nil {
		log.Fatalf("leader election error: %+v\n", err)
	}

	go func() {
		for {
			if leaderelector.IsLeader() {
				if err := fn(); err != nil {
					log.Println(err)
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()
	leaderelector.Run(ctx)
}
