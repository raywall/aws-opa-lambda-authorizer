package rules

default allow=false

allow = true {
    input.method == "GET"
    # regex.match("^[/]?user[/]?[^/?]*", input.path)
}

allow = true {
    input.method == "POST"
    input.userType == "admin"
    regex.match("^[/]?user[/]?[^/?]*", input.path)
}

allow = true {
    input.method == "GET"
	input.path == "hello"
}