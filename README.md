# finance-limits
## Context
In finance, it's common for accounts to have so-called "velocity limits". 
Each attempt to load funds will come as a single-line JSON payload, structured as follows:

```json
{
  "id": "1234",
  "customer_id": "1234",
  "load_amount": "$123.45",
  "time": "2018-01-01T00:00:00Z"
}
```

Each customer is subject to three limits:

- A maximum of $5,000 can be loaded per day
- A maximum of $20,000 can be loaded per week
- A maximum of 3 loads can be performed per day, regardless of amount

The return is a json string for each load telling it's accepted or not.
```json
{ "id": "1234", "customer_id": "1234", "accepted": true }
```


## Assumptions
