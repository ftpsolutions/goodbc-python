package goodbc_python_go

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

var sessionMutex sync.Mutex
var sessions map[uint64]*session
var lastSessionID uint64
var rowsMutex sync.Mutex
var rows map[uint64]*sql.Rows
var lastRowsID uint64
var resultMutex sync.Mutex
var results map[uint64]*sql.Result
var lastResultID uint64

func init() {
	sessions = make(map[uint64]*session)
	rows = make(map[uint64]*sql.Rows)
	results = make(map[uint64]*sql.Result)

	time.Sleep(time.Second) // give the Python side a little time to settle
}

// this is used to ensure the Go runtime keeps operating in the event of strange errors
func handleSessionPanic(extra string, sessionID uint64, s *session, err error) {
	log.Printf(
		fmt.Sprintf(
			"handleSessionPanic() for %v()\n\tSessionID: %v\n\tSession: %+v\n\tError: %v\n\nStack trace follows:\n\n%v",
			extra,
			sessionID,
			s,
			err,
			string(debug.Stack()),
		),
	)
}

// this is used to ensure the Go runtime keeps operating in the event of strange errors
func handleRowsPanic(extra string, rowsID uint64, r *sql.Rows, err error) {
	log.Printf(
		fmt.Sprintf(
			"handleSessionPanic() for %v()\n\tRowsID: %v\n\tRows: %+v\n\tError: %v\n\nStack trace follows:\n\n%v",
			extra,
			rowsID,
			r,
			err,
			string(debug.Stack()),
		),
	)
}

// this is used to ensure the Go runtime keeps operating in the event of strange errors
func handleResultPanic(extra string, resultID uint64, r *sql.Result, err error) {
	log.Printf(
		fmt.Sprintf(
			"handleSessionPanic() for %v()\n\tResultID: %v\n\tResult: %+v\n\tError: %v\n\nStack trace follows:\n\n%v",
			extra,
			resultID,
			r,
			err,
			string(debug.Stack()),
		),
	)
}

// NewRPCSession creates a new session and returns the sessionID
func NewRPCSession(dataSourceName string) uint64 {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	session := newSession(
		dataSourceName,
	)

	sessionMutex.Lock()
	sessionID := lastSessionID
	lastSessionID++
	sessions[sessionID] = &session
	sessionMutex.Unlock()

	return sessionID
}

// RPCConnect connects
func RPCConnect(sessionID uint64) error {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error

	sessionMutex.Lock()
	val, ok := sessions[sessionID]
	sessionMutex.Unlock()

	// permit recovering from a panic but return the error
	defer func(s *session) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handleSessionPanic("connect", sessionID, val, handledError)
				err = handledError
			}
		}
	}(val)

	if !ok {
		return fmt.Errorf("sessionID %v does not exist", sessionID)
	}

	return val.connect()
}

// RPCQuery executes a query
func RPCQuery(sessionID uint64, query string) (uint64, error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error
	var rowsID uint64

	sessionMutex.Lock()
	val, ok := sessions[sessionID]
	sessionMutex.Unlock()

	if !ok {
		return rowsID, fmt.Errorf("sessionID %v does not exist", sessionID)
	}

	// permit recovering from a panic but return the error
	defer func(s *session) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handleSessionPanic("query", sessionID, val, handledError)
				err = handledError
			}
		}
	}(val)

	queryRows, err := val.query(query)
	if err != nil {
		return rowsID, fmt.Errorf("RPCQuery() query error: %v", err)
	}

	rowsMutex.Lock()
	rowsID = lastRowsID
	lastRowsID++
	rows[rowsID] = queryRows
	rowsMutex.Unlock()

	return rowsID, err
}

// RPCFetchAll returns results
func RPCFetchAll(sessionID, rowsID uint64) (string, error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error
	var records [][]multiField

	sessionMutex.Lock()
	sessionVal, ok := sessions[sessionID]
	sessionMutex.Unlock()

	if !ok {
		return "", fmt.Errorf("sessionID %v does not exist", sessionID)
	}

	// permit recovering from a panic but return the error
	defer func(r *session) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handleSessionPanic("fetchAll", sessionID, sessionVal, handledError)
				err = handledError
			}
		}
	}(sessionVal)

	rowsMutex.Lock()
	rowsVal, ok := rows[rowsID]
	delete(rows, rowsID)
	rowsMutex.Unlock()

	if !ok {
		return "", fmt.Errorf("rowsID %v does not exist", sessionID)
	}

	// permit recovering from a panic but return the error
	defer func(r *sql.Rows) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handleRowsPanic("fetchAll", rowsID, rowsVal, handledError)
				err = handledError
			}
		}
	}(rowsVal)

	records, err = sessionVal.fetchAll(rowsVal)
	if err != nil {
		return "", err
	}

	recordsJSON, err := json.Marshal(records)
	if err != nil {
		return "", err
	}

	return string(recordsJSON), err
}

// RPCExecute executes a query
func RPCExecute(sessionID uint64, query string) (uint64, error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error
	var resultID uint64

	sessionMutex.Lock()
	val, ok := sessions[sessionID]
	sessionMutex.Unlock()

	if !ok {
		return resultID, fmt.Errorf("sessionID %v does not exist", sessionID)
	}

	// permit recovering from a panic but return the error
	defer func(s *session) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handleSessionPanic("query", sessionID, val, handledError)
				err = handledError
			}
		}
	}(val)

	queryResult, err := val.execute(query)
	if err != nil {
		return resultID, err
	}

	rowsMutex.Lock()
	resultID = lastResultID
	lastResultID++
	results[resultID] = queryResult
	rowsMutex.Unlock()

	return resultID, err
}

// RPCGetRowsAffected returns the last inserted and affected row counts
func RPCGetRowsAffected(sessionID, resultID uint64) (int64, error) {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	var err error

	sessionMutex.Lock()
	sessionVal, ok := sessions[sessionID]
	sessionMutex.Unlock()

	if !ok {
		return 0, fmt.Errorf("sessionID %v does not exist", sessionID)
	}

	// permit recovering from a panic but return the error
	defer func(r *session) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handleSessionPanic("fetchAll", sessionID, sessionVal, handledError)
				err = handledError
			}
		}
	}(sessionVal)

	rowsMutex.Lock()
	resultVal, ok := results[resultID]
	delete(results, resultID)
	rowsMutex.Unlock()

	if !ok {
		return 0, fmt.Errorf("resultID %v does not exist", resultID)
	}

	// permit recovering from a panic but return the error
	defer func(r *sql.Result) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handleResultPanic("fetchAll", resultID, resultVal, handledError)
				err = handledError
			}
		}
	}(resultVal)

	r := *resultVal

	rowsAffected, err := r.RowsAffected()
	if err != nil {

	}

	return rowsAffected, err
}

// RPCClose closes
func RPCClose(sessionID uint64) error {
	if !GetPyPy() {
		tState := releaseGIL()
		defer reacquireGIL(tState)
	}

	sessionMutex.Lock()
	val, ok := sessions[sessionID]
	sessionMutex.Unlock()

	if !ok {
		return nil
	}

	sessionMutex.Lock()
	delete(sessions, sessionID)
	sessionMutex.Unlock()

	// permit recovering from a panic silently (bury the error)
	defer func(s *session) {
		if r := recover(); r != nil {
			if handledError, _ := r.(error); handledError != nil {
				handleSessionPanic("close", sessionID, val, handledError)
			}
		}
	}(val)

	return val.close()
}
