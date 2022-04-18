# Leopard
A go web framework

---

## Todo:
 - [ ] Database related tasks
    - [ ] Migrations
    - [ ] Query builder
    - [ ] ...
 - [ ] Making Leopard more secure
    - [ ] CSRF protection
    - [ ] Session management
 - [ ] Making Leopard more customizable
    - [ ] Custom error handling
    - [ ] ...
 - [ ] Creating a CLI
    - [ ] ...
 - [ ] Create docs
 
## Usage
You can simply create a new leopard project by adding this to your main function:

```go
func main() {
    app, err := leopard.New()
    
    if err != nil {
        panic(err)
    }
    
    // Add everything here
    
    err = app.Serve()
    
    if err != nil {
        panic(err)
    }
}
```
Creating a project will later be simplified a lot with the CLI.  
