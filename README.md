# Honest Backend Engineer Technical Assessment

>
> 🤓 This repository contains a technical assessment to be used by candidates for the Backend Engineer position at Honest.
>

## Objective

Welcome to Honest, and thank you for taking the time to take part in our recruitment process. We're working hard to make
this a transparent, inclusive and positive process that lets everyone be their best (and have fun!) If you have any questions
or concerns please don't hesitate to raise them with your interviewers, who will be more than happy to help 🙂

For more information about our recruitment process please see our public [Honest Engineering Recruitment Process](https://www.notion.so/honestbank/Honest-Engineering-Recruitment-Process-0ddc3af604c14c6eba20399374edfd47)
page.

### Disclaimer

The scenario below is entirely fictitious and any resemblance to characters real or imaginary is purely coincidental. Please
don't sue us!

We will not use your code submission for any purpose other than evaluating your fit for our team. No engineers were 
harmed during the creation of this technical assessment ✌️.

## Assessment

### Background

An engineer on our team started work on a Decision Engine to approve/decline credit card applicants. However, they have
been unable to complete the project. Your task is to help complete the project, and improve the overall code quality as
you see fit.

---

### Requirements

Your task is to add rules as specified
and have the engine return "approved" or "declined" for the data provided. The engine should be able to be triggered by
a POST request and must handle errors gracefully.

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

## Evaluation Criteria

1. Problem Understanding
1. Problem Solving
1. Testing
1. Effective Architecture/Design

## Submission Instructions

* Create a new branch, check it out and add commits to that branch.
* Once you're done with the assessment create a patch onto `main` - to do this, run the following command:

  ```shell
  > git format-patch main
  ```

* The above command will produce some `.patch` files, simply archive them (ZIP, RAR, tarballs are all accepted).
