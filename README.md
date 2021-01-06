## QueryQL
_____
#### QueryQL is a handler for URL requests

### About QueryQL
QueryQL was first developed to help and making it simple to make dynamic querys on an API.

### What it does?
* QueryQL:
    * Guess and change the type of value you're passing through requests;
    * Group similar types of filters to help you to build your own db query;

### Why use QueryQL?
* Simplifies the way you handle parameters;
* Build controlled and advanced filters on your frontend; 

### Types
The following rules will be aplied to change the type of a parameter:
* **object_id type**: The parameter value is automatically typed:
  Regex Rule: ^[0-9a-fA-F]{24}$
  Example: _id=5ea31f55a6dea92aa7d27a8d

* **float type**: The parameter value is automatically typed:
  Regex Rule: ^[+-]?([0-9]*[.,])?[0-9]+$
  Example: product_id=998.21

* **int/long type**: The parameter value is automatically typed:
  Regex Rule: ^[-]?\d+$
  Example: product_id=998

* **bool type**: The parameter value is automatically typed:
  Rule: false or true value.
  Example: is_removed=true

* **datetime type**: The parameter value is automatically typed:
  Regex Rule: (?:\d{4}-\d{2}-\d{2}(?: T([Z]$)?(?:[.+]?\d{2,3})?(?:-?\d{2}:?\d{2}|:\d{2})?)?)$
  Example: product_id=2019-02-03 10:15:21.000

* **string type**: You can type value as string passing a text or wrapping the value in double-quotes.
  Example: product_id="998"


#### Operators:
* ~VALUE -> Contains;
* ~*VALUE -> Ends with;
* ~VALUE* -> Starts with;
* !~VALUE -> Do not contains;
* ~VALUE*EXAMPLE -> Starts with VALUE AND Ends with EXAMPLE;
* !~*VALUE -> Do not ends with;
* !~VALUE* -> Do not starts with;
* !VALUE -> Not equal to;
* \>VALUE -> Greater than;
* \>=VALUE -> Greater or Equal than;
* \<VALUE -> Lower than;
* <=VALUE -> Lower or Equal than;
* 3[]VALUE -> Belongs to default OR group;
* 3[1]VALUE -> Belongs to OR group number 1;
* !3[1]VALUE -> Belongs to OR group number 1 with operation Not Equal;


#### Example of Aggregations:

* ?corporate_name=!~*MARCOS&corporate_name=!~PATRICIA
Results the filter: corporate_name DOES NOT END with MARCOS and DOES NOT CONTAINS patricia

* ?corporate_name=~*MARCOS&corporate_name=~PATRICIA
Results the filter: corporate_name ENDS with MARCOS and CONTAINS patricia

* ?product_id=1&product_id=2
Results the filter: corporate_name IN 1,2

* ?product_id=!1&product_id=!2
Results the filter: corporate_name NOT IN 1,2

* ?product_id=!1&product_id=2
Results the filter: corporate_name NOT IN 1 and corporate_name IN 2

* ?corporate_name=!~PATRICIA[]&client_id=6[]&company_id=2
Results the filter: (corporate_name DOES NOT CONTAINS patricia OR client_id IS EQUAL TO 6) AND company_id IS EQUAL TO 2

* ?corporate_name=!~PATRICIA[1]&client_id=6[1]&product_id=!999[2]&territory_id=0[2]&company_id=2
Results the filter: (corporate_name DOES NOT CONTAINS patricia OR client_id IS EQUAL TO 6) AND (product_id NOT EQUAL TO 999 OR territory_id IS EQUAL TO 0) AND company_id IS EQUAL TO 2
  

### Usage
```
package main

import (
	"github.com/eibrunorodrigues/queryql"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		h := queryql.Handle{}
		err := h.AddList(c.Request.URL.Query(), false)

		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
		} else {
			c.JSON(200, gin.H{
				"message": h.Result,
			})
		}
	})

	if err := r.Run(":5002"); err != nil {
		log.Fatalf("Failed to serve GRPC 5052: %v", err)
	}
}

```