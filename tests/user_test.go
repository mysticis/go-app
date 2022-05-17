package tests

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/mysticis/go-dcktst-demo/demo"
	"github.com/stretchr/testify/require"
)

const DBUrl = "postgres://postgres:secret@localhost:5432/spin"

var testQueries *demo.Queries

func TestMain(m *testing.M) {

	db, err := sql.Open("pgx", DBUrl)
	if err != nil {
		log.Fatalf("couldn't connect to database, %v", err)
	}

	testQueries = demo.New(db)

	os.Exit(m.Run())
}

func createUser(t *testing.T) demo.User {

	args := demo.CreateUserParams{
		Name:  "John",
		Email: "john@test.com",
		Phone: "980-0986",
	}

	user, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.Name, user.Name)
	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.Phone, user.Phone)
	require.NotZero(t, user.ID)
	return user
}

func TestCreateUser(t *testing.T) {

	//initialize database
	err := testQueries.DeleteAllUsers(context.Background())

	require.NoError(t, err)

	args := demo.CreateUserParams{
		Name:  "Jack Still",
		Email: "jack@test.com",
		Phone: "9088-098",
	}

	testUser, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, testUser)

	require.Equal(t, args.Name, testUser.Name)
	require.Equal(t, args.Email, testUser.Email)
	require.Equal(t, args.Phone, testUser.Phone)
	require.NotZero(t, testUser.ID)

}

func TestGetUser(t *testing.T) {

	//initialize DB
	err := testQueries.DeleteAllUsers(context.Background())

	require.NoError(t, err)

	firstUser := createUser(t)
	newUser, err := testQueries.GetUser(context.Background(), firstUser.ID)

	require.NoError(t, err)

	require.NotEmpty(t, newUser)

	require.Equal(t, firstUser.ID, newUser.ID)
	require.Equal(t, firstUser.Name, newUser.Name)
	require.Equal(t, firstUser.Email, newUser.Email)
	require.Equal(t, firstUser.Phone, newUser.Phone)

}

func TestUpdateUser(t *testing.T) {

	//initialize DB
	err := testQueries.DeleteAllUsers(context.Background())

	require.NoError(t, err)

	initialUser := createUser(t)

	args := demo.UpdateUserParams{
		ID:    initialUser.ID,
		Name:  initialUser.Name,
		Email: "john@changed.org",
		Phone: "980-0987",
	}

	user, err := testQueries.UpdateUser(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, user)
	require.Equal(t, args.ID, user.ID)

	require.Equal(t, args.Name, user.Name)
	require.Equal(t, args.Phone, user.Phone)
	require.Equal(t, args.Email, user.Email)

}

func TestDeleteUser(t *testing.T) {

	//initialize DB
	err := testQueries.DeleteAllUsers(context.Background())

	require.NoError(t, err)

	userToBeDeleted := createUser(t)

	err = testQueries.DeleteUser(context.Background(), userToBeDeleted.ID)

	require.NoError(t, err)

	user, err := testQueries.GetUser(context.Background(), userToBeDeleted.ID)

	require.Error(t, err)

	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, user)

}

func TestListUsers(t *testing.T) {

	//initialize DB
	err := testQueries.DeleteAllUsers(context.Background())

	require.NoError(t, err)

	createUser(t)
	createUser(t)
	createUser(t)

	allUsers, err := testQueries.ListUsers(context.Background())

	require.NoError(t, err)

	require.Len(t, allUsers, 3)

	for _, tasks := range allUsers {
		require.NotEmpty(t, tasks)
	}
}

func TestCleanUp(t *testing.T) {

	err := testQueries.DeleteAllUsers(context.Background())

	require.NoError(t, err)

}
