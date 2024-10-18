package codegen

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"autocrud/src/config"
)

type TestGenerateBuffer struct {
	buffer bytes.Buffer
}

func (g *TestGenerateBuffer) CreateBuffer(destPath string) (io.Writer, error) {
	return &g.buffer, nil
}

func (g *TestGenerateBuffer) Close() {}

func TestGenerateMain(t *testing.T) {
	witness := TestGenerateBuffer{}
	BeginTest(&witness)

	expected := `/* This code is autogenerated by Autocrud v0.1.0 */

package main

import (
    "test/src/controller"
)

func main() {
    server := controller.NewController()
    server.Run(
        "localhost:8080",

        controller.NewUserController(),

    )
}

`

	conf := config.Config{
		Schema: config.Schema{
			Tables: []config.TableSchema{
				{
					Name: "user",
					Fields: []config.FieldSchema{
						{
							Name: "user_id",
							Type: "int",
						},
						{
							Name: "created_at",
							Type: "timestamp",
						},
					},
				},
			},
		},
	}

	err := GenerateMain("output.go", "test", conf)

	assert.NoError(t, err)

	assert.Equal(t, expected, witness.buffer.String())
}

func TestGenerateModel(t *testing.T) {
	witness := TestGenerateBuffer{}
	BeginTest(&witness)

	expected := `/* This code is autogenerated by Autocrud v0.1.0 */

package models



type User struct {

    User_id int

}


`

	table := config.TableSchema{
		Name: "user",
		Fields: []config.FieldSchema{
			{
				Name: "user_id",
				Type: "int",
			},
		},
	}
	err := GenerateModel("output.go", table)

	assert.NoError(t, err)

	assert.Equal(t, expected, witness.buffer.String())
}

func TestGenerateModelImportTime(t *testing.T) {
	witness := TestGenerateBuffer{}
	BeginTest(&witness)

	expected := `/* This code is autogenerated by Autocrud v0.1.0 */

package models


import (
    "time"
)


type User struct {

    User_id int

    Created_at time.Time

}


`

	table := config.TableSchema{
		Name: "user",
		Fields: []config.FieldSchema{
			{
				Name: "user_id",
				Type: "int",
			},
			{
				Name: "created_at",
				Type: "timestamp",
			},
		},
	}
	err := GenerateModel("output.go", table)

	assert.NoError(t, err)

	assert.Equal(t, expected, witness.buffer.String())
}

func TestGenerateDAO(t *testing.T) {
	witness := TestGenerateBuffer{}
	BeginTest(&witness)

	expected := `/* This code is autogenerated by Autocrud v0.1.0 */

package dao

import (
    "database/sql"

    _ "github.com/mattn/go-sqlite3"

    "test/src/models"
)

type UserDAO struct {}

func (r UserDAO) GetResource() ([]models.User, error) {
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        return nil, err
    }
    defer db.Close()

    query := "SELECT * FROM user;"

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    resources := make([]models.User, 0)

    for rows.Next() {
        resource := models.User{}
        err := rows.Scan(

            &resource.User_id,

            &resource.Created_at,

        )
        if err != nil {
            return nil, err
        }

        resources = append(resources, resource)
    }

    return resources, nil
}

func (r UserDAO) CreateResource(in *models.User) (int, error) {
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        return 0, err
    }
    defer db.Close()

    query := ` + "`" + `INSERT INTO user (created_at)
    VALUES (?);` + "`" + `
    result, err := db.Exec(
        query,




        in.Created_at,


    )
    if err != nil {
        return 0, err
    }

    newId, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int(newId), nil
}

func (r UserDAO) UpdateResource(in *models.User) error {
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        return err
    }
    defer db.Close()

    query := ` + "`" + `UPDATE user SET

    created_at=?

    WHERE user_id = ?
    ;` + "`" + `

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(

        in.Created_at,
in.User_id,
    )
    if err != nil {
        return err
    }

    return nil
}

func (r UserDAO) GetResourceById(resourceId string) (*models.User, error) {
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        return nil, err
    }
    defer db.Close()

    query := "SELECT * FROM user WHERE user_id=?;"

    resource := models.User{}

    err = db.QueryRow(query, resourceId).Scan(

        &resource.User_id,

        &resource.Created_at,

    )
    if err != nil {
        return nil, err
    }

    return &resource, nil
}

func (r UserDAO) DeleteResourceById(resourceId string) error {
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        return err
    }
    defer db.Close()

    stmt, err := db.Prepare("DELETE FROM user WHERE user_id = ?;")
    if err != nil {
        return err
    }

    _, err = stmt.Exec(resourceId)
    if err != nil {
        return err
    }

    return nil
}
`

	daoData := DAOData{
		ProjectName: "test",
		Table: config.TableSchema{
			Name: "user",
			Fields: []config.FieldSchema{
				{
					Name:         "user_id",
					Type:         "int",
					IsPrimaryKey: true,
				},
				{
					Name: "created_at",
					Type: "timestamp",
				},
			},
		},
		DatabasePath: "./test.db",
	}

	err := GenerateDAO("output.go", daoData)

	assert.NoError(t, err)

	assert.Equal(t, expected, witness.buffer.String())
}

func TestGenerateController(t *testing.T) {
	witness := TestGenerateBuffer{}
	BeginTest(&witness)

	expected := `/* This code is autogenerated by Autocrud v0.1.0 */

package controller

import (
    "net/http"
    "log"
    "strconv"

    "github.com/gin-gonic/gin"

    "test/src/dao"
    "test/src/models"
)

type UserController struct {
    UserDAO dao.UserDAO
}

func NewUserController() *UserController {
    return &UserController{
        UserDAO: dao.UserDAO{},
    }
}

func (c UserController) GetResource(ctx *gin.Context) {
    resources, err := c.UserDAO.GetResource()
    if err != nil {
        log.Printf("error %v\n", err)
        ctx.JSON(http.StatusInternalServerError, nil)
        return
    }

    ctx.JSON(http.StatusOK, resources)
}

func (c UserController) CreateResource(ctx *gin.Context) {
    in := models.User {}

    err := ctx.BindJSON(&in)
    if err != nil {
        log.Printf("error %v\n", err)
        ctx.JSON(http.StatusInternalServerError, nil)
        return
    }

    log.Printf("received %v\n", in)

    resourceId, err := c.UserDAO.CreateResource(&in)
    if err != nil {
        log.Printf("error %v\n", err)
        ctx.JSON(http.StatusInternalServerError, nil)
        return
    }

    ctx.JSON(http.StatusOK, resourceId)
}

func (c UserController) GetResourceById(ctx *gin.Context) {
    resourceId := ctx.Param("id")

    resource, err := c.UserDAO.GetResourceById(resourceId)
    if err != nil {
        log.Printf("error %v\n", err)
        ctx.JSON(http.StatusNotFound, nil)
        return
    }

    ctx.JSON(http.StatusOK, resource)
}

func (c UserController) UpdateResource(ctx *gin.Context) {
    in := models.User {}

    err := ctx.BindJSON(&in)
    if err != nil {
        log.Printf("error %v\n", err)
        ctx.JSON(http.StatusInternalServerError, nil)
        return
    }

    log.Printf("received %v\n", in)

    paramId, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        log.Printf("error %v\n", err)
        ctx.JSON(http.StatusBadRequest, nil)
        return
    }

    if in.User_id != paramId {
        log.Printf("error incompatible id: %d vs %d.\n", in.User_id, paramId)
        ctx.JSON(http.StatusForbidden, nil)
        return
    }

    err = c.UserDAO.UpdateResource(&in)
    if err != nil {
        log.Printf("error %v\n", err)
        ctx.JSON(http.StatusInternalServerError, nil)
        return
    }

    ctx.JSON(http.StatusOK, nil)
}

func (c UserController) DeleteResourceById(ctx *gin.Context) {
    paramId := ctx.Param("id")

    err := c.UserDAO.DeleteResourceById(paramId)
    if err != nil {
        log.Printf("error %v\n", err)
        ctx.JSON(http.StatusInternalServerError, nil)
        return
    }

    ctx.JSON(http.StatusOK, nil)
}

func (c UserController) RegisterResource(controller *Controller) {
    controller.Resources = append(controller.Resources, c)

    controller.Router.GET("/user", c.GetResource)
    controller.Router.POST("/user", c.CreateResource)
    controller.Router.GET("/user/:id", c.GetResourceById)
    controller.Router.PATCH("/user/:id", c.UpdateResource)
    controller.Router.DELETE("/user/:id", c.DeleteResourceById)
}
`

	table := config.TableSchema{
		Name: "user",
		Fields: []config.FieldSchema{
			{
				Name:         "user_id",
				Type:         "int",
				IsPrimaryKey: true,
			},
			{
				Name: "created_at",
				Type: "timestamp",
			},
		},
	}
	err := GenerateController("output.go", "test", table)

	assert.NoError(t, err)

	assert.Equal(t, expected, witness.buffer.String())
}

func TestGenerateControllerRouter(t *testing.T) {
	witness := TestGenerateBuffer{}
	BeginTest(&witness)

	expected := `/* This code is autogenerated by Autocrud v0.1.0 */

package controller

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"

    _ "test/src/models"
    _ "test/src/dao"
)

type ResourceController interface {
    RegisterResource(*Controller)
}

type Controller struct {
    Router    *gin.Engine
    Resources []ResourceController
}

func NewController() *Controller {
    router := gin.Default()
    router.Use(cors.Default())

    return &Controller{
        Router:    router,
        Resources: make([]ResourceController, 0),
    }
}

func (c *Controller) Run(addr string, resources ...ResourceController) {
    for _, controller := range resources {
        controller.RegisterResource(c)
    }
        c.Router.Run(addr)
}
`

	err := GenerateControllerRouter("output.go", "test")

	assert.NoError(t, err)

	assert.Equal(t, expected, witness.buffer.String())
}
