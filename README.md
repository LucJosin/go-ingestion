# Bank Scraper – Coding Challenge

This project includes two command-line tools written in Go using the Colly web scraping framework.

## Requirements

- Go 1.21 or higher
- [Colly](https://github.com/gocolly/colly)

## Part 1 – Bank List Scraper

Scrapes bank data from [forbes](https://www.forbes.com/lists/worlds-best-banks/).

### Input

- Output format: `stdout`, `json` or `csv` **(default: stdout)**
- Filename: The output filename **(default: banks)**

```bash
go run cmd/scraper/main.go -o csv
```

### Output

List of bank with these fields:
  - `name`
  - `city`
  - `country`
  - `founded`
  - `rank`
  - `profile` (URL to the bank's page)

### Example

```json
{"name":"1st Source Bank","city":"South Bend","country":"United States","founded":1863,"rank":44,"profile":"https://www.forbes.com/companies/1st-source-bank/?list=worlds-best-banks"}
```

## Part 2 – Bank Profile Scraper

This program accepts a single bank object (as produced by Part 1) and scrapes additional details from the profile page.

### Input

You can use this using `flag` or `stdin`

#### Using flag: 

- Only the profile url

```bash
go run cmd/finder/main.go -p https://www.forbes.com/companies/1st-source-bank/?list=worlds-best-banks
```

#### Using stdin: 

- A single JSON object representing a bank.

```bash
go run cmd/finder/main.go

> {"name":"1st Source Bank","city":"South Bend","country":"United States","founded":1863,"rank":44,"profile":"https://www.forbes.com/companies/1st-source-bank/?list=worlds-best-banks"}
```

### Output

Bank info with the following fields:
  - `name`
  - `ceo` (can be empty)
  - `employees` (can be empty)
  - `lists`: An array of lists the bank appears on, each with `name` and `url`

### Example
```json
{
  "name": "1st Source Bank",
  "lists": [
    {
      "name": "Best Employers for New Grads (2025)",
      "url": "https://www.forbes.com/best-employers-for-new-grads/"
    },
    {
      "name": "World's Best Banks (2025)",
      "url": "https://www.forbes.com/worlds-best-banks/"
    },
    {
      "name": "America's Best Midsize Employers (2025)",
      "url": "https://www.forbes.com/best-midsize-employers/"
    },
    {
      "name": "America's Best Banks (2025)",
      "url": "https://www.forbes.com/americas-best-banks/"
    },
    {
      "name": "America's Best Employers By State (2024)",
      "url": "https://www.forbes.com/best-employers-by-state/"
    },
    {
      "name": "Best-In-State Banks (2024)",
      "url": "https://www.forbes.com/best-in-state-banks/"
    },
    {
      "name": "Best Employers for Diversity (2022)",
      "url": "https://www.forbes.com/best-employers-diversity/"
    }
  ]
}
```