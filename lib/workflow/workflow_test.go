package workflow

import (
    "errors"
    "fmt"
    "testing"
)

func TestNewContext(t *testing.T) {
    c := New("test")
    if c.Name != "test" {
        t.Fatalf("Context name is %s, not test", c.Name)
    }
}

func TestNilHandler(t *testing.T) {
    c := New("test")

    c.Use(NilHandleFunc)

    err := c.Run()

    if err != nil {
        t.Fatal(err)
    }
}

func TestUsefulHandler(t *testing.T) {
    c := New("test")

    c.Use(func (c *Context) (err error) {
        println("Test message!")
        return
    })

    err := c.Run()
    if err != nil {
        t.Fatal(err)
    }
}

func TestManyHandlers(t *testing.T) {
    c := New("test")

    var i int
    for i = 0; i < 15; i++ {
        c.Use(func (c *Context) (err error) {
            fmt.Printf("Hi there %d\n", i)
            return nil
        })
    }

    err := c.Run()
    if err != nil {
        t.Fatal(err)
    }
}

func TestExpectedError(t *testing.T) {
    c := New("test")

    c.Use(func(c *Context) (err error) {
        return errors.New("Expected error")
    })

    err := c.Run()
    if err == nil {
        t.Fatal("Expected error, not nil")
    }
}
