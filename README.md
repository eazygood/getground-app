## Running the application
We added basic project skeleton with docker-compose. (optional)
Feel free to refactor but provide us with good instructions to start the application
```
make docker-up
```

Update the `docker/mysql/dump.sql` to initialize the mysql database


## Summary

- This task should take approximately 90-120 minutes.
- You can use the provided skeleton project structure (optional). If you are using the following supported languages:
  - `golang`
  - `java`

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
            "guest_id": "string",
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

### Count number of empty seats

```
GET /empty_seats
response:
{
    "empty_seats": int
}
```
