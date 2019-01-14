**#Go YAML API** 

A RESTful API that accepts and responds with YAML data. This project was written shortly after I began learning Go. 


**##Routes**

/applications  POST 

All fields are required, you may include more than one maintainer.

/applications/{title} GET

Full title required but queries are not case sensitive. All matching results will be returned.

/applications GET

This will retrieve all records.