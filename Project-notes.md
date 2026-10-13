# Project Notes

## Project Setup

Ran this command:

```bash
go mod init social-app
```

> TODO: change the module name to the actual GitHub repo.

### Setting up development patterns

1. `bin/` — contains the compiled application
2. `cmd/` — contains the executables or entry points
3. `cmd/api/`
4. `cmd/migrate/`
5. `cmd/migrate/migrations/`
6. `internal/` — interacts with Postgres, data validation, rate limiting, etc.
7. `docs/`
8. `scripts/`

---

# Go Study Notes: Structs, Functions, Methods, Pointers, and HTTP Servers

## 1. Structs

A struct groups related pieces of data into one type.

```go
type Person struct {
    name string
    age  int
}
```

Here, `Person` is a struct type with two fields:

- `name` — a string
- `age` — an integer

### Creating a struct

```go
person := Person{
    name: "John",
    age:  25,
}
```

### Accessing fields

Use the dot (`.`) operator.

```go
fmt.Println(person.name)
fmt.Println(person.age)
```

> **Remember:** A struct is a way to group related data together.

---

## 2. Functions

A function is a reusable block of code that performs a task.

```go
func greet(name string) {
    fmt.Println("Hello", name)
}
```

Call the function by using its name and providing an argument:

```go
greet("John")
```

### Functions with return values

```go
func add(a int, b int) int {
    return a + b
}
```

Usage:

```go
result := add(2, 3)
fmt.Println(result) // 5
```

The final `int` in the function declaration specifies the return type.

---

## 3. Methods

A method is a function associated with a particular type through a **receiver**.

```go
type Person struct {
    name string
}

func (p Person) introduce() {
    fmt.Println("My name is", p.name)
}
```

Usage:

```go
person := Person{name: "John"}
person.introduce()
```

### Understanding the method declaration

```go
func (p Person) introduce()
```

| Part        | Meaning                        |
| ----------- | ------------------------------ |
| `func`      | Declares a function or method  |
| `p`         | Receiver variable              |
| `Person`    | Receiver type                  |
| `introduce` | Method name                    |

The receiver allows the method to access the value associated with that type.

### Function vs. method

```go
func greet() {
    fmt.Println("Hello")
}
```

This is a regular function.

```go
func (p Person) greet() {
    fmt.Println("Hello", p.name)
}
```

This is a method associated with `Person`.

---

## 4. Pointers

A pointer stores the memory address of a value.

### The address operator: `&`

```go
person := Person{name: "John"}

ptr := &person
```

`&person` gives you the address of `person`.

### The pointer type: `*Person`

```go
var ptr *Person
```

This declares a variable that can hold a pointer to a `Person`.

The meaning of `*` depends on where it is used:

```go
var ptr *Person // Pointer type
ptr = &person   // Get the address of person
```

### Accessing data through a pointer

```go
person := Person{name: "John"}
ptr := &person

fmt.Println(ptr.name)
```

Go automatically dereferences the pointer when accessing struct fields, so `ptr.name` works.

You can also write:

```go
fmt.Println((*ptr).name)
```

Both access the same field.

### Modifying the original value

```go
type Person struct {
    age int
}

func updateAge(p *Person) {
    p.age = 30
}
```

Usage:

```go
person := Person{age: 20}

updateAge(&person)

fmt.Println(person.age) // 30
```

Because `p` points to the original `person`, the function can modify the original value.

> **Important:** A pointer does not automatically change a value. It gives you access to the original value, which you can then modify.

---

## 5. Value Receivers vs. Pointer Receivers

### Value receiver

```go
func (p Person) updateAge() {
    p.age = 30
}
```

The receiver gets a **copy** of the struct. Changing that copy does not modify the original struct.

### Pointer receiver

```go
func (p *Person) updateAge() {
    p.age = 30
}
```

The receiver points to the **original** struct, so changing its fields modifies the original.

| Value receiver                         | Pointer receiver                  |
| -------------------------------------- | --------------------------------- |
| `func (p Person)`                      | `func (p *Person)`                |
| Receives a copy                        | Receives a pointer                |
| Changes generally affect only the copy | Can modify the original           |
| Useful when mutation is unnecessary    | Useful when mutation is needed    |

---

## 6. Understanding Your Application Structs

Your code:

```go
type config struct {
    addr string
}

type application struct {
    config config
}
```

You have two struct types.

### `config`

```go
type config struct {
    addr string
}
```

Stores configuration data, such as the server address.

### `application`

```go
type application struct {
    config config
}
```

Stores a `config` value inside an `application`.

This is called **composition**: one struct contains another struct.

The relationship looks like this:

```text
application
    |
    └── config
          |
          └── addr: ":8080"
```

### Creating the values

```go
cfg := config{
    addr: ":8080",
}

app := &application{
    config: cfg,
}
```

- `cfg` is a `config` value.
- `app` is a pointer to an `application` value.
- `app.config.addr` accesses the address stored in the nested configuration.

---

## 7. Understanding the `run()` Method

Your code:

```go
func (app *application) run() error {
    mux := http.NewServeMux()

    srv := &http.Server{
        Addr:    app.config.addr,
        Handler: mux,
    }

    log.Printf("Server has start at %s", app.config.addr)

    return srv.ListenAndServe()
}
```

### Breaking down the method declaration

```go
func (app *application) run() error
```

| Part           | Meaning                      |
| -------------- | ---------------------------- |
| `func`         | Declares the method          |
| `app`          | Receiver variable            |
| `*application` | Pointer receiver type        |
| `run`          | Method name                  |
| `error`        | Return type                  |

This method belongs to the `application` type and uses a pointer receiver.

`*application` is not another struct. It means a **pointer to an `application` value**.

### Step 1: Create a router

```go
mux := http.NewServeMux()
```

Creates a new HTTP request multiplexer, commonly called a router. It matches incoming requests to registered handlers.

At this point, the router has no routes registered.

### Step 2: Create the HTTP server

```go
srv := &http.Server{
    Addr:    app.config.addr,
    Handler: mux,
}
```

This creates an HTTP server configuration.

- `Addr` specifies the address to listen on, such as `:8080`.
- `Handler` specifies which handler processes incoming requests.
- `mux` is the router responsible for dispatching requests.

The `&` creates a pointer to the `http.Server` value.

### Step 3: Log a message

```go
log.Printf("Server has start at %s", app.config.addr)
```

Prints a message containing the configured address. A more natural message would be:

```go
log.Printf("Server starting on %s", app.config.addr)
```

Note that this log runs *before* the server successfully starts listening.

### Step 4: Start the server

```go
return srv.ListenAndServe()
```

This starts listening for HTTP requests.

- If the server starts and continues running, the call blocks while it serves requests.
- If starting the server fails, it returns an error.
- If the server is shut down, it also returns an error, commonly `http.ErrServerClosed`.

The `return` passes that error back to whoever called `run()`.

---

## 8. Understanding `main()`

Your code:

```go
func main() {
    cfg := config{
        addr: ":8080",
    }

    app := &application{
        config: cfg,
    }

    log.Fatal(app.run())
}
```

### Execution flow

1. `main()` begins execution.
2. A `config` value is created with the address `:8080`.
3. An `application` value is created using that configuration.
4. `app.run()` is called.
5. The server starts listening on the configured address.
6. When `run()` returns an error, `log.Fatal()` logs it and exits the program with a non-zero status.

### What does `log.Fatal()` do?

```go
log.Fatal(app.run())
```

`log.Fatal()` logs its arguments and then calls `os.Exit(1)`.

Because `app.run()` is evaluated first, `log.Fatal()` waits while `ListenAndServe()` is running. It only gets the returned error when `run()` finishes.

> **Important:** `log.Fatal()` is not what keeps the server running. `ListenAndServe()` does that.

---

## 9. Complete Execution Flow

```text
main()
  |
  |-- Create config with addr ":8080"
  |
  |-- Create application with config
  |
  |-- Call app.run()
          |
          |-- Create ServeMux router
          |
          |-- Configure HTTP server
          |
          |-- Log startup message
          |
          |-- Call ListenAndServe()
                  |
                  |-- Listen for HTTP requests
                  |
                  |-- Return an error if it stops
  |
  |-- log.Fatal() logs the returned error
      and exits the program
```

---

## 10. One Important Detail: Your Router Has No Routes

Your current code creates a router:

```go
mux := http.NewServeMux()
```

But it does not register any handlers.

So if you visit `http://localhost:8080/`, the server may be running correctly, but the request will typically receive a `404 Not Found` response.

For example, you can add a simple route:

```go
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, world!")
})
```

Place this immediately after creating `mux`:

```go
mux := http.NewServeMux()

mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, world!")
})
```

Now requests to `/` will receive a response containing `Hello, world!`.

---

## Quick Revision

| Concept               | Remember                                   |
| --------------------- | ------------------------------------------ |
| Struct                | Groups related data                        |
| Function              | Reusable block of code                     |
| Method                | Function associated with a type            |
| Receiver              | Connects a method to a type                |
| `&value`              | Gets the address of a value                |
| `*Type`               | Declares a pointer type                    |
| Value receiver        | Works with a copy                          |
| Pointer receiver      | Can modify the original value              |
| `application`         | Your application's main struct             |
| `config`              | Stores configuration data                  |
| `app.config.addr`     | Accesses the configured server address     |
| `http.NewServeMux()`  | Creates a request router                   |
| `http.Server`         | Holds HTTP server configuration            |
| `ListenAndServe()`    | Starts serving HTTP requests               |
| `log.Fatal()`         | Logs and exits when execution reaches it   |

> **The key idea:** Your `main()` function creates the application's configuration and `application` value, then calls `app.run()`. The `run()` method uses that configuration to create and start an HTTP server.


1. Why do we use & before http.Server?
The key is understanding the difference between creating a value and creating a pointer to a value.
Without &:
srv := http.Server{
    Addr: app.config.addr,
}


Here, srv holds an http.Server value.
With &:
srv := &http.Server{
    Addr: app.config.addr,
}


Here, srv holds a pointer to an http.Server value.
Think of it this way
Without &
srv
http.Server value



With &
srv
Pointer (address)
http.Server value



Conceptual illustration: the pointer lets you refer to the server value indirectly.

2. Why use a pointer for an HTTP server?
There are three useful things to understand.
Reason 1: Go's HTTP server methods use pointer receivers.
For example, ListenAndServe() is defined on *http.Server, meaning a pointer to an HTTP server.
Conceptually, its method declaration looks like this:
func (srv *Server) ListenAndServe() error


When you write:
srv := &http.Server{
    Addr: app.config.addr,
}

return srv.ListenAndServe()


srv is already a pointer, so calling ListenAndServe() is straightforward.
Reason 2: The server has state.
An HTTP server manages things such as its listener, connections, and serving lifecycle. Using a pointer allows its methods to work with the same server instance rather than a copied server value.
Reason 3: It is the conventional way to construct an HTTP server.
You'll commonly see this pattern in Go:
srv := &http.Server{
    Addr:    app.config.addr,
    Handler: mux,
}


It creates a server value and gives you a pointer to it.
3. Is & actually required here?
No, not strictly. This is an important distinction.
You could also write:
srv := http.Server{
    Addr: app.config.addr,
}

return srv.ListenAndServe()


Go automatically takes the address of an addressable variable when calling a method with a pointer receiver. In this case, srv.ListenAndServe() works even though srv is not explicitly declared as a pointer.
So why do Go developers usually write &http.Server{...}?
- It makes it clear that srv is a pointer.
- It follows the conventional pattern for constructing an HTTP server.
- It makes it easy to pass the same server pointer to other functions that expect *http.Server.
4. What exactly happens in your code?
func (app *application) run() error {
    srv := &http.Server{
        Addr: app.config.addr,
    }

    return srv.ListenAndServe()
}


Follow the execution:
1. app is a pointer to your application.
2. app.config.addr retrieves the configured server address, such as ":8080".
3. &http.Server{...} creates an HTTP server value and obtains its pointer.
4. srv stores that pointer.
5. srv.ListenAndServe() starts the server.
6. When ListenAndServe() returns, its error is returned by run().
One additional detail: if you don't specify Handler, the server uses its default handler, http.DefaultServeMux.
The main takeaway
Don't think of & as something you must always use when creating a struct.
Instead, ask yourself:
- Do I need a pointer to this value?
- Do the methods I want to call use pointer receivers?
- Do I want to pass the same instance around without copying the struct?
For http.Server, using &http.Server{...} is a conventional and clear choice.