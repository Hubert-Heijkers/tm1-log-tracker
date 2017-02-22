# tm1-log-tracker

The tm1-log-tracker is one, of hopefully many soon, sample application, written against TM1 server's, OData v4.0 compliant, REST API.
It's purpose is to demonstrate how one, using the TM1 Server's REST API, can retrieve and 'listen' to, and subsequently 'process', new messages written into the server's message log.

## Getting Started 

This sample application was written in Go, a.k.a. Go-lang. Please make sure you have at least Go version 1.7.5 to use build this example.
If you don't know what Go is and want to learn more or if you don't have it installed just yet please go to the [golang.org](https://golang.org/) site for more information.

Presuming you have Go up and running you can grab the code for this example and build the app. To do so perform the following steps:
- open a concole/command box
- have go grab the tm1-log-tracker sources: `go get github.com/hubert-heijkers/tm1-log-tracker`
- navigate to the tm1-log-tracker source folder: `cd %GOPATH%\src\github.com\hubert-heijkers\tm1-log-tracker`

Now that you have the code of the example you can take a closer look at the files. The important files are:

- main.go

   This is, as the name suggests, the main file containing the bulk of the example code. The `main` function does the initialization,  
   makes an initial request to the TM1 Server requesting its version number, which has the nice side effect of getting authenticated at the same time,  
   and, last but not least, calls the `TrackCollection` utility as in: 

   ```Go
   client.TrackCollection(tm1ServiceRootURL, "MessageLogEntries", time.Duration(interval)*time.Second, processMessageLogEntries)
   ```

   This is what kick-starts what the example is all about, retrieving whatever is in the MessageLog, by iterating the MessageLogEntries, and subsequently,  
   after waiting a number of seconds defined by the interval, requesting any new message log entries, the delta, since the last request. After every request  
   the `processMessageLogEntries` function is being called. This is the function that gets the chance to process any message log entries returned by the  
   server and the function you most likely end up editing later.  

- utils/odata.go

   This file introduces some extensions to the standard http client making interaction with an OData compliant service in general more easy.  
   It also implements the logic, as defined by the OData protocol specification, to iterate a collection and, to iterate and track a collection,  
   `TrackCollection`, as used by this example.

   As you will notice, if you inspect the `TrackCollection` function, is that the only real difference with just interating a collection is the additional  
   `prefer` header that is being passed on the request, set to the value `odata.track-changes`. This causes the service, the TM1 Server in this case, to return  
   a so called delta link at the end of the response which the application then subsequently can use to find out if any changes have been made to the collection.  
   TM1 Server, to date, only supports track-changes on the message and transaction logs which, due to the nature of these collections, only receive new entries  
   that are being appended to the log. The delta responses are therefore of exactly the same shape then the initial response containing the complete collection.

- .env

   This sample make use of the [godotenv](https://github.com/joho/godotenv) package making grabbing and setting, using a .env file, of environment variables for  
   the application very easy. In the application itself we make use of the following environment variables:

   - `TM1_SERVICE_ROOT_URL`

      The service root URL of the TM1 Server you are planning to track the message log for, typically: `http[s]://tm1server:port/api/v1/`

   - `TM1_USER`

      The user name of the user to be used to log into the TM1 Server specified using the service root URL.

   - `TM1_PASSWORD`

      The password of the user.
 
   - `TM1_TRACKER_INTERVAL`

      The interval, in seconds, between requests to the server (if not specified, or a invalid value is specified, defaults to 5)
   
## Editing the Code

Now that you know where everything is, and perhaps even had a peek at the implementation of the `processMessageLogEntries` function, you likely want to define
some specific processing of message log entries to help you achieve your goals. The bare skeleton of the `processMessageLogEntries` function is:

```Go
func processMessageLogEntries(responseBody []byte) (string, string) {

    // Unmarshal the JSON response
    res := MessageLogEntriesResponse{}
    err := json.Unmarshal(responseBody, &res)
    if err != nil {
        log.Fatal(err)
    }

    // Interate over the message log entries retrieved from the server
    for _, entry := range res.MessageLogEntries {

        // YOUR CODE TO DO ANYTHING WITH A SINGLE MESSAGE LOG ENTRY GOES HERE!

    }

    // Return the nextLink and deltaLink, if there any
    return res.NextLink, res.DeltaLink
}
```

In the example as provided we are only interested in MDX queries that are being processed by the server. This implementation keeps track of the begin
and end times of the MDXViewCreate and dumps those time stamps, including the duration (time it took to create the view), into comma separated output 
and writes it out to the console. 

## Building the Code

Now that you have your code ready the last thing is to build it. Luckily for you we are using Go, so simply type `go build` in your console window and
Go will do the rest for you, grabbing dependencies, building any dependencies if so required and build your application.

After it is done building you have a tm1-log-tracker executable in your source folder. If you rather have the executable 'installed' into the bin folder
of your go path instead then use `go install`. Keep in mind that to be able to run the executable in that case you'll have to move/copy your `.env` file
to that bin folder as well.

## Running the application

Running the application is, after you have correctly set up your environment variables in the `.env` file, as easy as simply running the executable. The
application will run forever unless it runs into a communication issue with the server, the server no longer returns a delta link (which shouldn't happen),
or if you hit Ctrl-C to terminate the application.

Enjoy!
