# Locker

## small CLI app to store your password in safe

### I made it because was tired of having my passwords in excel.
### maybe add things later but works fine

- copy, build and use. it will get u through initialization.
- leave me your star
---
### after this quick path things u can go 
### Set on your .profile or .bashrc
        export LOCKER_PATH=path/to/lockerDir
        export PATH=$PATH:$LOCKER_PATH
#### in my case is
        export LOCKER_PATH=/home/<username>/.locker
        export PATH=$PATH:$LOCKER_PATH
#### .locker is the folder where i have the binary and the locker.txt (file where u data will be saved encrypted)
---
### After initialization u will get something like this:
      Locker do:
### Here u can do some actions like:
- #### clear -> clear terminal
- #### exit -> exit program
- #### save -> to save changes
- #### ls -> to list all your data 

## Getter

### this will print the item stored with the 'itemname'

    get itemname

### can gen a json report of your data with get -f, it will create a json in the directory where u call the app
### the third param will be the name (without ext), this is optional due that default name is the long date name
    get -f <filename> 

## Setter

### this will set a value

    set anyname anykey anyvalue

### if only value item

    set anyname anyvalue

### if only value item but with spaces ( imagine like metamask words )

#### use "-" as last param is necessary to the program know that all values are one

    set wallet1 a b c d -

### if want to overwrite only the value ( u change your mail password idk)

#### use "-" in as key this will remain the key the same

    set gmail - mysuperpassword

### Also can load your data from json with specific scheme

#### the file has to be recheable from where u calling the the program

     set -f mydata.json

---

        {
            "gmail": { "key": "anykey", "value": "value" },
            "github": { "key": "keyany", "value": "value" },
            "linkedin": { "key": "anykey", "value": "value" }
        }

---

## remove

### to remove a item
    rm itemname
### to remove all items
    rm *    

## ls
     [ linkedin ->  key : value ]
     [ github ->  key : value ]
     [ gmail    ->  key : value ]
#### to only print the names u have store without its data use 
    ls names
---
    gmail    
    github  
    linkedin  

#### To print all items that name matches a substring use "ls find" (imagine have multiple mails) 
    ls find gmail
---
     [ gmail3 ->  key : value ]
     [ gmail4 ->  key : value ]
     [ gmail2 ->  key : value ]
     [ gmail    ->  key : value ]

## CMD

- clear
- exit
- save
- ls 
- rm < name | * >
- get < value >
- set < params >

