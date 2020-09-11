package session

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/boundary/internal/db"
	"github.com/hashicorp/boundary/internal/iam"
	"github.com/hashicorp/boundary/internal/session/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestState_Create(t *testing.T) {
	t.Parallel()
	conn, _ := db.TestSetup(t, "postgres")
	wrapper := db.TestWrapper(t)
	iamRepo := iam.TestRepo(t, conn, wrapper)
	session := TestDefaultSession(t, conn, wrapper, iamRepo)

	type args struct {
		sessionId string
		status    Status
	}
	tests := []struct {
		name          string
		args          args
		want          *State
		wantErr       bool
		wantIsErr     error
		create        bool
		wantCreateErr bool
	}{
		{
			name: "valid",
			args: args{
				sessionId: session.PublicId,
				status:    StatusPending,
			},
			want: &State{
				State: &store.State{
					SessionId: session.PublicId,
					Status:    StatusPending.String(),
				},
			},
			create: true,
		},
		{
			name: "empty-sessionId",
			args: args{
				status: StatusPending,
			},
			wantErr:   true,
			wantIsErr: db.ErrInvalidParameter,
		},
		{
			name: "empty-status",
			args: args{
				sessionId: session.PublicId,
			},
			wantErr:   true,
			wantIsErr: db.ErrInvalidParameter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			got, err := NewState(tt.args.sessionId, tt.args.status)
			if tt.wantErr {
				require.Error(err)
				assert.True(errors.Is(err, tt.wantIsErr))
				return
			}
			require.NoError(err)
			assert.Equal(tt.want, got)
			if tt.create {
				err = db.New(conn).Create(context.Background(), got)
				if tt.wantCreateErr {
					assert.Error(err)
					return
				} else {
					assert.NoError(err)
				}
			}
		})
	}
}

func TestState_Delete(t *testing.T) {
	t.Parallel()
	conn, _ := db.TestSetup(t, "postgres")
	rw := db.New(conn)
	wrapper := db.TestWrapper(t)
	iamRepo := iam.TestRepo(t, conn, wrapper)
	session := TestDefaultSession(t, conn, wrapper, iamRepo)
	session2 := TestDefaultSession(t, conn, wrapper, iamRepo)

	tests := []struct {
		name            string
		state           *State
		deleteSessionId string
		wantRowsDeleted int
		wantErr         bool
		wantErrMsg      string
	}{
		{
			name:            "valid",
			state:           TestState(t, conn, session.PublicId, StatusTerminated),
			wantErr:         false,
			wantRowsDeleted: 1,
		},
		{
			name:  "bad-id",
			state: TestState(t, conn, session2.PublicId, StatusTerminated),
			deleteSessionId: func() string {
				id, err := db.NewPublicId(SessionPrefix)
				require.NoError(t, err)
				return id
			}(),
			wantErr:         false,
			wantRowsDeleted: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)

			var initialState State
			err := rw.LookupWhere(context.Background(), &initialState, "session_id = ? and state = ?", tt.state.SessionId, tt.state.Status)
			require.NoError(err)

			deleteState := allocState()
			if tt.deleteSessionId != "" {
				deleteState.SessionId = tt.deleteSessionId
			} else {
				deleteState.SessionId = tt.state.SessionId
			}
			deleteState.StartTime = initialState.StartTime
			deletedRows, err := rw.Delete(context.Background(), &deleteState)
			if tt.wantErr {
				require.Error(err)
				return
			}
			require.NoError(err)
			if tt.wantRowsDeleted == 0 {
				assert.Equal(tt.wantRowsDeleted, deletedRows)
				return
			}
			assert.Equal(tt.wantRowsDeleted, deletedRows)
			foundState := allocState()
			err = rw.LookupWhere(context.Background(), &foundState, "session_id = ? and start_time = ?", tt.state.SessionId, initialState.StartTime)
			fmt.Println(foundState)
			require.Error(err)
			assert.True(errors.Is(db.ErrRecordNotFound, err))
		})
	}
}

func TestState_Clone(t *testing.T) {
	t.Parallel()
	conn, _ := db.TestSetup(t, "postgres")
	wrapper := db.TestWrapper(t)
	iamRepo := iam.TestRepo(t, conn, wrapper)
	t.Run("valid", func(t *testing.T) {
		assert := assert.New(t)
		s := TestDefaultSession(t, conn, wrapper, iamRepo)
		state := TestState(t, conn, s.PublicId, StatusPending)
		cp := state.Clone()
		assert.True(proto.Equal(cp.(*State).State, state.State))
	})
	t.Run("not-equal", func(t *testing.T) {
		assert := assert.New(t)
		s := TestDefaultSession(t, conn, wrapper, iamRepo)
		s2 := TestDefaultSession(t, conn, wrapper, iamRepo)
		state := TestState(t, conn, s.PublicId, StatusPending)
		state2 := TestState(t, conn, s2.PublicId, StatusPending)

		cp := state.Clone()
		assert.True(!proto.Equal(cp.(*State).State, state2.State))
	})
}

func TestState_SetTableName(t *testing.T) {
	t.Parallel()
	defaultTableName := defaultStateTableName
	tests := []struct {
		name      string
		setNameTo string
		want      string
	}{
		{
			name:      "new-name",
			setNameTo: "new-name",
			want:      "new-name",
		},
		{
			name:      "reset to default",
			setNameTo: "",
			want:      defaultTableName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			def := allocState()
			require.Equal(defaultTableName, def.TableName())
			s := allocState()
			s.SetTableName(tt.setNameTo)
			assert.Equal(tt.want, s.TableName())
		})
	}
}
