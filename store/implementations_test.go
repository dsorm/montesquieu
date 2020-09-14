package store_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/david-sorm/goblog/store"
	"github.com/david-sorm/goblog/store/postgres"
	"github.com/david-sorm/goblog/users"
	"github.com/jackc/pgx/v4"
	"github.com/ory/dockertest"
	"log"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

var storesToTest []store.Store

var storeConfig store.StoreConfig

type postgresDocker struct {
	postgresResource *dockertest.Resource
	pool             *dockertest.Pool
	id               string
}

type stores struct {
	currentStore int
	init         bool
}

// postgres_prepare makes sure there is a suitable environment for testing the
// postgres store
func (pd *postgresDocker) prepare() {

	// start the container
	cmd := exec.Command("bash", "-c",
		"docker run -e POSTGRES_USER=goblog -e POSTGRES_PASSWORD=goblog -e POSTGRES_DB=goblog --rm -d -p 5432:5432 postgres")
	outputBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("error while reading docker/bash's stdout:", err)
	}
	pd.id = string(outputBytes)
	// strings.TrimSuffix(string(outputBytes), "\n")
	fmt.Println("Postgres container ID: ", pd.id)

	// wait until it starts up
	err = errors.New("")
	for err != nil {
		_, err = pgx.Connect(context.Background(), "postgres://goblog:goblog@localhost:5432/goblog")
		if err != nil {
			fmt.Println("waiting until postgres container is up...")
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Postgres container running")
}

// stops the docker container
func (pd *postgresDocker) stop() {
	cmd := exec.Command("bash", "-c", "docker stop "+pd.id)
	if err := cmd.Start(); err != nil {
		fmt.Println("error while running bash script for docker container:", err)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("error while waiting for bash docker container script to execute:", err)
	}
	fmt.Println("Postgres container stopped")
}

// used in for loops for quick traversal of stores
func (s *stores) Next() bool {
	// make sure we begin with 0
	if !s.init {
		s.init = true
		s.currentStore = 0
		return true
	}

	if s.currentStore+1 >= len(storesToTest) {
		return false
	} else {
		s.currentStore++
		return true
	}
}

func (s *stores) Current() store.Store {
	return storesToTest[s.currentStore]
}

// overrides the default main from testing
func TestMain(m *testing.M) {

	// load all stores that are meant to be tested
	storesToTest = make([]store.Store, 0, 0)
	postgresStore := &postgres.Store{}
	storesToTest = append(storesToTest, postgresStore)

	// load mock config
	storeConfig = store.StoreConfig{
		Host:                 "127.0.0.1",
		Database:             "goblog",
		Username:             "goblog",
		Password:             "goblog",
		ArticlesPerIndexPage: 0,
	}

	// prepare all stores and their dependencies
	pd := postgresDocker{}
	pd.prepare()

	// run the tests
	exitCode := m.Run()

	pd.stop()

	// end
	os.Exit(exitCode)
}

func Test_Users(t *testing.T) {

	// just like reflect.DeepEqual but ignores ids (implementations might have different id's)
	//and password (they are not needed for ListUsers)
	usersEqual := func(a []users.User, b []users.User) bool {
		if len(a) != len(b) {
			return false
		}

		for k, v := range a {
			if (v.DisplayName == b[k].DisplayName) &&
				(v.Login == b[k].Login) {
				continue
			}
			return false
		}
		return true
	}

	checkGotWant := func(operation string, got []users.User, want []users.User, compareFunc func(a []users.User, b []users.User) bool) {
		if !compareFunc(got, want) {
			t.Errorf("%v; got: \n%#v\nwant: \n%#v\n\n", operation, got, want)
		}
	}

	strs := stores{}
	for strs.Next() {
		s := strs.Current()
		if err := s.Init(func() {}, storeConfig); err != nil {
			log.Fatal("Error has happened while initialising", s.Info().Name, "store:", err)
		}

		// add some mock users and check their passwords
		var passwordWant string
		var passwordGot string
		var login string
		var id uint64
		var exists bool
		for i := 1; i < 101; i++ {
			passwordWant = "neco" + strconv.Itoa(i)
			login = "nekdo" + strconv.Itoa(i)
			s.AddUser("Nekdo", login, passwordWant)
			id, exists = s.GetUserID(login)
			if !exists {
				t.Fatalf("User could not be found using GetUserID")
			}
			passwordGot = s.GetUser(id).Password
			if passwordWant != passwordGot {
				t.Errorf("password: got %#v, want %#v", passwordGot, passwordWant)
			}
		}

		// check if all users are ok first
		checkGotWant("ListUsers(0,100)",
			s.ListUsers(0, 100),
			[]users.User{{ID: 0x1, DisplayName: "Nekdo", Login: "nekdo1", Password: ""}, {ID: 0x2, DisplayName: "Nekdo", Login: "nekdo2", Password: ""}, {ID: 0x3, DisplayName: "Nekdo", Login: "nekdo3", Password: ""}, {ID: 0x4, DisplayName: "Nekdo", Login: "nekdo4", Password: ""}, {ID: 0x5, DisplayName: "Nekdo", Login: "nekdo5", Password: ""}, {ID: 0x6, DisplayName: "Nekdo", Login: "nekdo6", Password: ""}, {ID: 0x7, DisplayName: "Nekdo", Login: "nekdo7", Password: ""}, {ID: 0x8, DisplayName: "Nekdo", Login: "nekdo8", Password: ""}, {ID: 0x9, DisplayName: "Nekdo", Login: "nekdo9", Password: ""}, {ID: 0xa, DisplayName: "Nekdo", Login: "nekdo10", Password: ""}, {ID: 0xb, DisplayName: "Nekdo", Login: "nekdo11", Password: ""}, {ID: 0xc, DisplayName: "Nekdo", Login: "nekdo12", Password: ""}, {ID: 0xd, DisplayName: "Nekdo", Login: "nekdo13", Password: ""}, {ID: 0xe, DisplayName: "Nekdo", Login: "nekdo14", Password: ""}, {ID: 0xf, DisplayName: "Nekdo", Login: "nekdo15", Password: ""}, {ID: 0x10, DisplayName: "Nekdo", Login: "nekdo16", Password: ""}, {ID: 0x11, DisplayName: "Nekdo", Login: "nekdo17", Password: ""}, {ID: 0x12, DisplayName: "Nekdo", Login: "nekdo18", Password: ""}, {ID: 0x13, DisplayName: "Nekdo", Login: "nekdo19", Password: ""}, {ID: 0x14, DisplayName: "Nekdo", Login: "nekdo20", Password: ""}, {ID: 0x15, DisplayName: "Nekdo", Login: "nekdo21", Password: ""}, {ID: 0x16, DisplayName: "Nekdo", Login: "nekdo22", Password: ""}, {ID: 0x17, DisplayName: "Nekdo", Login: "nekdo23", Password: ""}, {ID: 0x18, DisplayName: "Nekdo", Login: "nekdo24", Password: ""}, {ID: 0x19, DisplayName: "Nekdo", Login: "nekdo25", Password: ""}, {ID: 0x1a, DisplayName: "Nekdo", Login: "nekdo26", Password: ""}, {ID: 0x1b, DisplayName: "Nekdo", Login: "nekdo27", Password: ""}, {ID: 0x1c, DisplayName: "Nekdo", Login: "nekdo28", Password: ""}, {ID: 0x1d, DisplayName: "Nekdo", Login: "nekdo29", Password: ""}, {ID: 0x1e, DisplayName: "Nekdo", Login: "nekdo30", Password: ""}, {ID: 0x1f, DisplayName: "Nekdo", Login: "nekdo31", Password: ""}, {ID: 0x20, DisplayName: "Nekdo", Login: "nekdo32", Password: ""}, {ID: 0x21, DisplayName: "Nekdo", Login: "nekdo33", Password: ""}, {ID: 0x22, DisplayName: "Nekdo", Login: "nekdo34", Password: ""}, {ID: 0x23, DisplayName: "Nekdo", Login: "nekdo35", Password: ""}, {ID: 0x24, DisplayName: "Nekdo", Login: "nekdo36", Password: ""}, {ID: 0x25, DisplayName: "Nekdo", Login: "nekdo37", Password: ""}, {ID: 0x26, DisplayName: "Nekdo", Login: "nekdo38", Password: ""}, {ID: 0x27, DisplayName: "Nekdo", Login: "nekdo39", Password: ""}, {ID: 0x28, DisplayName: "Nekdo", Login: "nekdo40", Password: ""}, {ID: 0x29, DisplayName: "Nekdo", Login: "nekdo41", Password: ""}, {ID: 0x2a, DisplayName: "Nekdo", Login: "nekdo42", Password: ""}, {ID: 0x2b, DisplayName: "Nekdo", Login: "nekdo43", Password: ""}, {ID: 0x2c, DisplayName: "Nekdo", Login: "nekdo44", Password: ""}, {ID: 0x2d, DisplayName: "Nekdo", Login: "nekdo45", Password: ""}, {ID: 0x2e, DisplayName: "Nekdo", Login: "nekdo46", Password: ""}, {ID: 0x2f, DisplayName: "Nekdo", Login: "nekdo47", Password: ""}, {ID: 0x30, DisplayName: "Nekdo", Login: "nekdo48", Password: ""}, {ID: 0x31, DisplayName: "Nekdo", Login: "nekdo49", Password: ""}, {ID: 0x32, DisplayName: "Nekdo", Login: "nekdo50", Password: ""}, {ID: 0x33, DisplayName: "Nekdo", Login: "nekdo51", Password: ""}, {ID: 0x34, DisplayName: "Nekdo", Login: "nekdo52", Password: ""}, {ID: 0x35, DisplayName: "Nekdo", Login: "nekdo53", Password: ""}, {ID: 0x36, DisplayName: "Nekdo", Login: "nekdo54", Password: ""}, {ID: 0x37, DisplayName: "Nekdo", Login: "nekdo55", Password: ""}, {ID: 0x38, DisplayName: "Nekdo", Login: "nekdo56", Password: ""}, {ID: 0x39, DisplayName: "Nekdo", Login: "nekdo57", Password: ""}, {ID: 0x3a, DisplayName: "Nekdo", Login: "nekdo58", Password: ""}, {ID: 0x3b, DisplayName: "Nekdo", Login: "nekdo59", Password: ""}, {ID: 0x3c, DisplayName: "Nekdo", Login: "nekdo60", Password: ""}, {ID: 0x3d, DisplayName: "Nekdo", Login: "nekdo61", Password: ""}, {ID: 0x3e, DisplayName: "Nekdo", Login: "nekdo62", Password: ""}, {ID: 0x3f, DisplayName: "Nekdo", Login: "nekdo63", Password: ""}, {ID: 0x40, DisplayName: "Nekdo", Login: "nekdo64", Password: ""}, {ID: 0x41, DisplayName: "Nekdo", Login: "nekdo65", Password: ""}, {ID: 0x42, DisplayName: "Nekdo", Login: "nekdo66", Password: ""}, {ID: 0x43, DisplayName: "Nekdo", Login: "nekdo67", Password: ""}, {ID: 0x44, DisplayName: "Nekdo", Login: "nekdo68", Password: ""}, {ID: 0x45, DisplayName: "Nekdo", Login: "nekdo69", Password: ""}, {ID: 0x46, DisplayName: "Nekdo", Login: "nekdo70", Password: ""}, {ID: 0x47, DisplayName: "Nekdo", Login: "nekdo71", Password: ""}, {ID: 0x48, DisplayName: "Nekdo", Login: "nekdo72", Password: ""}, {ID: 0x49, DisplayName: "Nekdo", Login: "nekdo73", Password: ""}, {ID: 0x4a, DisplayName: "Nekdo", Login: "nekdo74", Password: ""}, {ID: 0x4b, DisplayName: "Nekdo", Login: "nekdo75", Password: ""}, {ID: 0x4c, DisplayName: "Nekdo", Login: "nekdo76", Password: ""}, {ID: 0x4d, DisplayName: "Nekdo", Login: "nekdo77", Password: ""}, {ID: 0x4e, DisplayName: "Nekdo", Login: "nekdo78", Password: ""}, {ID: 0x4f, DisplayName: "Nekdo", Login: "nekdo79", Password: ""}, {ID: 0x50, DisplayName: "Nekdo", Login: "nekdo80", Password: ""}, {ID: 0x51, DisplayName: "Nekdo", Login: "nekdo81", Password: ""}, {ID: 0x52, DisplayName: "Nekdo", Login: "nekdo82", Password: ""}, {ID: 0x53, DisplayName: "Nekdo", Login: "nekdo83", Password: ""}, {ID: 0x54, DisplayName: "Nekdo", Login: "nekdo84", Password: ""}, {ID: 0x55, DisplayName: "Nekdo", Login: "nekdo85", Password: ""}, {ID: 0x56, DisplayName: "Nekdo", Login: "nekdo86", Password: ""}, {ID: 0x57, DisplayName: "Nekdo", Login: "nekdo87", Password: ""}, {ID: 0x58, DisplayName: "Nekdo", Login: "nekdo88", Password: ""}, {ID: 0x59, DisplayName: "Nekdo", Login: "nekdo89", Password: ""}, {ID: 0x5a, DisplayName: "Nekdo", Login: "nekdo90", Password: ""}, {ID: 0x5b, DisplayName: "Nekdo", Login: "nekdo91", Password: ""}, {ID: 0x5c, DisplayName: "Nekdo", Login: "nekdo92", Password: ""}, {ID: 0x5d, DisplayName: "Nekdo", Login: "nekdo93", Password: ""}, {ID: 0x5e, DisplayName: "Nekdo", Login: "nekdo94", Password: ""}, {ID: 0x5f, DisplayName: "Nekdo", Login: "nekdo95", Password: ""}, {ID: 0x60, DisplayName: "Nekdo", Login: "nekdo96", Password: ""}, {ID: 0x61, DisplayName: "Nekdo", Login: "nekdo97", Password: ""}, {ID: 0x62, DisplayName: "Nekdo", Login: "nekdo98", Password: ""}, {ID: 0x63, DisplayName: "Nekdo", Login: "nekdo99", Password: ""}, {ID: 0x64, DisplayName: "Nekdo", Login: "nekdo100", Password: ""}},
			usersEqual)

		// check if ranges are ok
		checkGotWant("ListUsers(22,74)",
			s.ListUsers(22, 74),
			[]users.User{users.User{ID: 0x17, DisplayName: "Nekdo", Login: "nekdo23", Password: ""}, users.User{ID: 0x18, DisplayName: "Nekdo", Login: "nekdo24", Password: ""}, users.User{ID: 0x19, DisplayName: "Nekdo", Login: "nekdo25", Password: ""}, users.User{ID: 0x1a, DisplayName: "Nekdo", Login: "nekdo26", Password: ""}, users.User{ID: 0x1b, DisplayName: "Nekdo", Login: "nekdo27", Password: ""}, users.User{ID: 0x1c, DisplayName: "Nekdo", Login: "nekdo28", Password: ""}, users.User{ID: 0x1d, DisplayName: "Nekdo", Login: "nekdo29", Password: ""}, users.User{ID: 0x1e, DisplayName: "Nekdo", Login: "nekdo30", Password: ""}, users.User{ID: 0x1f, DisplayName: "Nekdo", Login: "nekdo31", Password: ""}, users.User{ID: 0x20, DisplayName: "Nekdo", Login: "nekdo32", Password: ""}, users.User{ID: 0x21, DisplayName: "Nekdo", Login: "nekdo33", Password: ""}, users.User{ID: 0x22, DisplayName: "Nekdo", Login: "nekdo34", Password: ""}, users.User{ID: 0x23, DisplayName: "Nekdo", Login: "nekdo35", Password: ""}, users.User{ID: 0x24, DisplayName: "Nekdo", Login: "nekdo36", Password: ""}, users.User{ID: 0x25, DisplayName: "Nekdo", Login: "nekdo37", Password: ""}, users.User{ID: 0x26, DisplayName: "Nekdo", Login: "nekdo38", Password: ""}, users.User{ID: 0x27, DisplayName: "Nekdo", Login: "nekdo39", Password: ""}, users.User{ID: 0x28, DisplayName: "Nekdo", Login: "nekdo40", Password: ""}, users.User{ID: 0x29, DisplayName: "Nekdo", Login: "nekdo41", Password: ""}, users.User{ID: 0x2a, DisplayName: "Nekdo", Login: "nekdo42", Password: ""}, users.User{ID: 0x2b, DisplayName: "Nekdo", Login: "nekdo43", Password: ""}, users.User{ID: 0x2c, DisplayName: "Nekdo", Login: "nekdo44", Password: ""}, users.User{ID: 0x2d, DisplayName: "Nekdo", Login: "nekdo45", Password: ""}, users.User{ID: 0x2e, DisplayName: "Nekdo", Login: "nekdo46", Password: ""}, users.User{ID: 0x2f, DisplayName: "Nekdo", Login: "nekdo47", Password: ""}, users.User{ID: 0x30, DisplayName: "Nekdo", Login: "nekdo48", Password: ""}, users.User{ID: 0x31, DisplayName: "Nekdo", Login: "nekdo49", Password: ""}, users.User{ID: 0x32, DisplayName: "Nekdo", Login: "nekdo50", Password: ""}, users.User{ID: 0x33, DisplayName: "Nekdo", Login: "nekdo51", Password: ""}, users.User{ID: 0x34, DisplayName: "Nekdo", Login: "nekdo52", Password: ""}, users.User{ID: 0x35, DisplayName: "Nekdo", Login: "nekdo53", Password: ""}, users.User{ID: 0x36, DisplayName: "Nekdo", Login: "nekdo54", Password: ""}, users.User{ID: 0x37, DisplayName: "Nekdo", Login: "nekdo55", Password: ""}, users.User{ID: 0x38, DisplayName: "Nekdo", Login: "nekdo56", Password: ""}, users.User{ID: 0x39, DisplayName: "Nekdo", Login: "nekdo57", Password: ""}, users.User{ID: 0x3a, DisplayName: "Nekdo", Login: "nekdo58", Password: ""}, users.User{ID: 0x3b, DisplayName: "Nekdo", Login: "nekdo59", Password: ""}, users.User{ID: 0x3c, DisplayName: "Nekdo", Login: "nekdo60", Password: ""}, users.User{ID: 0x3d, DisplayName: "Nekdo", Login: "nekdo61", Password: ""}, users.User{ID: 0x3e, DisplayName: "Nekdo", Login: "nekdo62", Password: ""}, users.User{ID: 0x3f, DisplayName: "Nekdo", Login: "nekdo63", Password: ""}, users.User{ID: 0x40, DisplayName: "Nekdo", Login: "nekdo64", Password: ""}, users.User{ID: 0x41, DisplayName: "Nekdo", Login: "nekdo65", Password: ""}, users.User{ID: 0x42, DisplayName: "Nekdo", Login: "nekdo66", Password: ""}, users.User{ID: 0x43, DisplayName: "Nekdo", Login: "nekdo67", Password: ""}, users.User{ID: 0x44, DisplayName: "Nekdo", Login: "nekdo68", Password: ""}, users.User{ID: 0x45, DisplayName: "Nekdo", Login: "nekdo69", Password: ""}, users.User{ID: 0x46, DisplayName: "Nekdo", Login: "nekdo70", Password: ""}, users.User{ID: 0x47, DisplayName: "Nekdo", Login: "nekdo71", Password: ""}, users.User{ID: 0x48, DisplayName: "Nekdo", Login: "nekdo72", Password: ""}, users.User{ID: 0x49, DisplayName: "Nekdo", Login: "nekdo73", Password: ""}, users.User{ID: 0x4a, DisplayName: "Nekdo", Login: "nekdo74", Password: ""}, users.User{ID: 0x4b, DisplayName: "Nekdo", Login: "nekdo75", Password: ""}, users.User{ID: 0x4c, DisplayName: "Nekdo", Login: "nekdo76", Password: ""}, users.User{ID: 0x4d, DisplayName: "Nekdo", Login: "nekdo77", Password: ""}, users.User{ID: 0x4e, DisplayName: "Nekdo", Login: "nekdo78", Password: ""}, users.User{ID: 0x4f, DisplayName: "Nekdo", Login: "nekdo79", Password: ""}, users.User{ID: 0x50, DisplayName: "Nekdo", Login: "nekdo80", Password: ""}, users.User{ID: 0x51, DisplayName: "Nekdo", Login: "nekdo81", Password: ""}, users.User{ID: 0x52, DisplayName: "Nekdo", Login: "nekdo82", Password: ""}, users.User{ID: 0x53, DisplayName: "Nekdo", Login: "nekdo83", Password: ""}, users.User{ID: 0x54, DisplayName: "Nekdo", Login: "nekdo84", Password: ""}, users.User{ID: 0x55, DisplayName: "Nekdo", Login: "nekdo85", Password: ""}, users.User{ID: 0x56, DisplayName: "Nekdo", Login: "nekdo86", Password: ""}, users.User{ID: 0x57, DisplayName: "Nekdo", Login: "nekdo87", Password: ""}, users.User{ID: 0x58, DisplayName: "Nekdo", Login: "nekdo88", Password: ""}, users.User{ID: 0x59, DisplayName: "Nekdo", Login: "nekdo89", Password: ""}, users.User{ID: 0x5a, DisplayName: "Nekdo", Login: "nekdo90", Password: ""}, users.User{ID: 0x5b, DisplayName: "Nekdo", Login: "nekdo91", Password: ""}, users.User{ID: 0x5c, DisplayName: "Nekdo", Login: "nekdo92", Password: ""}, users.User{ID: 0x5d, DisplayName: "Nekdo", Login: "nekdo93", Password: ""}, users.User{ID: 0x5e, DisplayName: "Nekdo", Login: "nekdo94", Password: ""}, users.User{ID: 0x5f, DisplayName: "Nekdo", Login: "nekdo95", Password: ""}, users.User{ID: 0x60, DisplayName: "Nekdo", Login: "nekdo96", Password: ""}},
			usersEqual)

		checkGotWant("ListUsers(31,5)",
			s.ListUsers(31, 5),
			[]users.User{users.User{ID: 0x20, DisplayName: "Nekdo", Login: "nekdo32", Password: ""}, users.User{ID: 0x21, DisplayName: "Nekdo", Login: "nekdo33", Password: ""}, users.User{ID: 0x22, DisplayName: "Nekdo", Login: "nekdo34", Password: ""}, users.User{ID: 0x23, DisplayName: "Nekdo", Login: "nekdo35", Password: ""}, users.User{ID: 0x24, DisplayName: "Nekdo", Login: "nekdo36", Password: ""}},
			usersEqual)

		checkGotWant("ListUsers(98, 120)",
			s.ListUsers(98, 120),
			[]users.User{users.User{ID: 0x63, DisplayName: "Nekdo", Login: "nekdo99", Password: ""}, users.User{ID: 0x64, DisplayName: "Nekdo", Login: "nekdo100", Password: ""}},
			usersEqual)

		// change exactly half of the mock users
		var str string
		for i := 1; i < 101; i += 2 {
			str = strconv.Itoa(i)
			id, exists = s.GetUserID("nekdo" + str)
			if !exists {
				t.Fatalf("User could not be found using GetUserID")
			}
			s.EditUser(users.User{
				ID:          id,
				DisplayName: "Nekdo++" + str,
				Login:       "nekdo++" + str,
				Password:    "neco" + str,
			})
		}

		checkGotWant("ListUsers(0,130)",
			s.ListUsers(0, 130),
			[]users.User{users.User{ID: 0x1, DisplayName: "Nekdo++1", Login: "nekdo++1", Password: ""}, users.User{ID: 0x2, DisplayName: "Nekdo", Login: "nekdo2", Password: ""}, users.User{ID: 0x3, DisplayName: "Nekdo++3", Login: "nekdo++3", Password: ""}, users.User{ID: 0x4, DisplayName: "Nekdo", Login: "nekdo4", Password: ""}, users.User{ID: 0x5, DisplayName: "Nekdo++5", Login: "nekdo++5", Password: ""}, users.User{ID: 0x6, DisplayName: "Nekdo", Login: "nekdo6", Password: ""}, users.User{ID: 0x7, DisplayName: "Nekdo++7", Login: "nekdo++7", Password: ""}, users.User{ID: 0x8, DisplayName: "Nekdo", Login: "nekdo8", Password: ""}, users.User{ID: 0x9, DisplayName: "Nekdo++9", Login: "nekdo++9", Password: ""}, users.User{ID: 0xa, DisplayName: "Nekdo", Login: "nekdo10", Password: ""}, users.User{ID: 0xb, DisplayName: "Nekdo++11", Login: "nekdo++11", Password: ""}, users.User{ID: 0xc, DisplayName: "Nekdo", Login: "nekdo12", Password: ""}, users.User{ID: 0xd, DisplayName: "Nekdo++13", Login: "nekdo++13", Password: ""}, users.User{ID: 0xe, DisplayName: "Nekdo", Login: "nekdo14", Password: ""}, users.User{ID: 0xf, DisplayName: "Nekdo++15", Login: "nekdo++15", Password: ""}, users.User{ID: 0x10, DisplayName: "Nekdo", Login: "nekdo16", Password: ""}, users.User{ID: 0x11, DisplayName: "Nekdo++17", Login: "nekdo++17", Password: ""}, users.User{ID: 0x12, DisplayName: "Nekdo", Login: "nekdo18", Password: ""}, users.User{ID: 0x13, DisplayName: "Nekdo++19", Login: "nekdo++19", Password: ""}, users.User{ID: 0x14, DisplayName: "Nekdo", Login: "nekdo20", Password: ""}, users.User{ID: 0x15, DisplayName: "Nekdo++21", Login: "nekdo++21", Password: ""}, users.User{ID: 0x16, DisplayName: "Nekdo", Login: "nekdo22", Password: ""}, users.User{ID: 0x17, DisplayName: "Nekdo++23", Login: "nekdo++23", Password: ""}, users.User{ID: 0x18, DisplayName: "Nekdo", Login: "nekdo24", Password: ""}, users.User{ID: 0x19, DisplayName: "Nekdo++25", Login: "nekdo++25", Password: ""}, users.User{ID: 0x1a, DisplayName: "Nekdo", Login: "nekdo26", Password: ""}, users.User{ID: 0x1b, DisplayName: "Nekdo++27", Login: "nekdo++27", Password: ""}, users.User{ID: 0x1c, DisplayName: "Nekdo", Login: "nekdo28", Password: ""}, users.User{ID: 0x1d, DisplayName: "Nekdo++29", Login: "nekdo++29", Password: ""}, users.User{ID: 0x1e, DisplayName: "Nekdo", Login: "nekdo30", Password: ""}, users.User{ID: 0x1f, DisplayName: "Nekdo++31", Login: "nekdo++31", Password: ""}, users.User{ID: 0x20, DisplayName: "Nekdo", Login: "nekdo32", Password: ""}, users.User{ID: 0x21, DisplayName: "Nekdo++33", Login: "nekdo++33", Password: ""}, users.User{ID: 0x22, DisplayName: "Nekdo", Login: "nekdo34", Password: ""}, users.User{ID: 0x23, DisplayName: "Nekdo++35", Login: "nekdo++35", Password: ""}, users.User{ID: 0x24, DisplayName: "Nekdo", Login: "nekdo36", Password: ""}, users.User{ID: 0x25, DisplayName: "Nekdo++37", Login: "nekdo++37", Password: ""}, users.User{ID: 0x26, DisplayName: "Nekdo", Login: "nekdo38", Password: ""}, users.User{ID: 0x27, DisplayName: "Nekdo++39", Login: "nekdo++39", Password: ""}, users.User{ID: 0x28, DisplayName: "Nekdo", Login: "nekdo40", Password: ""}, users.User{ID: 0x29, DisplayName: "Nekdo++41", Login: "nekdo++41", Password: ""}, users.User{ID: 0x2a, DisplayName: "Nekdo", Login: "nekdo42", Password: ""}, users.User{ID: 0x2b, DisplayName: "Nekdo++43", Login: "nekdo++43", Password: ""}, users.User{ID: 0x2c, DisplayName: "Nekdo", Login: "nekdo44", Password: ""}, users.User{ID: 0x2d, DisplayName: "Nekdo++45", Login: "nekdo++45", Password: ""}, users.User{ID: 0x2e, DisplayName: "Nekdo", Login: "nekdo46", Password: ""}, users.User{ID: 0x2f, DisplayName: "Nekdo++47", Login: "nekdo++47", Password: ""}, users.User{ID: 0x30, DisplayName: "Nekdo", Login: "nekdo48", Password: ""}, users.User{ID: 0x31, DisplayName: "Nekdo++49", Login: "nekdo++49", Password: ""}, users.User{ID: 0x32, DisplayName: "Nekdo", Login: "nekdo50", Password: ""}, users.User{ID: 0x33, DisplayName: "Nekdo++51", Login: "nekdo++51", Password: ""}, users.User{ID: 0x34, DisplayName: "Nekdo", Login: "nekdo52", Password: ""}, users.User{ID: 0x35, DisplayName: "Nekdo++53", Login: "nekdo++53", Password: ""}, users.User{ID: 0x36, DisplayName: "Nekdo", Login: "nekdo54", Password: ""}, users.User{ID: 0x37, DisplayName: "Nekdo++55", Login: "nekdo++55", Password: ""}, users.User{ID: 0x38, DisplayName: "Nekdo", Login: "nekdo56", Password: ""}, users.User{ID: 0x39, DisplayName: "Nekdo++57", Login: "nekdo++57", Password: ""}, users.User{ID: 0x3a, DisplayName: "Nekdo", Login: "nekdo58", Password: ""}, users.User{ID: 0x3b, DisplayName: "Nekdo++59", Login: "nekdo++59", Password: ""}, users.User{ID: 0x3c, DisplayName: "Nekdo", Login: "nekdo60", Password: ""}, users.User{ID: 0x3d, DisplayName: "Nekdo++61", Login: "nekdo++61", Password: ""}, users.User{ID: 0x3e, DisplayName: "Nekdo", Login: "nekdo62", Password: ""}, users.User{ID: 0x3f, DisplayName: "Nekdo++63", Login: "nekdo++63", Password: ""}, users.User{ID: 0x40, DisplayName: "Nekdo", Login: "nekdo64", Password: ""}, users.User{ID: 0x41, DisplayName: "Nekdo++65", Login: "nekdo++65", Password: ""}, users.User{ID: 0x42, DisplayName: "Nekdo", Login: "nekdo66", Password: ""}, users.User{ID: 0x43, DisplayName: "Nekdo++67", Login: "nekdo++67", Password: ""}, users.User{ID: 0x44, DisplayName: "Nekdo", Login: "nekdo68", Password: ""}, users.User{ID: 0x45, DisplayName: "Nekdo++69", Login: "nekdo++69", Password: ""}, users.User{ID: 0x46, DisplayName: "Nekdo", Login: "nekdo70", Password: ""}, users.User{ID: 0x47, DisplayName: "Nekdo++71", Login: "nekdo++71", Password: ""}, users.User{ID: 0x48, DisplayName: "Nekdo", Login: "nekdo72", Password: ""}, users.User{ID: 0x49, DisplayName: "Nekdo++73", Login: "nekdo++73", Password: ""}, users.User{ID: 0x4a, DisplayName: "Nekdo", Login: "nekdo74", Password: ""}, users.User{ID: 0x4b, DisplayName: "Nekdo++75", Login: "nekdo++75", Password: ""}, users.User{ID: 0x4c, DisplayName: "Nekdo", Login: "nekdo76", Password: ""}, users.User{ID: 0x4d, DisplayName: "Nekdo++77", Login: "nekdo++77", Password: ""}, users.User{ID: 0x4e, DisplayName: "Nekdo", Login: "nekdo78", Password: ""}, users.User{ID: 0x4f, DisplayName: "Nekdo++79", Login: "nekdo++79", Password: ""}, users.User{ID: 0x50, DisplayName: "Nekdo", Login: "nekdo80", Password: ""}, users.User{ID: 0x51, DisplayName: "Nekdo++81", Login: "nekdo++81", Password: ""}, users.User{ID: 0x52, DisplayName: "Nekdo", Login: "nekdo82", Password: ""}, users.User{ID: 0x53, DisplayName: "Nekdo++83", Login: "nekdo++83", Password: ""}, users.User{ID: 0x54, DisplayName: "Nekdo", Login: "nekdo84", Password: ""}, users.User{ID: 0x55, DisplayName: "Nekdo++85", Login: "nekdo++85", Password: ""}, users.User{ID: 0x56, DisplayName: "Nekdo", Login: "nekdo86", Password: ""}, users.User{ID: 0x57, DisplayName: "Nekdo++87", Login: "nekdo++87", Password: ""}, users.User{ID: 0x58, DisplayName: "Nekdo", Login: "nekdo88", Password: ""}, users.User{ID: 0x59, DisplayName: "Nekdo++89", Login: "nekdo++89", Password: ""}, users.User{ID: 0x5a, DisplayName: "Nekdo", Login: "nekdo90", Password: ""}, users.User{ID: 0x5b, DisplayName: "Nekdo++91", Login: "nekdo++91", Password: ""}, users.User{ID: 0x5c, DisplayName: "Nekdo", Login: "nekdo92", Password: ""}, users.User{ID: 0x5d, DisplayName: "Nekdo++93", Login: "nekdo++93", Password: ""}, users.User{ID: 0x5e, DisplayName: "Nekdo", Login: "nekdo94", Password: ""}, users.User{ID: 0x5f, DisplayName: "Nekdo++95", Login: "nekdo++95", Password: ""}, users.User{ID: 0x60, DisplayName: "Nekdo", Login: "nekdo96", Password: ""}, users.User{ID: 0x61, DisplayName: "Nekdo++97", Login: "nekdo++97", Password: ""}, users.User{ID: 0x62, DisplayName: "Nekdo", Login: "nekdo98", Password: ""}, users.User{ID: 0x63, DisplayName: "Nekdo++99", Login: "nekdo++99", Password: ""}, users.User{ID: 0x64, DisplayName: "Nekdo", Login: "nekdo100", Password: ""}},
			usersEqual)

		// remove the other half
		for i := 2; i < 101; i += 2 {
			str = strconv.Itoa(i)
			id, exists = s.GetUserID("nekdo" + str)
			if !exists {
				t.Fatalf("User could not be found using GetUserID")
			}
			s.RemoveUser(id)
		}

		// check if the other half was deleted
		checkGotWant("ListUsers(0,120)",
			s.ListUsers(0, 120),
			[]users.User{users.User{ID: 0x1, DisplayName: "Nekdo++1", Login: "nekdo++1", Password: ""}, users.User{ID: 0x3, DisplayName: "Nekdo++3", Login: "nekdo++3", Password: ""}, users.User{ID: 0x5, DisplayName: "Nekdo++5", Login: "nekdo++5", Password: ""}, users.User{ID: 0x7, DisplayName: "Nekdo++7", Login: "nekdo++7", Password: ""}, users.User{ID: 0x9, DisplayName: "Nekdo++9", Login: "nekdo++9", Password: ""}, users.User{ID: 0xb, DisplayName: "Nekdo++11", Login: "nekdo++11", Password: ""}, users.User{ID: 0xd, DisplayName: "Nekdo++13", Login: "nekdo++13", Password: ""}, users.User{ID: 0xf, DisplayName: "Nekdo++15", Login: "nekdo++15", Password: ""}, users.User{ID: 0x11, DisplayName: "Nekdo++17", Login: "nekdo++17", Password: ""}, users.User{ID: 0x13, DisplayName: "Nekdo++19", Login: "nekdo++19", Password: ""}, users.User{ID: 0x15, DisplayName: "Nekdo++21", Login: "nekdo++21", Password: ""}, users.User{ID: 0x17, DisplayName: "Nekdo++23", Login: "nekdo++23", Password: ""}, users.User{ID: 0x19, DisplayName: "Nekdo++25", Login: "nekdo++25", Password: ""}, users.User{ID: 0x1b, DisplayName: "Nekdo++27", Login: "nekdo++27", Password: ""}, users.User{ID: 0x1d, DisplayName: "Nekdo++29", Login: "nekdo++29", Password: ""}, users.User{ID: 0x1f, DisplayName: "Nekdo++31", Login: "nekdo++31", Password: ""}, users.User{ID: 0x21, DisplayName: "Nekdo++33", Login: "nekdo++33", Password: ""}, users.User{ID: 0x23, DisplayName: "Nekdo++35", Login: "nekdo++35", Password: ""}, users.User{ID: 0x25, DisplayName: "Nekdo++37", Login: "nekdo++37", Password: ""}, users.User{ID: 0x27, DisplayName: "Nekdo++39", Login: "nekdo++39", Password: ""}, users.User{ID: 0x29, DisplayName: "Nekdo++41", Login: "nekdo++41", Password: ""}, users.User{ID: 0x2b, DisplayName: "Nekdo++43", Login: "nekdo++43", Password: ""}, users.User{ID: 0x2d, DisplayName: "Nekdo++45", Login: "nekdo++45", Password: ""}, users.User{ID: 0x2f, DisplayName: "Nekdo++47", Login: "nekdo++47", Password: ""}, users.User{ID: 0x31, DisplayName: "Nekdo++49", Login: "nekdo++49", Password: ""}, users.User{ID: 0x33, DisplayName: "Nekdo++51", Login: "nekdo++51", Password: ""}, users.User{ID: 0x35, DisplayName: "Nekdo++53", Login: "nekdo++53", Password: ""}, users.User{ID: 0x37, DisplayName: "Nekdo++55", Login: "nekdo++55", Password: ""}, users.User{ID: 0x39, DisplayName: "Nekdo++57", Login: "nekdo++57", Password: ""}, users.User{ID: 0x3b, DisplayName: "Nekdo++59", Login: "nekdo++59", Password: ""}, users.User{ID: 0x3d, DisplayName: "Nekdo++61", Login: "nekdo++61", Password: ""}, users.User{ID: 0x3f, DisplayName: "Nekdo++63", Login: "nekdo++63", Password: ""}, users.User{ID: 0x41, DisplayName: "Nekdo++65", Login: "nekdo++65", Password: ""}, users.User{ID: 0x43, DisplayName: "Nekdo++67", Login: "nekdo++67", Password: ""}, users.User{ID: 0x45, DisplayName: "Nekdo++69", Login: "nekdo++69", Password: ""}, users.User{ID: 0x47, DisplayName: "Nekdo++71", Login: "nekdo++71", Password: ""}, users.User{ID: 0x49, DisplayName: "Nekdo++73", Login: "nekdo++73", Password: ""}, users.User{ID: 0x4b, DisplayName: "Nekdo++75", Login: "nekdo++75", Password: ""}, users.User{ID: 0x4d, DisplayName: "Nekdo++77", Login: "nekdo++77", Password: ""}, users.User{ID: 0x4f, DisplayName: "Nekdo++79", Login: "nekdo++79", Password: ""}, users.User{ID: 0x51, DisplayName: "Nekdo++81", Login: "nekdo++81", Password: ""}, users.User{ID: 0x53, DisplayName: "Nekdo++83", Login: "nekdo++83", Password: ""}, users.User{ID: 0x55, DisplayName: "Nekdo++85", Login: "nekdo++85", Password: ""}, users.User{ID: 0x57, DisplayName: "Nekdo++87", Login: "nekdo++87", Password: ""}, users.User{ID: 0x59, DisplayName: "Nekdo++89", Login: "nekdo++89", Password: ""}, users.User{ID: 0x5b, DisplayName: "Nekdo++91", Login: "nekdo++91", Password: ""}, users.User{ID: 0x5d, DisplayName: "Nekdo++93", Login: "nekdo++93", Password: ""}, users.User{ID: 0x5f, DisplayName: "Nekdo++95", Login: "nekdo++95", Password: ""}, users.User{ID: 0x61, DisplayName: "Nekdo++97", Login: "nekdo++97", Password: ""}, users.User{ID: 0x63, DisplayName: "Nekdo++99", Login: "nekdo++99", Password: ""}},
			usersEqual)

		// delete all users
		usrs := s.ListUsers(0, 101)
		for _, v := range usrs {
			s.RemoveUser(v.ID)
		}
	}
}

func Test_Authors(t *testing.T) {
	strs := stores{}

	authorsEqual := func(a []users.Author, b []users.Author) bool {
		if len(a) != len(b) {
			return false
		}

		for k, v := range a {
			if (v.DisplayName == b[k].DisplayName) &&
				(v.Login == b[k].Login) &&
				(v.AuthorName == b[k].AuthorName) {
				continue
			}
			return false
		}
		return true
	}

	checkGotWant := func(operation string, got []users.Author, want []users.Author, compareFunc func(a []users.Author, b []users.Author) bool) {
		if !compareFunc(got, want) {
			t.Errorf("%v; got: \n%#v\nwant: \n%#v\n\n", operation, got, want)
		}
	}
	for strs.Next() {
		s := strs.Current()

		// add a few users
		s.AddUser("Nekdo 1", "nekdo1", "")
		s.AddUser("Nekdo 2", "nekdo2", "")
		s.AddUser("Nekdo 3", "nekdo3", "")

		// make user 1 and user 2 authors
		id, _ := s.GetUserID("nekdo1")
		s.AddAuthor(id, "Nekdo Author 1")

		id, _ = s.GetUserID("nekdo2")
		s.AddAuthor(id, "Nekdo Author 2")

		got := s.ListAuthors(0, 100)
		want := []users.Author{users.Author{User: users.User{ID: 0x65, DisplayName: "Nekdo 1", Login: "nekdo1", Password: ""}, AuthorID: 0x1, AuthorName: "Nekdo Author 1"}, users.Author{User: users.User{ID: 0x66, DisplayName: "Nekdo 2", Login: "nekdo2", Password: ""}, AuthorID: 0x2, AuthorName: "Nekdo Author 2"}}
		checkGotWant("ListAuthors(0, 100)", got, want, authorsEqual)

		uid, _ := s.GetUserID("nekdo1")
		got2 := s.GetAuthor(uid)

		if !(got2.AuthorName == "Nekdo Author 1") {
			t.Errorf("GetAuthor(id of nekdo1); wanted author.AuthorName == `Nekdo Author 1, got author.AuthorName = %v", got2.AuthorName)
		}

		// relink author 2 to nekdo3 from nekdo2
		uid, _ = s.GetUserID("nekdo2")
		uid2, _ := s.GetUserID("nekdo3")
		a := s.GetAuthor(uid)
		s.LinkAuthor(a.AuthorID, uid2)

		// add another author
		s.AddAuthor(uid, "nekdo 2")

		got = s.ListAuthors(2, 1)
		want = []users.Author{users.Author{User: users.User{ID: 0x67, DisplayName: "Nekdo 3", Login: "nekdo3", Password: ""}, AuthorID: 0x2, AuthorName: "Nekdo Author 2"}}
		checkGotWant("ListAuthors(1,1)", got, want, authorsEqual)

		got = s.ListAuthors(0, 3)
		want = []users.Author{users.Author{User: users.User{ID: 0x65, DisplayName: "Nekdo 1", Login: "nekdo1", Password: ""}, AuthorID: 0x1, AuthorName: "Nekdo Author 1"}, users.Author{User: users.User{ID: 0x66, DisplayName: "Nekdo 2", Login: "nekdo2", Password: ""}, AuthorID: 0x3, AuthorName: "nekdo 2"}, users.Author{User: users.User{ID: 0x67, DisplayName: "Nekdo 3", Login: "nekdo3", Password: ""}, AuthorID: 0x2, AuthorName: "Nekdo Author 2"}}
		checkGotWant("ListAuthors(0,3)", got, want, authorsEqual)

		uid, _ = s.GetUserID("nekdo1")
		aid := s.GetAuthor(uid)
		s.RemoveAuthor(aid.AuthorID)

		got = s.ListAuthors(0, 6)
		want = []users.Author{users.Author{User: users.User{ID: 0x66, DisplayName: "Nekdo 2", Login: "nekdo2", Password: ""}, AuthorID: 0x3, AuthorName: "nekdo 2"}, users.Author{User: users.User{ID: 0x67, DisplayName: "Nekdo 3", Login: "nekdo3", Password: ""}, AuthorID: 0x2, AuthorName: "Nekdo Author 2"}}
		checkGotWant("ListAuthors(0,6)", got, want, authorsEqual)

	}
}
