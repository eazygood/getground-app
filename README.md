## Running the application
We added basic project skeleton with docker-compose. (optional)
Feel free to refactor but provide us with good instructions to start the application
```
make
```

 `docker/mysql/dump.sql` has initializion of the mysql database


## Summary

- This task has taken approximately 2 days due to some issue with testing.
- Languages taken:
  - `golang`

## Requirements

- Please use MySQL version 5.7 as a database
- We are expecting well tested and structured code
- Good documentation to get us started and understand your implementation

**NOTE: Please DO NOT upload this task or your implementation to a public repository. If you must use version control, consider using a private repository instead.**

## Task Overview

Hopefully this tech task allows you to strut your stuff as much as you decide to!

We'd like to implement a guestlist service for the GetGround year end party!
We haven't decide on the venue yet so the number of tables and the capacity are subject to change.

When the party begins, guests will arrive with an entourage. This party may not be the size indicated on the guest list. 
However, if it is expected that the guest's table can accommodate the extra people, then the whole party should be let in. Otherwise, they will be turned away.
Guests will also leave throughout the course of the party. Note that when a guest leaves, their accompanying guests will leave with them.

At any point in the party, we should be able to know:
- Our guests at the party
- How many empty seats there are

## Sample API guide

This is a directional API guide.

### Add a guest to the guestlist

If there is insufficient space at the specified table, then an error should be thrown.

```
POST /guestlist/:guest_id
body: 
{
    "accompanying_guests": int
}
response: 
{
    "guest_id": int
}
```

### Get the guest list

```
GET /guestlist
response: 
{
    "guests": [
        {
            "table_id": int,
            "seats" int
            "guest_id": string,
        }, ...
    ]
}
```

### Guest Arrives

A guest may arrive with an entourage that is not the size indicated at the guest list.
If the table is expected to have space for the extras, allow them to come. Otherwise, this method should throw an error.

```
POST /guests
body:
{
    "name": string
    "accompanying_guests": int
}
response:
{
    "id": int
    "name": string
    "accompanying_guests": int,
    "time_arrived": date,
    "is_arrived": boolean
}
```

### Guest Update

```
PUT /guests/:guest_id
body:
{
    "name": string
    "accompanying_guests": int
    "time_arrived": datetime
    "is_arrived": boolean
}
response:
{
    "message": string
}
```

### Guest Leaves

When a guest leaves, all their accompanying guests leave as well.

```
DELETE /guests/:guest_id

response:
{
    "message": string
}
```

### Get arrived guests

You can provide filter with query parameter `?arrived` to filter out only arrived guests

```
GET /guests
response: 
{
    "guests": [
        {
            "id": int,
            "name": string,
            "accompanying_guests": int,
            "time_arrived": string,
            "is_arrived": boolean
        }
    ]
}
```

### Add Table

```
POST /tables
body:
{
    "seats": int
}
response:
{
    "id": int
    "seats": int
    "guest_id": int | null,
}
```
### Update Table

```
PUT /tables/:table_id
body:
{
    "seats": int
    "guest_id": int
}

response:
{
    "message": string
}
```
### Count number of empty seats from tables

```
GET /empty_seats
response:
{
    "empty_seats": int
}
```

### Delete Table

```
DELETE /tables/:table_id

response:
{
    "message": string
}
```
