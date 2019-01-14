#Go YAML REST API 

A RESTful API that accepts and responds with YAML data. This project was written shortly after I began learning Go. 


##Routes
POST https//localhost:8000/applications  

All fields are required, you may include more than one maintainer.

GET one https//localhost:8000/applications/{title} 

Full title required but queries are not case sensitive. All matching results will be returned.

GET all https//localhost:8000/applications

This will retrieve all records.