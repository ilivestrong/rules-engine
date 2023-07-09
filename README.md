# Credit Card Rules Engine in Go

>
> 🤓 This repository is my journey into DevOps including deployments on K8s.
>


### Requirements

This repo is a simple implementation of a Go HTTP server, which processes a credit card application and based on 
rules mentioned below, declines or approves it.

1. Implement a `POST` HTTP endpoint that:

   * Receives [a request with JSON body](#request-body).

   * Runs through the [Decision Engine Rules](#decision-rules).

   * Returns [a response with JSON body](#response-body):

     * status = `approved` if all rules are passed.

     * status = `declined` if any rules are failed.

1. Handle errors gracefully, without stopping the process.

### Specifications

#### Request Body

| Fields                   | Type        |
| -----------              | ----------- |
| income                   | number      |
| number_of_credit_cards   | number      |
| age                      | number      |
| politically_exposed      | bool        |
| job_industry_code        | string      |
| phone_number             | string      |

##### Example

```json
{
  "income": 82428,
  "number_of_credit_cards": 3,
  "age": 9,
  "politically_exposed": true,
  "job_industry_code": "2-930 - Exterior Plants",
  "phone_number": "486-356-0375"
}
```

#### Response Body

| Fields                   | Type        |
| -----------              | ----------- |
| status                   | string      |

##### Example

###### Approved:

```json
{
  "status": "approved"
}
```

###### Declined:

```json
{
  "status": "declined"
}
```

#### Decision Rules

The application is approved if it evaluates as `true` on the following rules:

1. The applicant must earn more than 100000.
1. The applicant must be at least 18 years old.
1. The applicant must not hold more than 3 credit cards and their `credit_risk_score` must be `LOW`.
1. The applicant must not be involved in any political activities (must not be a Politically Exposed Person or PEP).
1. The applicant's phone number must be in an area that is allowed to apply for this product. The area code is denoted by first digit of phone number. The allowed area codes are `0`, `2`, `5`, and `8`.
1. A pre-approved list of phone numbers should cause the application to be automatically approved without evaluation of the above rules. This list must be able to be updated at runtime without needing to restart the process.

#### External Data Sources

Values for the `credit_risk_score` field can be retrieved by calling the existing functions in the provided `risk` module.
