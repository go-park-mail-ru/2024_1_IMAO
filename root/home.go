package root

/*
func (active *ActiveUsers) Root(writer http.ResponseWriter, request *http.Request) {
	authorized := false
	session, err := request.Cookie("session_id")

	if err == nil && session != nil {
		authorized = active.sessionExists(session.Value)
	}

	if authorized {
		user, err := active.GetUserBySession(session.Value)

		if err != nil {
			fmt.Println("No such session")
			return
		}

		writer.Write([]byte("You authorized as user "))
		writer.Write([]byte(user.Email))
	} else {
		writer.Write([]byte("You have no authorization"))
	}
} */
