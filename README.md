# gmsend

Go library for sending emails through a gmail account.  Support is included for a default
JSON configuration.

The simplest usage assumes the use a nil to imply a default *UserEmail, 
the content of which comes from a JSON file located in one of:

	1) a file specified in the environment variable GMSEND
	2) $HOME/.gmsend.json
	3) /opt/etc/gmsend.json
	
Whether using the default nil pointer to fetch from JSON, or supplied explicity, 
the UserMail will end up like:

```
SMTPAuthentication { 	UserName    : "pierre@gmail.com", 
						Password    : "7nuit,pas3,stp",
						EmailServer : "smtp.gmail.com"
						Port        : 587 }
```

Then you need a message.

```
Message {
	Subject        : "meeting tonight", // required 
	Content        : someText,          // required
	From           : string,            // defaults to the UserName above
	HideRecipients : false,             // default false
}
```

Finally supply recipients and hopefully the error return is nil.

```
recipients := []string{ "i812@gmail.com","ralph@malph.com"}
err := gmail_send.Send(nil, theMessage, recipients)
```


This modifies work from 
[Nathan LeClaire blog](https://nathanleclaire.com/blog/2013/12/17/sending-email-from-gmail-using-golang/).

