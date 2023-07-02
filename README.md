# hotels-data-merge
This service implement a golang http server to merge hotel source from multi source

## Repository Structure: 
```
➜  hotels-data-merge git:(main) ✗ tree .              
.
├── Dockerfile                              // for docker build
├── Makefile                                // cmd shortcut
├── README.md                               // README file
├── app
│   └── main.go                             // main function to start application
├── etc
│   └── local.yaml                          // local config file
├── go.mod                                  // go dependency
├── go.sum
├── internal 
│   ├── config
│   │   └── config.go                       // handle config relate logic
│   ├── context
│   │   ├── context.go                      // generate and store traceID in context
│   │   └── context_test.go                 // unit test file 
│   ├── domain
│   │   └── hotel
│   │       ├── hotel.go                    // different data structure from different data source & our repository data structure & merge logic for our data structure
│   │       ├── hotel_test.go               // unit test file 
│   │       ├── repo.go                     // get from different data source and convert data to our structure
│   │       ├── repo.mock.gen.go            // mock file for unit test
│   │       ├── repo_test.go                // unit test file 
│   │       ├── service.go                  // hotel service for fetch from all source and periodically refresh
│   │       ├── service.mock.gen.go         // mock file for unit test
│   │       └── service_test.go             // unit test file 
│   ├── handler
│   │   ├── hotels.go                       // handler file for endpoint /hotels
│   │   └── hotels_test.go                  // unit test file 
│   ├── log
│   │   └── log.go                          // log related tool
│   └── util
│       ├── util.go                         // util package
│       ├── util.mock.gen.go                // mock file for unit test
│       └── util_test.go                    // unit test file 
└── log                                     // log file
    └── application.log

```

## Requirement:
- docker

## Installation & Run:
- install docker
- ```git clone https://github.com/YunchengHua/hotels-data-merge.git```
- ```make build```
- ```make run``` 

If you would like to check the log of the server 
- make log 

## Request:
There is only one endpoint: GET /hotels
Params:
- **hotel_ids**
- **destination_ids**

Providing at least one of the parameters is required. If both of them are given, then we will use **hotel_ids**.

```
curl -X GET 'http://localhost:8080/hotels?hotel_ids=iJhz,f8c9'
```
```
curl -X GET 'http://localhost:8080/hotels?destination_ids=5432'
```
```
curl -X GET 'http://localhost:8080/hotels?hotel_ids=iJhz,f8c9&destination_ids=5432'
```

## Response
Properties: 
* **id** (string)
* **name** (string)
* **destination_id** (int64)
* **description** (string)
* **booking_conditions** (string array)
* **amenities** (string array)
* **images** (Object)
  * amenities (Object array)
    * url (string)
    * caption (string)
  * rooms (Object array)
    * url (string)
    * caption (string)
  * site (Object array)
    * url (string)
    * caption (string)
* **location** (Object)
  * address (string)
  * city (string)
  * country (string)
  * latitude (string)
  * longitude (string)
* **missing** (bool) 

PS:
- When get by hotel_ids and id is not exist, **missing** will return true and ID will be set in the structure

## Merge Strategy:
- id: no merging
- destination_id: no merging
- name: longer wins
- description: default choose longer wins, could also config to concatenate together
- bookingConditions: append together and remove duplication
- Amenities: append together and remove duplication
- Images: append together and use url as key to remove duplication
- Location: Non-Null value wins, For address & city & country & postalCode longer wins, because consider as more detailed

## Caching Strategy:
- Fetch all data sources when the service start
- Fetch all data sources and update periodically

## Achievements
- Fullfil all the requirements
- Over 70% unit test coverage on main package(folder)
```
ok      github.com/YunchengHua/hotels-data-merge/internal/context       0.152s  coverage: 84.6% of statements
ok      github.com/YunchengHua/hotels-data-merge/internal/domain/hotel  0.769s  coverage: 70.9% of statements
ok      github.com/YunchengHua/hotels-data-merge/internal/handler       0.904s  coverage: 78.3% of statements
```
- Self implement an log wrap on original log package, able to log the req and resp and each request has a unique traceID for filter
- Add build/test/lint pipeline for every pull requests to **main**
- Use docker to build the server for easy deployment

## Future Improvement:
- Use Redis or other method to share the cache across all instances
- Add cache invalidation method to update or remove the out-data data instead of periodically fetching
- Add config allowing to config log setting etc...
- Use tcp service and protobuf, reducing the size of req and resp, saving the bandwidth
- Add integration test to test end to end flow
- Add custom linter